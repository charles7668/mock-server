package src

import (
	"encoding/json"
	"io"
	"log"
	"mock-server/src/models"
	"net/http"
)

func getMethodHandler(routeConfig models.RouteConfig) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		for key, value := range routeConfig.Headers {
			writer.Header().Set(key, value)
		}
		writer.WriteHeader(routeConfig.StatusCode)

		_, err := writer.Write([]byte(routeConfig.Body))
		if err != nil {
			return
		}
	}
}

func postMethodHandler(routeConfig models.RouteConfig) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("POST request received")
		// set default
		overrides := routeConfig.Overrides
		// get content type
		contentType := request.Header.Get("Content-Type")
		log.Println("Content-Type:", contentType)
		var requestBody interface{}
		isJson := true
		switch contentType {
		case "application/json":
			err := json.NewDecoder(request.Body).Decode(&requestBody)
			if err != nil {
				http.Error(writer, "Invalid request body", http.StatusBadRequest)
				return
			}
			break
		default:
			isJson = false
			bodyBytes, err := io.ReadAll(request.Body)
			if err != nil {
				http.Error(writer, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			requestBody = string(bodyBytes)
		}
		// Read and parse the request body
		foundOverride := false
		overrideConfig := make(map[string]interface{})
		for _, override := range overrides {
			if temp, ok := override["verify"]; ok {
				// verify
				match := true
				if isJson {
					verify := temp.(map[string]interface{})
					for key, value := range verify {
						if requestBody.(map[string]interface{})[key] != value {
							match = false
							break
						}
					}
				} else {
					verify := temp.(string)
					if requestBody != verify {
						match = false
					}
				}
				if match {
					foundOverride = true
					overrideConfig = override
				}
			}
		}
		if overrideHeaders, ok := overrideConfig["headers"]; ok && foundOverride {
			log.Println("Overriding headers found")
			for key, value := range overrideHeaders.(map[string]interface{}) {
				log.Println("Setting header:", key, value)
				writer.Header().Set(key, value.(string))
			}
		} else {
			for key, value := range routeConfig.Headers {
				log.Println("Setting header:", key, value)
				writer.Header().Set(key, value)
			}
		}
		if overrideStatus, ok := overrideConfig["status"]; ok && foundOverride {
			log.Println("Overriding status code found")
			log.Println("Setting status code:", int(overrideStatus.(float64)))
			writer.WriteHeader(int(overrideStatus.(float64)))
		} else {
			log.Println("Setting status code:", routeConfig.StatusCode)
			writer.WriteHeader(routeConfig.StatusCode)
		}

		if overrideBody, ok := overrideConfig["body"]; ok && foundOverride {
			log.Println("Overriding body found")
			var body string
			switch v := overrideBody.(type) {
			case string:
				body = overrideBody.(string)
				break
			case map[string]interface{}:
				content, err := json.Marshal(v)
				if err == nil {
					body = string(content)
				} else {
					body = ""
				}
			case []interface{}:
				content, err := json.Marshal(v)
				if err == nil {
					body = string(content)
				} else {
					body = ""
				}
			default:
				body = ""
			}
			_, _ = writer.Write([]byte(body))
		} else {
			_, _ = writer.Write([]byte(routeConfig.Body))
		}
		return
	}
}
