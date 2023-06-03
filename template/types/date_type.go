package types

import (
	"time"
)

var dateTypes = make(map[string]func([]string) Type)

func init() {
	dateTypes["now"] = dateNow
	dateTypes["increment"] = dateIncrement
}

func Date(parameters []string) Type {
	if len(parameters) == 0 {
		return dateNow([]string{})
	}
	return dateTypes[parameters[0]](parameters[1:])
}

type DateNow struct {
}

func dateNow([]string) Type {
	return &DateNow{}
}

func (dn *DateNow) Value() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.45Z")
}

type DateIncrement struct {
	startingDate time.Time
	increment    time.Duration
}

func dateIncrement(parameters []string) Type {
	return &DateIncrement{
		startingDate: time.Date(2021, time.January, 10, 11, 0, 0, 0, time.UTC).UTC(),
		increment:    time.Duration(1 * time.Second),
	}
}

func (di *DateIncrement) Value() string {
	di.startingDate = di.startingDate.Add(di.increment)
	return di.startingDate.Format("2006-01-02T15:04:05.45Z")
}
