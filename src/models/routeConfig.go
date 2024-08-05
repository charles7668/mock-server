package models

type RouteConfig struct {
	Method     string
	Path       string
	Body       string
	StatusCode int
	Headers    map[string]string
}
