package connectors

type Connector interface {
	Close() error
	Post(body string) error
}
