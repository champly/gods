package skiplist

import "testing"

func TestNodeRemove(t *testing.T) {
	list := buildSortDubboLinkedList(1, 2, 3, 4)
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
