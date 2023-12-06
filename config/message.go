package config

// Body - message body configuration
type Body struct {
	File   []string
	String []string
}

// Options - message options configuration
type Options struct {
	Quantity    int
	Delay       string
	RandomDelay RandomDelay
}

// Delay - message delay configuration
type RandomDelay struct {
	Min string
	Max string
}
