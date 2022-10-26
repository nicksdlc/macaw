package types

import "time"

func Date(parameters []string) Type {
	return &DateNow{}
}

type DateNow struct {
}

func (dn *DateNow) Value() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.45Z")
}
