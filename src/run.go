package src

import (
	"github.com/go-chi/chi/v5"
	"log"
	"mock-server/src/models"
	"net/http"
	"strconv"
)

func Start(options models.LaunchOptions) {
	startServer(LoadSettings(options.File))
}

func startServer(settings map[string]interface{}) {
	routeConfigs := make([]models.RouteConfig, 0)
	serverOptions := models.NewServerOptions()
	serverOptions.Parse(settings)
	parsers := []models.Parser{
		models.GetRoutesParser{},
	}
	for _, parser := range parsers {
		routeConfigs = append(routeConfigs, parser.ParseRoutes(settings)...)
	}
	for _, routeConfig := range routeConfigs {
		println("routed : ", routeConfig.Method, routeConfig.Path)
	}

	router := chi.NewRouter()
	for _, routeConfig := range routeConfigs {
		handler := func(writer http.ResponseWriter, request *http.Request) {
			for key, value := range routeConfig.Headers {
				writer.Header().Set(key, value)
			}
			writer.WriteHeader(routeConfig.StatusCode)

			_, err := writer.Write([]byte(routeConfig.Body))
			if err != nil {
				return
			}
		}
		router.Method(
			routeConfig.Method,
			routeConfig.Path,
			http.HandlerFunc(handler),
		)
	}
	err := http.ListenAndServe(":"+strconv.Itoa(serverOptions.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}
