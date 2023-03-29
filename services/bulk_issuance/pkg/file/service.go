package file_service

import (
	"encoding/csv"
	"log"
	"os"
)

func CreateFile(fileName string, values [][]string) {
	csvFile, err := os.Create(fileName)
	defer csvFile.Close()
	if err != nil {
		log.Printf("Error creating csv file : %v", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()
	for _, value := range values {
		csvwriter.Write(value)
	}
	csvwriter.Flush()
	csvFile.Close()
}
