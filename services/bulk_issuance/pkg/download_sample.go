package pkg

import (
	file_service "bulk_issuance/pkg/file"
	"bulk_issuance/swagger_gen/restapi/operations/sample_template"
	"bulk_issuance/utils"
	"os"

	"github.com/go-openapi/runtime/middleware"
)

func downloadSampleFile(params sample_template.GetV1BulkSampleSchemaNameParams) middleware.Responder {
	response := sample_template.GetV1BulkSampleSchemaNameOK{}
	schemaHeader := getFieldHeaders(params.SchemaName)
	headers := make([][]string, 0)
	headers = append(headers, schemaHeader)
	file_service.CreateFile(params.SchemaName+".csv", headers)
	f, _ := os.Open(params.SchemaName + ".csv")
	response.WithContentDisposition("attachment; filename=\"" + params.SchemaName + ".csv\"").WithPayload(f)
	return &response
}

func getFieldHeaders(schemaName string) []string {
	properties := utils.GetSchemaProperties(schemaName)
	headers := make([]string, 0)
	for k := range properties {
		headers = append(headers, k)
	}
	return headers
}
