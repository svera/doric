package columns

const (
	ended = iota
	removed
)

type Event struct {
	name int
	data interface{}
}

func NewEvent(name int, data interface{}) *Event {
	return &Event{name, data}
}

func (e *Event) Name() int {
	return e.name
}

func (e *Event) Data() interface{} {
	return e.data
}
