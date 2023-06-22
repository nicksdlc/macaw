package config

// Body - message body configuration
type Body struct {
	File   string
	String string
}

// Options - message options configuration
type Options struct {
	Quantity int
	Delay    int
}
