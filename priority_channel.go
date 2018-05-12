package prioritychannel

import (
	"container/heap"
	"sync"
)

type Item struct {
	Priority int
	Value    interface{}
}

type items []Item

type PriorityChannel struct {
	items items
	In    chan<- Item
	Out   <-chan Item
}

func (pc PriorityChannel) Len() int { return pc.items.Len() }

func (is items) Len() int { return len(is) }

func (is items) Less(i, j int) bool {
	return is[i].Priority >= is[j].Priority
}

func (is *items) Swap(i, j int) {
	isOld := *is
	isOld[i], isOld[j] = isOld[j], isOld[i]
	*is = isOld
}

func (is *items) Push(x interface{}) {
	*is = append(*is, x.(Item))
}

func (is *items) Pop() interface{} {
	isOld := *is
	n := len(isOld)
	i := isOld[n-1]
	*is = isOld[0 : n-1]
	return i
}

func New() *PriorityChannel {
	pc := PriorityChannel{}
	in := make(chan Item)
	out := make(chan Item)
	quit := make(chan struct{})
	mtx := sync.Mutex{}
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go func() {
		for {
			item, ok := <-in
			if !ok {
				quit <- struct{}{}
				return
			}

			mtx.Lock()
			if pc.Len() == 0 {
				wg.Done()
			}
			heap.Push(&pc.items, item)
			mtx.Unlock()
		}
	}()

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
			}

			wg.Wait()

			mtx.Lock()
			item := heap.Pop(&pc.items).(Item)
			if pc.Len() == 0 {
				wg.Add(1)
			}
			mtx.Unlock()

			out <- item
		}
	}()

	pc.In = in
	pc.Out = out
	return &pc
}
