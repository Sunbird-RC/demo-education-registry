package file_service

import (
	"bulk_issuance/utils"
	"encoding/csv"
	"os"
)

func CreateFile(fileName string, values [][]string) {
	csvFile, err := os.Create(fileName)
	defer csvFile.Close()
	utils.LogErrorIfAny("Error creating csv file with name %v : %v", err, fileName)
	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()
	for _, value := range values {
		csvwriter.Write(value)
	}
	csvwriter.Flush()
	err = csvFile.Close()
	utils.LogErrorIfAny("Error while closing the file with name %v : %v", err, fileName)
}
