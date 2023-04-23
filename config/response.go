package config

// Response configuration
type Response struct {
	ResponseRequest ResponseRequest `json:"request" yaml:"request" mapstructure:"request"`
	File            string
	String          string
	Amount          int
}

// ResponseRequest configuration
type ResponseRequest struct {
	To    string
	Field Field
}

// Field configuration
type Field struct {
	Name  string
	Value string
}
