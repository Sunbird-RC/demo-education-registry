package utils

import (
	"bulk_issuance/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	log "github.com/sirupsen/logrus"
)

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func GetSchemaProperties(key string) []string {
	resp, err := http.Get(config.Config.Registry.BaseUrl + "api/docs/swagger.json")
	LogErrorIfAny("Error creating a get request for %v : %v", err, config.Config.Registry.BaseUrl+"api/docs/swagger.json")
	body, _ := ioutil.ReadAll(resp.Body)
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	LogErrorIfAny("Error creating request body for %v : %v", err, config.Config.Registry.BaseUrl+"api/docs/swagger.json")
	schemaProperties := responseMap["definitions"].(map[string]interface{})[key].(map[string]interface{})["properties"].(map[string]interface{})
	properties := make([]string, 0)
	for k := range schemaProperties {
		properties = append(properties, k)
	}
	sort.Strings(properties)
	return properties
}

func LogErrorIfAny(description string, err error, data ...interface{}) {
	if err != nil {
		log.Errorf(description, data, err)
	}
}
