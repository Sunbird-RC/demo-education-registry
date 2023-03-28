package pkg

import (
	"bulk_issuance/swagger_gen/restapi/operations"
	"bulk_issuance/swagger_gen/restapi/operations/download_file_report"
	"bulk_issuance/swagger_gen/restapi/operations/sample_template"
	"bulk_issuance/swagger_gen/restapi/operations/upload_and_create_records"
	"bulk_issuance/swagger_gen/restapi/operations/uploaded_files"
	"encoding/csv"
	"os"
	"time"

	"log"

	"github.com/go-openapi/runtime/middleware"
)

func SetupHandlers(api *operations.BulkIssuanceAPI) {
	api.SampleTemplateGetV1BulkSampleSchemaNameHandler = sample_template.GetV1BulkSampleSchemaNameHandlerFunc(downloadSampleFile)
	api.UploadedFilesGetV1BulkUploadedFilesHandler = uploaded_files.GetV1BulkUploadedFilesHandlerFunc(listFiles)
	api.DownloadFileReportGetV1DownloadFileNameHandler = download_file_report.GetV1DownloadFileNameHandlerFunc(downloadReportFile)
	api.UploadAndCreateRecordsPostV1UploadFilesHandler = upload_and_create_records.PostV1UploadFilesHandlerFunc(createRecords)
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

func createRecords(params upload_and_create_records.PostV1UploadFilesParams) middleware.Responder {
	return upload_and_create_records.NewPostV1UploadFilesOK()
}
