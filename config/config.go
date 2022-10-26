package config

type Configuration struct {
	Mock             string
	Rabbit           RabbitMQ
	ResponseTemplate string
}

type RabbitMQ struct {
	Host          string
	Port          string
	User          string
	Password      string
	RequestQueue  string
	ResponseQueue string
}
