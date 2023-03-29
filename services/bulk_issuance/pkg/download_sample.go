package pkg

import (
	file_service "bulk_issuance/pkg/file"
	"bulk_issuance/swagger_gen/restapi/operations/sample_template"
	"bulk_issuance/utils"
	"os"

	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
)

func downloadSampleFile(params sample_template.GetV1BulkSampleSchemaNameParams) middleware.Responder {
	log.Info("Downloading sample file")
	response := sample_template.GetV1BulkSampleSchemaNameOK{}
	schemaProperties := utils.GetSchemaProperties(params.SchemaName)
	headers := make([][]string, 0)
	headers = append(headers, schemaProperties)
	log.Infof("Headers for schema %v : %v", params.SchemaName, headers)
	file_service.CreateFile(params.SchemaName+".csv", headers)
	f, err := os.Open(params.SchemaName + ".csv")
	utils.LogErrorIfAny("Error while opening file : %v", err)
	response.WithContentDisposition("attachment; filename=\"" + params.SchemaName + ".csv\"").WithPayload(f)
	return &response
}
