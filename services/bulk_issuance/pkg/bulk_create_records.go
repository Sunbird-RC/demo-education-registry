package pkg

import (
	"bulk_issuance/config"
	"bulk_issuance/db"
	"bulk_issuance/swagger_gen/models"
	"bulk_issuance/swagger_gen/restapi/operations/upload_and_create_records"
	"bulk_issuance/utils"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
)

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
	totalSuccess, totalErrors, rows := processDataFromCSV(&data, params.HTTPRequest.Header, params.VCName)
	successFailureCount := map[string]int{
		"success":   totalSuccess,
		"error":     totalErrors,
		"totalRows": totalSuccess + totalErrors,
	}
	response := upload_and_create_records.NewPostV1UploadFilesVCNameOK()
	response.SetPayload(successFailureCount)
	_, fileHeader, _ := params.HTTPRequest.FormFile("file")
	fileName := fileHeader.Filename
	err := addEntryForDbUploadToDatabase(fileName, totalSuccess, totalErrors, principal)
	logErrorIfAny("DBUpload", err)
	addEntryForDbFilesToDatabase(rows, fileName, data)
	logErrorIfAny("DBFilesUpload", err)
	return response
}

func logErrorIfAny(entity string, err error) {
	if err != nil {
		log.Printf("Error in adding %v to database : %v", entity, err)
	}
}

func addEntryForDbFilesToDatabase(rows [][]string, fileName string, data Scanner) error {
	bytes, _ := json.Marshal(rows)
	fileUpload := db.DBFilesUpload{
		Filename: fileName,
		Headers:  getHeaders(data.Head),
		RowData:  bytes,
	}
	return db.CreateDBFilesUpload(&fileUpload)
}

func addEntryForDbUploadToDatabase(fileName string, totalSuccess int, totalErrors int, principal *models.JWTClaimBody) error {
	dbUpload := db.DBFileUpload{
		Filename:     fileName,
		TotalRecords: totalSuccess + totalErrors,
		UserID:       principal.PreferredUsername,
		Date:         time.Now().Format("2006-01-02"),
	}
	return db.CreateDBFileUpload(&dbUpload)
}

func processDataFromCSV(data *Scanner, header http.Header, vcName string) (int, int, [][]string) {
	var (
		totalSuccess int = 0
		totalErrors  int = 0
	)
	rows := make([][]string, 0)
	for data.Scan() {
		jsonBody := make(map[string]interface{})
		properties := utils.GetSchemaProperties(vcName)
		bytes := createReqBody(properties, jsonBody, data)
		res, err := createSingleRecord(vcName, bytes, header)
		if err != nil {
			log.Printf("Error in creating a record : %v", err)
		}
		currRow := make([]string, 0)
		currRow = data.Row
		if res.StatusCode != 200 {
			lastColIndex := len(properties) + 1
			currRow = appendErrorsToCurrentRow(res, data, lastColIndex, currRow)
			totalErrors += 1
		} else {
			totalSuccess += 1
		}
		rows = append(rows, currRow)
	}
	return totalSuccess, totalErrors, rows
}

func appendErrorsToCurrentRow(res *http.Response, data *Scanner, lastColIndex int, currRow []string) []string {
	resBody, _ := ioutil.ReadAll(res.Body)
	data.Head["Errors"] = lastColIndex
	currRow = append(currRow, string(resBody))
	return currRow
}

func createSingleRecord(vcName string, bytes []byte, header http.Header) (*http.Response, error) {
	methodName := "POST"
	req, err := http.NewRequest(methodName, config.Config.Registry.BaseUrl+"api/v1/"+vcName, strings.NewReader(string(bytes)))
	if err != nil {
		log.Printf("Error in creating request : %v", err)
	}
	req.Header.Set("Authorization", header["Authorization"][0])
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

func createReqBody(properties map[string]interface{}, jsonBody map[string]interface{}, data *Scanner) []byte {
	for k, _ := range properties {
		jsonBody[k] = data.Row[data.Head[k]]
	}
	bytes, _ := json.Marshal(jsonBody)
	return bytes
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
