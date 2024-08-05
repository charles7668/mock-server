package models

type ServerOptions struct {
	Port int
}

// NewServerOptions creates a new ServerOptions with default values.
func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		Port: 3000,
	}
}

// Parse parses the settings map and returns a ServerOptions struct.
func (options *ServerOptions) Parse(settings map[string]interface{}) {
	if port, ok := settings["port"]; ok {
		options.Port = int(port.(float64))
	}
}
