package utils

import (
	"bulk_issuance/config"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func GetSchemaProperties(key string) map[string]interface{} {
	resp, err := http.Get(config.Config.Registry.BaseUrl + "api/docs/swagger.json")
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var responseMap map[string]interface{}
	if err = json.Unmarshal(body, &responseMap); err != nil {
		log.Fatal("IN JSON ", err)
	}
	return responseMap["definitions"].(map[string]interface{})[key].(map[string]interface{})["properties"].(map[string]interface{})
}
