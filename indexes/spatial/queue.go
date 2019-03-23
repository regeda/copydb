package spatial

import (
	"github.com/golang/geo/s1"
	"github.com/regeda/copydb"
)

type queueItem struct {
	copydb.Item
	angle s1.Angle
	i     int
}

type itemsQueue []queueItem

func newItemsQueue(capacity int) itemsQueue {
	return make([]queueItem, 0, capacity+1)
}

func (q itemsQueue) Len() int { return len(q) }

func (q itemsQueue) Less(i, j int) bool { return q[i].angle > q[j].angle }

func (q itemsQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].i = i
	q[j].i = j
}

func (q *itemsQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(queueItem)
	item.i = n
	*q = append(*q, item)
}

func (q *itemsQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.i = -1
	*q = old[0 : n-1]
	return item
}
