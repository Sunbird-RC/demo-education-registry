package services

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func CallPOST(url string, payload string) map[string]interface{} {
	log.Info("Call POST API triggered")
	method := "POST"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		log.Error("Failed creating new client request", err)
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Error("POST API failed with", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Error while closing the request body", err)
		}
	}(res.Body)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Error while reading the request body", err)
		return nil
	}
	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Error("Error while unmarshalling the response data", err)
		return nil
	}
	log.Info("Call POST API completed")
	return data
}
