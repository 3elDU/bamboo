package event

type Type int
type Event interface {
	Type() Type
	Args() interface{}
}

type CustomEvent struct {
	t    Type
	args interface{}
}

func (e CustomEvent) Type() Type {
	return e.t
}
func (e CustomEvent) Args() interface{} {
	return e.args
}

func NewEvent(t Type, args interface{}) CustomEvent {
	return CustomEvent{
		t:    t,
		args: args,
	}
}

var eventList []Event

func FireEvent(event Event) {
	eventList = append(eventList, event)
}

func GetEvents() []Event {
	events := make([]Event, len(eventList))
	copy(events, eventList)
	eventList = []Event{}
	return events
}
