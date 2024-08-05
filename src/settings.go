package src

import (
	"encoding/json"
	"log"
	"os"
)

func LoadSettings(file string) map[string]interface{} {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var settings = make(map[string]interface{})
	err = json.Unmarshal(fileContent, &settings)
	if err != nil {
		log.Fatal(err)
	}
	return settings
}
