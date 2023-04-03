package pkg

import (
	file_service "bulk_issuance/pkg/file"
	"bulk_issuance/swagger_gen/models"
	"bulk_issuance/swagger_gen/restapi/operations/sample_template"
	"bulk_issuance/utils"
	"os"

	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
)

func downloadSampleFile(params sample_template.GetV1BulkSampleSchemaNameParams, principal *models.JWTClaimBody) middleware.Responder {
	log.Info("Downloading sample file")
	response := sample_template.GetV1BulkSampleSchemaNameOK{}
	schemaProperties, sampleValues, err := utils.GetSchemaPropertiesAndSampleValues(params.SchemaName)
	if err != nil {
		response := sample_template.NewGetV1BulkSampleSchemaNameNotFound()
		response.SetPayload(err.Error())
		return response
	}
	csvData := make([][]string, 0)
	csvData = append(csvData, schemaProperties)
	csvData = append(csvData, sampleValues)
	log.Infof("Headers for schema %v : %v", params.SchemaName, csvData)
	file_service.CreateFile(params.SchemaName+".csv", csvData)
	f, err := os.Open(params.SchemaName + ".csv")
	utils.LogErrorIfAny("Error while opening file : %v", err)
	response.WithContentDisposition("attachment; filename=\"" + params.SchemaName + ".csv\"").WithPayload(f)
	return &response
}
