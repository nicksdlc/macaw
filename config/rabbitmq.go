package config

// RabbitMQ configuration of RabbitMQ
type RabbitMQ struct {
	Host            string
	Port            string
	User            string
	Password        string
	Exchanges       []Exchange
	Queues          []Queue
	ConnectionRetry Retry
}

type Queue struct {
	Name string
	Args map[string]interface{}
}

type Exchange struct {
	Name       string
	RoutingKey string
}
