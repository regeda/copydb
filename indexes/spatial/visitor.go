package spatial

import (
	"container/heap"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

type visitor struct {
	angle   s1.Angle
	latlng  s2.LatLng
	k       int
	closest itemsQueue
}

func (v *visitor) visit(c cell) {
	for i := range c {
		dist := i.latlng.Distance(v.latlng)
		if dist < v.angle {
			heap.Push(&v.closest, queueItem{Item: i.it, angle: dist})
			if v.closest.Len() > v.k {
				heap.Pop(&v.closest)

				top := v.closest[0]

				v.angle = top.angle
			}
		}
	}
}
