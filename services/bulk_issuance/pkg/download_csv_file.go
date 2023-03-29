package pkg

import (
	"bulk_issuance/db"
	file_service "bulk_issuance/pkg/file"
	"bulk_issuance/swagger_gen/restapi/operations/download_file_report"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

func downloadReportFile(params download_file_report.GetV1DownloadIDParams) middleware.Responder {
	response := download_file_report.GetV1DownloadFileNameOK{}
	file := db.GetDBFilesUpload(int(params.ID))
	var data [][]string
	if err := json.Unmarshal(file.RowData, &data); err != nil {
		log.Fatal("Error : %v", err)
	}
	data = append([][]string{strings.Split(file.Headers, ",")}, data...)
	file_service.CreateFile(file.Filename, data)
	f, _ := os.Open(file.Filename)
	response.WithContentDisposition("attachment; filename=\"" + file.Filename + "\"").WithPayload(f)
	return &response
}
