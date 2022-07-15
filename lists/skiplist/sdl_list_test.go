package skiplist

import (
	"math/rand"
	"testing"
	"time"
)

var (
	r = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
)

func TestSortDubboLinkedListSetAndSet(t *testing.T) {
	list := &SortDubboLinkedList{}

	count := 10000
	currentData := map[int]int{}

	var node *Node
	var findNode bool

	for i := count; i >= 0; i-- {
		index := random(100)
		value := random(count)

		if findNode {
			node.Set(list, int64(index), value)
		} else {
			list.Set(int64(index), value)
		}

		if node != nil {
			node, findNode = node.FindNode(int64(index))
		} else {
			node, findNode = list.FindNode(int64(index))
		}
		if !findNode {
			t.Error("insert into list, but not found.")
			return
		}

		currentData[index] = value
	}
	list.Print()

	pHead := list.Head
	for pHead != nil {
		if currentData[int(pHead.Index)] != pHead.Value.(int) {
			t.Errorf("index %d want %d but got %d", pHead.Index, currentData[int(pHead.Index)], pHead.Value.(int))
			return
		}
		pHead = pHead.Next
	}
}

func BenchmarkSortDubboLinkedListSet(b *testing.B) {
	list := &SortDubboLinkedList{}
	count := b.N

	b.ResetTimer()
	for i := count; i > 0; i-- {
		list.Set(int64(random(100)), i)
	}
	b.Log("count", count)
}

func TestSortDubboLinkedListRemove(t *testing.T) {
	list := buildSortDubboLinkedList(1, 2, 3, 4, 10)
	currentNode, ok := list.Tail.FindNode(1)
	if !ok {
		t.Error("should found node")
		return
	}

	// remove head
	currentNode.Remove(list, 1)
	list.Print()

	// remove middle
	currentNode, ok = list.Tail.FindNode(3)
	if !ok {
		t.Error("should found node")
		return
	}
	currentNode.Remove(list, 3)
	list.Print()

	// remove tail
	currentNode, ok = list.Tail.FindNode(4)
	if !ok {
		t.Error("should found node")
		return
	}
	currentNode.Remove(list, 4)
	list.Print()
}

func buildSortDubboLinkedList(dataList ...int) *SortDubboLinkedList {
	list := &SortDubboLinkedList{}
	for _, data := range dataList {
		list.Set(int64(data), data)
	}
	return list
}
func random(max int) int {
	return r.Intn(max)
}
