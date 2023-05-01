package config

// Body - message body configuration
type Body struct {
	File   string
	String string
}

// Matchers - message matcher configuration
type Matchers struct {
	Field    FieldMatcher        `json:"field" yaml:"field" mapstructure:"field"`
	Contains BodyContainsMatcher `json:"contains" yaml:"contains" mapstructure:"contains"`
}

// Field configuration
type FieldMatcher struct {
	Name  string
	Value string
}

// BodyContainsMatcher configuration
type BodyContainsMatcher struct {
	Value string
}

// Options - message options configuration
type Options struct {
	Quantity int
	Delay    int
}
