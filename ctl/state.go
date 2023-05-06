package ctl

type Health struct {
	Listener Listener `json:"listener"`
}

type Listener struct {
	Listener  string `json:"listener"`
	ListensOn string `json:"listens_on"`
	State     string `json:"state"`
	Reason    string `json:"reason"`
}
