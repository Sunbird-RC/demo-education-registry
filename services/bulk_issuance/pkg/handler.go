package pkg

import (
	"bulk_issuance/config"
	"bulk_issuance/db"
	"bulk_issuance/swagger_gen/models"
	"bulk_issuance/swagger_gen/restapi/operations"
	"bulk_issuance/swagger_gen/restapi/operations/download_file_report"
	"bulk_issuance/swagger_gen/restapi/operations/sample_template"
	"bulk_issuance/swagger_gen/restapi/operations/upload_and_create_records"
	"bulk_issuance/swagger_gen/restapi/operations/uploaded_files"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"log"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func SetupHandlers(api *operations.BulkIssuanceAPI) {
	api.SampleTemplateGetV1BulkSampleSchemaNameHandler = sample_template.GetV1BulkSampleSchemaNameHandlerFunc(downloadSampleFile)
	api.UploadedFilesGetV1BulkUploadedFilesHandler = uploaded_files.GetV1BulkUploadedFilesHandlerFunc(listFiles)
	api.DownloadFileReportGetV1DownloadFileNameHandler = download_file_report.GetV1DownloadFileNameHandlerFunc(downloadReportFile)
	api.UploadAndCreateRecordsPostV1UploadFilesVCNameHandler = upload_and_create_records.PostV1UploadFilesVCNameHandlerFunc(createRecords)
}

func NewGenericJSONResponse(body interface{}) middleware.Responder {
	return &GenericJsonResponse{body: body}
}

type GenericJsonResponse struct {
	body interface{}
}

func (o *GenericJsonResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	bytes, err := json.Marshal(o.body)
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte("JSON Marshalling error"))
	}
	rw.WriteHeader(200)
	rw.Write(bytes)
}

