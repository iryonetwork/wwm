package status

//go:generate sh ../bin/mockgen.sh status Component $GOFILE

type Value string

// Status values
const (
	OK      Value = "ok"
	Warning       = "warning"
	Error         = "error"
)

type Response struct {
	Status     Value                `json:"status"`
	Msg        string               `json:"description,omitempty"`
	Components map[string]*Response `json:"components,omitempty"`
}

func (s Value) Int() int8 {
	switch s {
	case OK:
		return 0
	case Warning:
		return 1
	case Error:
		return 2
	}

	// invalid status value
	return -1
}

// Component is an interface for a component that exposes its status
type Component interface {
	Status() *Response
}
