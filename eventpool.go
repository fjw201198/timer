package timer

type eventPool struct {
	events   []*sortedNode
	freeList []int
}

// newEventPool create pool instance.
func newEventPool(newDataFn NewDataFuncType) *eventPool {
	ep := new(eventPool)
	ep.events = make([]*sortedNode, 0, TIMER_QUEUE_SIZE)
	ep.freeList = make([]int, 0, TIMER_QUEUE_SIZE)
	for i := 0; i < TIMER_QUEUE_SIZE; i++ {
		sn := &sortedNode{newDataFn(), nil, nil, 0}
		ep.events = append(ep.events, sn)
		ep.freeList = append(ep.freeList, 0+i)
	}
	return ep
}

// clear() will clear all events
func (ep *eventPool) clear() {
	if len(ep.freeList) > 0 {
		ep.freeList = ep.freeList[0:0]
	}
	for i := 0; i < TIMER_QUEUE_SIZE; i++ {
		ep.events[i].Data.Clear()
		ep.events[i].LChild = nil
		ep.events[i].RChild = nil
		ep.freeList = append(ep.freeList, i)
	}
}

// alloc() will get an events from the pool, if the pool is full,
// we'll new one, but it will not put to the pool
func (ep *eventPool) alloc() (*sortedNode, error) {
	var flen = len(ep.freeList)
	if flen == 0 {
		sn := new(sortedNode)
		sn.id = INVALID_TICK
		sn.LChild = nil
		sn.RChild = nil
		return sn, nil
		// return nil, ErrInfo[E_QFull]
	}
	var idx = ep.freeList[flen-1]
	ep.freeList = ep.freeList[0 : flen-1]
	sn := ep.events[idx]
	sn.LChild = nil
	sn.RChild = nil
	sn.id = idx
	return sn, nil
}

// free() will free the event, if the event owned by the pool,
// it will be append to the free list, else the event will be free by gc.
func (ep *eventPool) free(sn *sortedNode) {
	sn.Data.Clear()
	if sn.id == INVALID_TICK {
		sn = nil
		return
	}
	ep.freeList = append(ep.freeList, sn.id)
}
