package pkg

import (
	"bulk_issuance/config"
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
	"strings"
	"time"

	"log"

	"github.com/go-openapi/runtime/middleware"
)

func SetupHandlers(api *operations.BulkIssuanceAPI) {
	api.SampleTemplateGetV1BulkSampleSchemaNameHandler = sample_template.GetV1BulkSampleSchemaNameHandlerFunc(downloadSampleFile)
	api.UploadedFilesGetV1BulkUploadedFilesHandler = uploaded_files.GetV1BulkUploadedFilesHandlerFunc(listFiles)
	api.DownloadFileReportGetV1DownloadFileNameHandler = download_file_report.GetV1DownloadFileNameHandlerFunc(downloadReportFile)
	api.UploadAndCreateRecordsPostV1UploadFilesVCNameHandler = upload_and_create_records.PostV1UploadFilesVCNameHandlerFunc(createRecords)
}

func downloadSampleFile(params sample_template.GetV1BulkSampleSchemaNameParams) middleware.Responder {
	response := sample_template.GetV1BulkSampleSchemaNameOK{}
	headers := [][]string{
		{"Record Name", "Email", "Phone", "Date of Birth", "Validity"},
	}
	csvFile, err := os.Create(params.SchemaName + ".csv")
	if err != nil {
		log.Printf("Error creating csv file : %v", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	for _, header_row := range headers {
		log.Print(header_row)
		csvwriter.Write(header_row)
	}
	csvwriter.Flush()
	response.SetPayload(csvFile)
	return &response
}

func listFiles(params uploaded_files.GetV1BulkUploadedFilesParams) middleware.Responder {
	response := uploaded_files.GetV1BulkUploadedFilesOK{}
	uploadedFiles := make([]map[string]interface{}, 0)
	uploadedFiles = append(uploadedFiles, map[string]interface{}{
		"name":              "file1",
		"uploadedBy":        "temp_user1",
		"number_of_records": "152",
		"date":              time.Now().Format("2006-01-02"),
	})
	uploadedFiles = append(uploadedFiles, map[string]interface{}{
		"name":              "file2",
		"uploadedBy":        "temp_user2",
		"number_of_records": "152",
		"date":              time.Now().Format("2006-01-02"),
	})
	uploadedFiles = append(uploadedFiles, map[string]interface{}{
		"name":              "file3",
		"uploadedBy":        "temp_user3",
		"number_of_records": "152",
		"date":              time.Now().Format("2006-01-02"),
	})
	uploadedFiles = append(uploadedFiles, map[string]interface{}{
		"name":              "file4",
		"uploadedBy":        "temp_user4",
		"number_of_records": "152",
		"date":              time.Now().Format("2006-01-02"),
	})
	uploadedFiles = append(uploadedFiles, map[string]interface{}{
		"name":              "file5",
		"uploadedBy":        "temp_user5",
		"number_of_records": "152",
		"date":              time.Now().Format("2006-01-02"),
	})
	response.SetPayload(uploadedFiles)
	return &response
}

func downloadReportFile(params download_file_report.GetV1DownloadFileNameParams) middleware.Responder {
	return download_file_report.NewGetV1DownloadFileNameOK()
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
	totalSuccess, totalErrors := Process(&data, params.HTTPRequest.Header, params.VCName)
	log.Printf("Principal : %v", principal.PreferredUsername)
	successFailureCount := map[string]int{
		"success":   totalSuccess,
		"error":     totalErrors,
		"totalRows": totalSuccess + totalErrors,
	}
	response := upload_and_create_records.NewPostV1UploadFilesVCNameOK()
	response.SetPayload(successFailureCount)
	return response
}

func Process(data *Scanner, header http.Header, vcName string) (int, int) {
	var (
		totalSuccess int = 0
		totalErrors  int = 0
	)
	for data.Scan() {
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
		if res.StatusCode == 200 {
			totalSuccess += 1
		} else {
			totalErrors += 1
		}
		log.Printf("Response Status : %v", res)
	}
	return totalSuccess, totalErrors
}
