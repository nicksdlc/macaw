package data

// Request represents request that will arrive
type Request struct {
	Version          string
	EventTimeUtc     string
	SystemID         int
	PolicyID         int
	DateTimeRangeUtc DateTimeRange
	RequiredFields   []string
}

// DateTimeRange to search policy in
type DateTimeRange struct {
	From string
	To   string
}
