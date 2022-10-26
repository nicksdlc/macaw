package data

// Response is a response representation
type Response struct {
	Version            string
	EventTimeUtc       string
	Status             string
	Reason             string
	PolicyID           int
	PolicyEventTimeUtc string
	BulkID             int
	BulksCount         int
	// Data               ResponseData
}

// ResponseData is the list of fields and list of values
type ResponseData struct {
	Fields []string
	Values [][]string
}
