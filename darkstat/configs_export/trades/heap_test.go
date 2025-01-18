package trades

/*
Content of this file is copy pasting somewhere from https://pkg.go.dev/container/heap
*/

import (
	"container/heap"
	"fmt"
	"testing"
)

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func TestHeap(t *testing.T) {
	// Some items and their priorities.
	items := map[int]int{
		0: 3, 1: 2, 2: 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value_weight: value,
			priority:     priority,
			index:        i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value_weight: 5,
		priority:     1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value_weight, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%d ", item.priority, item.value_weight)
	}
}
