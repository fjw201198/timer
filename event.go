package timer

type TimeoutCallback func()

type eventObject struct {
	// who  interface{};
	fn   TimeoutCallback
	tick int
}

func (e *eventObject) isValid() bool {
	return e.tick == INVALID_TICK
}

func (e *eventObject) Less(rh TreeDataType) bool {
	rh2 := rh.(*eventObject)
	return e.tick < rh2.tick
}

func NewEventObject(p ...interface{}) TreeDataType {
	plen := len(p)
	var fn TimeoutCallback = nil
	var tick int = INVALID_TICK
	if p != nil {
		if plen > 0 {
			fn = p[0].(TimeoutCallback)
		}
		if plen > 1 {
			tick = p[1].(int)
		}
	}
	return &eventObject{fn, tick}
}

func (e *eventObject) Clear() {
	e.fn = nil
	e.tick = INVALID_TICK
}

func (e *eventObject) SetParam(p ...interface{}) {
	plen := len(p)
	if plen > 0 {
		e.fn = p[0].(TimeoutCallback)
	}
	if plen > 1 {
		e.tick = p[1].(int)
	}
}

func (e *eventObject) DataLess(rh interface{}) bool {
	tick := rh.(int)
	return e.tick < tick
}
