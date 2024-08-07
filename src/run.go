package src

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
		models.PostRoutesParser{},
	}
	for _, parser := range parsers {
		temp := parser.ParseRoutes(settings)
		routeConfigs = append(routeConfigs, temp...)
	}
	for _, routeConfig := range routeConfigs {
		println("routed : ", routeConfig.Method, routeConfig.Path)
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	for _, routeConfig := range routeConfigs {
		var handler http.HandlerFunc
		switch routeConfig.Method {
		case "GET":
			handler = getMethodHandler(routeConfig)
			break
		case "POST":
			handler = postMethodHandler(routeConfig)
			break
		default:
			log.Fatal("Unknown method:", routeConfig.Method)
		}
		router.Method(
			routeConfig.Method,
			routeConfig.Path,
			handler,
		)
	}
	log.Println("Starting server on port", serverOptions.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(serverOptions.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}
