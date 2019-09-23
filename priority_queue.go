package mpb

import (
	"container/heap"
	"sync"
)

// A priorityQueue implements heap.Interface
type priorityQueue []*Bar

func (pq *priorityQueue) Len() int {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	return len(pq.items)
}

func (pq *priorityQueue) Less(i, j int) bool {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	return pq.items[i].priority < pq.items[j].priority
}

func (pq *priorityQueue) Swap(i, j int) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	n := len(pq.items)
	bar := x.(*Bar)
	bar.index = n
	pq.items = append(pq.items, bar)
}

func (pq *priorityQueue) Pop() interface{} {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	old := pq.items
	n := len(old)
	bar := old[n-1]
	bar.index = -1 // for safety
	pq.items = old[0 : n-1]
	return bar
}

// update modifies the priority of a Bar in the queue.
func (pq *priorityQueue) update(bar *Bar, priority int) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	bar.priority = priority
	heap.Fix(pq, bar.index)
}

func (pq *priorityQueue) maxNumP() int {
	if pq.Len() == 0 {
		return 0
	}

	pq.lock.Lock()
	max := pq.items[0].NumOfPrependers()
	pq.lock.Unlock()
	for i := 1; i < pq.Len(); i++ {
		pq.lock.Lock()
		n := pq.items[i].NumOfPrependers()
		pq.lock.Unlock()
		if n > max {
			max = n
		}
	}
	return max
}

func (pq *priorityQueue) maxNumA() int {
	if pq.Len() == 0 {
		return 0
	}

	pq.lock.Lock()
	max := pq.items[0].NumOfAppenders()
	pq.lock.Unlock()
	for i := 1; i < pq.Len(); i++ {
		pq.lock.Lock()
		n := pq.items[i].NumOfAppenders()
		pq.lock.Unlock()
		if n > max {
			max = n
		}
	}
	return max
}
