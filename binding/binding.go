package binding

import "net/http"

const (
	MIMEJSON     = "application/json"
	MIMEPOSTForm = "application/x-www-form-urlencoded"
)

type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

var (
	JSON = jsonBinding{}
)

func New(method, contentType string) Binding {

	switch contentType {

	default:
		return JSON
	}
}
