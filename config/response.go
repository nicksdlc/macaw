package config

// Response configuration
type Response struct {
	To              string
	ResponseRequest ResponseRequest `json:"request" yaml:"request" mapstructure:"request"`
	Body            Body
	Options         Options
}

// ResponseRequest configuration
type ResponseRequest struct {
	To       string
	Matchers []Matcher `json:"matchers" yaml:"matchers" mapstructure:"matchers"`
}
