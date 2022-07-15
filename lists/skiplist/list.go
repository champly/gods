package skiplist

import (
	"log"
	"math/rand"
	"time"
)

const (
	Level0 = iota
	Level1
	Level2
	Level3
	Level4
	LevelLast
)

var (
	// https://github.com/redis/redis/blob/unstable/src/t_zset.c#L122
	SkipListPoint    = 0.25
	MaxRandThreshold = 10000
	RandThreshold    = int(SkipListPoint * float64(MaxRandThreshold))
)

// New build SkipList
func New(indexLevel int) *List {
	if indexLevel > LevelLast || indexLevel < Level1 {
		log.Fatalf("Max support %d index level, or less than %d", indexLevel, Level1)
	}
	list := &List{
		Rand:      rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		MaxLevel:  indexLevel,
		LevelList: make([]*IndexList, indexLevel),
	}

	for i := range list.LevelList {
		list.LevelList[i] = &IndexList{
			Level: i,
		}
	}

	return list
}

type List struct {
	Rand      *rand.Rand
	MaxLevel  int
	LevelList []*IndexList
}

type IndexList struct {
	DataList *SortDubboLinkedList
	Level    int
}

type IndexData struct {
	Reference *Node
}

func (list *List) Query(index int64) (value interface{}, exist bool) {
	// var startNode *Node
	for i := list.MaxLevel; i > Level1; i-- {
		if list.LevelList[i].DataList == nil {
			list.LevelList[i].DataList = &SortDubboLinkedList{}
		}
		currentIndexList := list.LevelList[i]

		currentIndexList.DataList.Head.FindNode(index)
	}

	panic("not implement")
}

func (list *List) Set(index int64, value interface{}) {
	level := list.randomLevel()
	if level == Level0 {
		// set to data
		list.directSetData(index, value)
		return
	}

	// build index
	list.setWithIndex(level, index, value)
}

func (list *List) directSetData(index int64, value interface{}) {
}

func (list *List) setWithIndex(level int, index int64, value interface{}) {
	// create current level index
}

func (list *List) randomLevel() int {
	var level int
	for list.Rand.Intn(MaxRandThreshold) < RandThreshold {
		level++
	}

	if level > list.MaxLevel {
		return list.MaxLevel
	}
	return level
}
