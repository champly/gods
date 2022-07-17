package skiplist

import (
	"log"
	"testing"
)

func TestSkipListSet(t *testing.T) {
	list := New(IndexLevelMax)

	count := 1000000
	currentData := map[int]int{}

	for i := count; i >= 0; i-- {
		index := random(100)
		value := random(count)

		list.Set(int64(index), value)

		currentData[index] = value
	}
	for i := IndexLevelMax; i >= 0; i-- {
		log.Println("----> level:", i)
		l := list.DataLevelList[i]
		pl := l.Head
		for pl != nil {
			id, ok := pl.Value.(*IndexData)
			if !ok {
				break
			}
			id.Println()
			pl = pl.Next

		}
	}
	list.DataLevelList[DataLevel0].Print()

	pHead := list.DataLevelList[DataLevel0].Head
	for pHead != nil {
		if currentData[int(pHead.Index)] != pHead.Value.(int) {
			t.Errorf("index %d want %d but got %d", pHead.Index, currentData[int(pHead.Index)], pHead.Value.(int))
			return
		}
		pHead = pHead.Next
	}
}

func BenchmarkSkipList(b *testing.B) {
	list := New(IndexLevelMax)
	count := b.N

	b.ResetTimer()
	for i := count; i > 0; i-- {
		list.Set(int64(random(1_0000)), i)
	}
	b.Log("count", count)
}

func TestSkipListFindLargestNodeNotLargerThanIndex(t *testing.T) {
	list := buildSkipList(1, 2, 3, 4, 5, 6, 9)
	node := list.FindLargestNodeNotLargerThanIndex(3)
	if node == nil {
		t.Error("should found node")
		return
	}
	if node.Index != 3 {
		t.Errorf("found node index should %d but got %d", 3, node.Index)
		return
	}

	node = list.FindLargestNodeNotLargerThanIndex(7)
	if node == nil {
		t.Error("should found node")
		return
	}
	if node.Index != 6 {
		t.Errorf("found node index should %d but got %d", 6, node.Index)
		return
	}

	list2 := buildSkipList(2)
	if list2.FindLargestNodeNotLargerThanIndex(1) != nil {
		t.Error("should not found node")
		return
	}
	node = list2.FindLargestNodeNotLargerThanIndex(3)
	if node == nil {
		t.Error("should found node")
		return
	}
	if node.Index != 2 {
		t.Errorf("found node index should %d but got %d", 2, node.Index)
	}
}

func TestSkipListFindSmallestNodeNotSmallerThanIndexForward(t *testing.T) {
	list := buildSkipList(1, 2, 3, 4, 5, 6, 9)
	node := list.FindSmallestNodeNotSmallerThanIndex(3)
	if node == nil {
		t.Error("should found node")
		return
	}
	if node.Index != 3 {
		t.Errorf("found node index should %d but got %d", 3, node.Index)
		return
	}

	node = list.FindSmallestNodeNotSmallerThanIndex(7)
	if node == nil {
		t.Error("should found node")
		return
	}
	if node.Index != 9 {
		t.Errorf("found node index should %d but got %d", 9, node.Index)
		return
	}

	list2 := buildSkipList(2)
	if list2.FindSmallestNodeNotSmallerThanIndex(3) != nil {
		t.Error("should not found node")
		return
	}
	node = list2.FindSmallestNodeNotSmallerThanIndex(1)
	if node == nil {
		t.Error("should found node")
		return
	}
	if node.Index != 2 {
		t.Errorf("found node index should %d but got %d", 2, node.Index)
	}
}

func buildSkipList(dataList ...int) *SkipList {
	list := New(IndexLevelMax)
	for _, data := range dataList {
		list.Set(int64(data), data)
	}
	return list
}
