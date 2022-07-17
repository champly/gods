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
		list.Set(int64(random(100)), i)
	}
	b.Log("count", count)
}
