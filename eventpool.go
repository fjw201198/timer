package timer

type eventPool struct {
	events   []*sortedNode
	freeList []int
}

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

func (ep *eventPool) clear() {
	for i := 0; i < TIMER_QUEUE_SIZE; i++ {
		ep.events[i].Data.Clear()
		ep.events[i].LChild = nil
		ep.events[i].RChild = nil
	}
}

func (ep *eventPool) alloc() (*sortedNode, error) {
	var flen = len(ep.freeList)
	if flen == 0 {
		return nil, ErrInfo[E_QFull]
	}
	var idx = ep.freeList[flen-1]
	ep.freeList = ep.freeList[0 : flen-1]
	sn := ep.events[idx]
	sn.LChild = nil
	sn.RChild = nil
	sn.id = idx
	return sn, nil
}

func (ep *eventPool) free(sn *sortedNode) {
	sn.Data.Clear()
	ep.freeList = append(ep.freeList, sn.id)
}
