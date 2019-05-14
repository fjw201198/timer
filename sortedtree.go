package timer

type TreeDataType interface {
	Less(rh TreeDataType) bool
	Clear()
	SetParam(param ...interface{})
	DataLess(rh interface{}) bool
}

type NewDataFuncType func(param ...interface{}) TreeDataType

type sortedNode struct {
	Data   TreeDataType
	LChild *sortedNode
	RChild *sortedNode
	id     int
}

type SortedTree struct {
	pool *eventPool
	Root *sortedNode
}

func NewSortedTree(newDataFn NewDataFuncType) *SortedTree {
	st := new(SortedTree)
	st.pool = newEventPool(newDataFn)
	st.Root, _ = st.pool.alloc()
	return st
}

func (st *SortedTree) findById(where *sortedNode, id int, parent **sortedNode, dst **sortedNode, right *bool) {
	if *dst != nil || where == nil {
		return
	}
	if where.LChild != nil {
		if where.LChild.id == id {
			*parent = where
			*dst = where.LChild
			*right = false
			return
		}
	}
	if where.RChild != nil {
		if where.RChild.id == id {
			*parent = where
			*dst = where.RChild
			*right = true
			return
		}
	}
	st.findById(where.LChild, id, parent, dst, right)
	st.findById(where.RChild, id, parent, dst, right)
}

func (st *SortedTree) findAndInsert(where *sortedNode, what *sortedNode) {
	// middle order
	bret := what.Data.Less(where.Data)
	if bret {
		if where.LChild == nil {
			where.LChild = what
			return
		}
		st.findAndInsert(where.LChild, what)
	} else {
		if where.RChild == nil {
			where.RChild = what
			return
		}
		st.findAndInsert(where.RChild, what)
	}
}

func (st *SortedTree) Insert(param ...interface{}) (int, error) {
	sn, err := st.pool.alloc()
	if err != nil {
		return INVALID_TICK, err
	}
	sn.Data.SetParam(param...)
	st.findAndInsert(st.Root, sn)
	return sn.id, nil
}

func (st *SortedTree) EraseNode(sn *sortedNode, pn *sortedNode, right bool) {
	var tmpPar = sn.LChild
	if right {
		pn.RChild = tmpPar
	} else {
		pn.LChild = tmpPar
	}

	if tmpPar == nil {
		if right {
			pn.RChild = sn.RChild
		} else {
			pn.LChild = sn.RChild
		}
		return
	}
	if tmpPar.LChild == nil {
		tmpPar.LChild = sn.RChild
	} else {
		temp := tmpPar.LChild
		for temp.RChild != nil {
			temp = temp.RChild
		}
		if temp != nil {
			temp.RChild = sn.RChild
		}
	}
	sn.LChild = nil
	sn.RChild = nil
	st.pool.free(sn)
}

func (st *SortedTree) EraseById(tid int) {
	var sn *sortedNode = nil
	var pn *sortedNode = nil
	var right bool = false
	st.findById(st.Root, tid, &pn, &sn, &right)
	if sn == nil {
		return
	}
	st.EraseNode(sn, pn, right)
}

func (st *SortedTree) Erase(sn *sortedNode) {
	st.EraseById(sn.id)
}

func (st *SortedTree) Each(fn func(data TreeDataType), asynced bool) {
	if asynced {
		st.eachHelperAsync(st.Root.LChild, fn)
		st.eachHelperAsync(st.Root.RChild, fn)
	} else {
		st.eachHelper(st.Root.LChild, fn)
		st.eachHelper(st.Root.RChild, fn)
	}
}

func (st *SortedTree) eachHelperAsync(curNode *sortedNode, fn func(data TreeDataType)) {
	if curNode == nil {
		return
	}
	go fn(curNode.Data)
	st.eachHelperAsync(curNode.LChild, fn)
	st.eachHelperAsync(curNode.RChild, fn)
}

func (st *SortedTree) eachHelper(curNode *sortedNode, fn func(data TreeDataType)) {
	if curNode == nil {
		return
	}
	fn(curNode.Data)
	st.eachHelper(curNode.LChild, fn)
	st.eachHelper(curNode.RChild, fn)
}

func (st *SortedTree) Find(data interface{}) *sortedNode {
	var ret *sortedNode
	st.findDataHelper(data, st.Root.LChild, &ret)
	st.findDataHelper(data, st.Root.RChild, &ret)
	return ret
}

func (st *SortedTree) findDataHelper(data interface{}, src *sortedNode, dst **sortedNode) {
	if src == nil {
		return
	}
	if *dst != nil {
		return
	}
	if src.Data.DataLess(data) {
		*dst = src
		return
	}
	st.findDataHelper(data, src.LChild, dst)
	st.findDataHelper(data, src.RChild, dst)
}
