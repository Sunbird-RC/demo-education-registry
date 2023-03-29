package pkg

import (
	"bulk_issuance/db"
	"bulk_issuance/swagger_gen/restapi/operations/uploaded_files"

	"github.com/go-openapi/runtime/middleware"
)

func listFiles(params uploaded_files.GetV1BulkUploadedFilesParams) middleware.Responder {
	response := uploaded_files.GetV1BulkUploadedFilesOK{}
	files, err := db.GetAllUploadedFilesData()
	if err != nil {
		return &response
	}
	response.SetPayload(files)
	return &response
}