func downloadSampleFile(params sample_template.GetV1BulkSampleSchemaNameParams) middleware.Responder {
	response := sample_template.GetV1BulkSampleSchemaNameOK{}
	resp, err := http.Get(config.Config.Registry.BaseUrl + "api/docs/swagger.json")
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var responseMap map[string]interface{}
	if err = json.Unmarshal(body, &responseMap); err != nil {
		log.Fatal("IN JSON ", err)
	}
	properties := responseMap["definitions"].(map[string]interface{})[params.SchemaName].(map[string]interface{})["properties"].(map[string]interface{})
	headers := make([]string, 0)
	for k, _ := range properties {
		headers = append(headers, k)
	}
	fileName := params.SchemaName
	csvFile, err := os.Create(fileName)
	defer csvFile.Close()

	if err != nil {
		log.Printf("Error creating csv file : %v", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()
	csvwriter.Write(headers)
	csvwriter.Flush()
	csvFile.Close()
	f, _ := os.Open(fileName)
	response.WithContentDisposition("attachment; filename=\"" + params.SchemaName + ".csv\"").WithPayload(f)
	return &response
}

func listFiles(params uploaded_files.GetV1BulkUploadedFilesParams) middleware.Responder {
	response := uploaded_files.GetV1BulkUploadedFilesOK{}
	files, err := db.GetAllUploadedFilesData()
	if err != nil {
		return &response
	}
	response.SetPayload(files)
	return &response
}

func downloadReportFile(params download_file_report.GetV1DownloadFileNameParams) middleware.Responder {
	response := download_file_report.GetV1DownloadFileNameOK{}
	file := db.GetDBFilesUpload(params.FileName)
	var data [][]string
	if err := json.Unmarshal(file.RowData, &data); err != nil {
		log.Fatal("Error : %v", err)
	}
	f, _ := os.Create(params.FileName)
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write(strings.Split(file.Headers, ","))
	for _, d := range data {
		if err := w.Write(d); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
	f.Close()
	f, _ = os.Open(params.FileName)
	response.WithContentDisposition("attachment; filename=\"" + params.FileName + "\"").WithPayload(f)
	return &response
}

type Scanner struct {
	Reader *csv.Reader
	Head   map[string]int
	Row    []string
}

func NewScanner(o io.Reader) Scanner {
	csv_o := csv.NewReader(o)
	header, e := csv_o.Read()
	if e != nil {
		return Scanner{}
	}
	m := map[string]int{}
	for n, s := range header {
		m[strings.TrimSpace(s)] = n
	}
	return Scanner{Reader: csv_o, Head: m}
}

func (o *Scanner) Scan() bool {
	a, e := o.Reader.Read()
	o.Row = a
	return e == nil
}

func createRecords(params upload_and_create_records.PostV1UploadFilesVCNameParams, principal *models.JWTClaimBody) middleware.Responder {
	data := NewScanner(params.File)
	totalSuccess, totalErrors, rows := Process(&data, params.HTTPRequest.Header, params.VCName)
	successFailureCount := map[string]int{
		"success":   totalSuccess,
		"error":     totalErrors,
		"totalRows": totalSuccess + totalErrors,
	}
	response := upload_and_create_records.NewPostV1UploadFilesVCNameOK()
	response.SetPayload(successFailureCount)
	_, fileHeader, _ := params.HTTPRequest.FormFile("file")
	fileName := fileHeader.Filename
	dbUpload := db.DBFileUpload{
		Filename:     fileName,
		TotalRecords: totalSuccess + totalErrors,
		UserID:       principal.PreferredUsername,
		Date:         time.Now().Format("2006-01-02"),
	}
	bytes, _ := json.Marshal(rows)
	fileUpload := db.DBFilesUpload{
		Filename: fileName,
		Headers:  getHeaders(data.Head),
		RowData:  bytes,
	}
	db.CreateDBFileUpload(&dbUpload)
	db.CreateDBFilesUpload(&fileUpload)
	return response
}

func getHeaders(head map[string]int) string {
	headers := make([]string, 0)
	for k := range head {
		headers = append(headers, k)
	}
	sort.SliceStable(headers, func(i, j int) bool {
		return head[headers[i]] < head[headers[j]]
	})
	return strings.Join(headers, ",")
}

func Process(data *Scanner, header http.Header, vcName string) (int, int, [][]string) {
	var (
		totalSuccess    int = 0
		totalErrors     int = 0
		lastColumnIndex int = 0
	)
	rows := make([][]string, 0)
	for data.Scan() {
		currRow := make([]string, 0)
		jsonBody := make(map[string]interface{})
		resp, err := http.Get(config.Config.Registry.BaseUrl + "api/docs/swagger.json")
		if err != nil {
			log.Fatal(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		var responseMap map[string]interface{}
		if err = json.Unmarshal(body, &responseMap); err != nil {
			log.Fatal("IN JSON ", err)
		}
		properties := responseMap["definitions"].(map[string]interface{})[vcName].(map[string]interface{})["properties"].(map[string]interface{})
		for k, _ := range properties {
			if lastColumnIndex < data.Head[k] {
				lastColumnIndex = data.Head[k]
			}
			jsonBody[k] = data.Row[data.Head[k]]
		}
		bytes, _ := json.Marshal(jsonBody)
		req, err := http.NewRequest("POST", config.Config.Registry.BaseUrl+"api/v1/"+vcName, strings.NewReader(string(bytes)))
		req.Header.Set("Authorization", header["Authorization"][0])
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Fatal("Error : %v", err)
		}
		currRow = data.Row
		if res.StatusCode == 200 {
			totalSuccess += 1
		} else {
			totalErrors += 1
			resBody, _ := ioutil.ReadAll(res.Body)
			data.Head["Errors"] = lastColumnIndex + 1
			currRow = append(currRow, string(resBody))
		}
		rows = append(rows, currRow)
	}
	log.Printf("Rows : %v", rows)
	return totalSuccess, totalErrors, rows
}
