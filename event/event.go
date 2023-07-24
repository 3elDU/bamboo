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

	for _, subscriber := range subscribers {
		if subscriber.eventType == event.Type() {
			subscriber.callback(event.Args())
		}
	}
}

func GetEvents() []Event {
	events := make([]Event, len(eventList))
	copy(events, eventList)
	eventList = []Event{}
	return events
}

type eventSubscriber struct {
	eventType Type
	callback  func(args interface{})
}

var subscribers []eventSubscriber

// Register a custom callback for a specific event, that will be called when the event is fired.
func Subscribe(eventType Type, callback func(args interface{})) {
	subscribers = append(subscribers, eventSubscriber{
		eventType: eventType,
		callback:  callback,
	})
}
