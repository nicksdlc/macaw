package config

// Response configuration
type Response struct {
	Type            string
	FromOpenAPI     string
	Alias           string
	To              string
	ResponseRequest ResponseRequest `json:"request" yaml:"request" mapstructure:"request"`
	Body            *Body
	Options         *Options
}

// ResponseRequest configuration
type ResponseRequest struct {
	To       string
	Match    string    `json:"match" yaml:"match" mapstructure:"match"`
	Matchers []Matcher `json:"matchers" yaml:"matchers" mapstructure:"matchers"`
}
