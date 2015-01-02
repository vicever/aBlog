package sys

var Event *eventCaller = newEventCaller()

type EventMsg struct {
	Name  string
	Value interface{}
}

type eventCaller struct {
	handlers    map[string][]func(EventMsg)
	msgChan     chan EventMsg
	EnableAsync bool
}

func newEventCaller() *eventCaller {
	caller := &eventCaller{
		handlers:    make(map[string][]func(EventMsg)),
		msgChan:     make(chan EventMsg),
		EnableAsync: true,
	}
	go caller.init()
	return caller
}

func (caller *eventCaller) init() {
	for {
		msg := <-caller.msgChan
		caller.callMessage(msg)
	}
}

// event on
func (caller *eventCaller) On(name string, fn func(EventMsg)) {
	if len(caller.handlers[name]) == 0 {
		caller.handlers[name] = []func(EventMsg){fn}
		return
	}
	caller.handlers[name] = append(caller.handlers[name], fn)
}

func (caller *eventCaller) callMessage(msg EventMsg) {
	handlers := caller.handlers[msg.Name]
	if len(handlers) == 0 {
		return
	}
	for _, fn := range handlers {
		fn(msg)
	}
}

// run event sync
func (caller *eventCaller) CallSync(name string, value interface{}) {
	msg := EventMsg{name, value}
	caller.callMessage(msg)
}

// trigger means broadcast event message, not run sync
func (caller *eventCaller) Call(name string, value interface{}) {
	if !caller.EnableAsync {
		caller.CallSync(name, value)
		return
	}
	msg := EventMsg{name, value}
	caller.msgChan <- msg
}
