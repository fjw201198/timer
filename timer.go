// package timer defined some convince method to add/delete timeout event.
package timer

import (
	"sync"
	"time"
)

const (
	// Max timer events.
	TIMER_QUEUE_SIZE = 1024

	// Invalid tick.
	INVALID_TICK = -1

	// second contains nanoseconds.
	SECOND = 1000000000

	// the interval check the timer.
	TICK_INTERVAL = SECOND / 10
)

// Timer type
type Timer struct {
	tick    int
	tickObj *time.Ticker
	lock    sync.Mutex
	evtTree *SortedTree
}

// NewTimer will create and iniitial a Timer
func NewTimer() *Timer {
	tmr := new(Timer)
	tmr.tick = 0
	tmr.tickObj = time.NewTicker(TICK_INTERVAL)
	tmr.evtTree = NewSortedTree(NewEventObject)
	go tmr.checkTick()
	return tmr
}

// Stop will stop a timer, not good, the go runtime will not be stoppped,
// we'll fixed it later.
func (t *Timer) Stop() {
	t.tickObj.Stop()
}

func (t *Timer) checkTick() {
	for {
		select {
		case <-t.tickObj.C:
			t.tick++
			var sn = t.evtTree.Find(t.tick)
			for sn != nil {
				evt := sn.Data.(*eventObject)
				go evt.fn()
				pn := sn.LChild
				t.evtTree.Erase(sn)
				sn = pn
			}
			break
		}
	}
}

// SetTimeout will set a callback to the timer, the callback will be called
// after secs seconds timeout.
func (t *Timer) SetTimeout(secs int, cb TimeoutCallback) (int, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	var tid int
	var err error
	tid, err = t.evtTree.Insert(cb, t.tick+int(secs*SECOND/TICK_INTERVAL))
	return tid, err
}

// KillTimer will delete the timeout event which specified by `tid` on the
// event tree. it's the only way to cancel the timeout.
func (t *Timer) KillTimer(tid int) {
	if tid == INVALID_TICK {
		return
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	t.evtTree.EraseById(tid)
}

// PrintTree just for debug. it will print the current event tree.
func (t *Timer) PrintTree() {
	t.evtTree.Each(func(ev TreeDataType) {
		evt := ev.(*eventObject)
		println(evt.tick)
	}, false)
}
