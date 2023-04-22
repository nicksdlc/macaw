package config

// HTTP represents HTTP connector configuration
type HTTP struct {
	Serve  Serve
	Remote Remote
}

// Remote represents remote HTTP server configuration
type Remote struct {
	Host string
	Port uint16
}

// Serve represents HTTP response server configuration
type Serve struct {
	Host string
	Port uint16
}
