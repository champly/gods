package skiplist

import (
	"log"
	"testing"
)

func TestSkipListSet(t *testing.T) {
	list := newSkipList(IndexLevelMax)

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
	list := NewSkipList(IndexLevelMax)
	count := b.N

	b.ResetTimer()
	for i := count; i > 0; i-- {
		list.Set(int64(random(1_0000)), i)
	}
	b.Log("count", count)
}

func TestSkipListFindLargestNodeNotLargerThanIndex(t *testing.T) {
	list := buildSkipList(1, 2, 3, 4, 5, 6, 9)
	index, _, ok := list.FindLargestNodeNotLargerThanIndex(3).GetValue()
	if !ok {
		t.Error("should found node")
		return
	}
	if index != 3 {
		t.Errorf("found node index should %d but got %d", 3, index)
		return
	}

	index, _, ok = list.FindLargestNodeNotLargerThanIndex(7).GetValue()
	if !ok {
		t.Error("should found node")
		return
	}
	if index != 6 {
		t.Errorf("found node index should %d but got %d", 6, index)
		return
	}

	list2 := buildSkipList(2)
	_, _, ok = list2.FindLargestNodeNotLargerThanIndex(1).GetValue()
	if ok {
		t.Error("should not found node")
		return
	}
	index, _, ok = list2.FindLargestNodeNotLargerThanIndex(3).GetValue()
	if !ok {
		t.Error("should found node")
		return
	}
	if index != 2 {
		t.Errorf("found node index should %d but got %d", 2, index)
	}
}

func TestSkipListFindSmallestNodeNotSmallerThanIndexForward(t *testing.T) {
	list := buildSkipList(1, 2, 3, 4, 5, 6, 9)
	index, _, ok := list.FindSmallestNodeNotSmallerThanIndex(3).GetValue()
	if !ok {
		t.Error("should found node")
		return
	}
	if index != 3 {
		t.Errorf("found node index should %d but got %d", 3, index)
		return
	}

	index, _, ok = list.FindSmallestNodeNotSmallerThanIndex(7).GetValue()
	if !ok {
		t.Error("should found node")
		return
	}
	if index != 9 {
		t.Errorf("found node index should %d but got %d", 9, index)
		return
	}

	list2 := buildSkipList(2)
	_, _, ok = list2.FindSmallestNodeNotSmallerThanIndex(3).GetValue()
	if ok {
		t.Error("should not found node")
		return
	}
	index, _, ok = list2.FindSmallestNodeNotSmallerThanIndex(1).GetValue()
	if !ok {
		t.Error("should found node")
		return
	}
	if index != 2 {
		t.Errorf("found node index should %d but got %d", 2, index)
	}
}

func buildSkipList(dataList ...int) *SkipList {
	list := newSkipList(IndexLevelMax)
	for _, data := range dataList {
		list.Set(int64(data), data)
	}
	return list
}
