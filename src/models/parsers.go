package models

import (
	"encoding/json"
	"fmt"
)

type Parser interface {
	ParseRoutes(settings map[string]interface{}) []RouteConfig
}

type GetRoutesParser struct {
	Parser
}

func (p GetRoutesParser) ParseRoutes(settings map[string]interface{}) []RouteConfig {
	routeConfigs := make([]RouteConfig, 0)
	if setting, ok := settings["GET"]; ok {
		getRouteSettings := setting.([]interface{})
		for _, getRouteSetting := range getRouteSettings {
			route := getRouteSetting.(map[string]interface{})
			routeConfig := RouteConfig{
				Method: "GET",
			}
			if path, ok := route["path"]; ok {
				routeConfig.Path = path.(string)
			} else {
				continue
			}
			// parse body
			if bodySettings, ok := route["body"]; ok {
				switch v := bodySettings.(type) {
				case string:
					routeConfig.Body = v
					break
				case map[string]interface{}:
					content, err := json.Marshal(v)
					if err == nil {
						routeConfig.Body = string(content)
					} else {
						fmt.Println("Error marshaling JSON:", err)
						routeConfig.Body = ""
					}
				default:
					routeConfig.Body = ""
					fmt.Println("Unknown type:", v)
				}
			} else {
				routeConfig.Body = ""
			}

			// parse status code
			if statusCode, ok := route["status"]; ok {
				routeConfig.StatusCode = int(statusCode.(float64))
			} else {
				routeConfig.StatusCode = 200
			}

			// parse headers
			if headers, ok := route["headers"]; ok {
				temp := headers.(map[string]interface{})
				routeConfig.Headers = make(map[string]string)
				for key, value := range temp {
					routeConfig.Headers[key] = value.(string)
				}
			} else {
				routeConfig.Headers = nil
			}

			routeConfigs = append(routeConfigs, routeConfig)
		}
	}
	return routeConfigs
}
