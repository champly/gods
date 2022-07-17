package skiplist

import (
	"log"
	"math/rand"
	"time"
)

const (
	Level0    = 0
	Level1    = 1
	LevelLast = 5
)

var (
	// https://github.com/redis/redis/blob/unstable/src/t_zset.c#L122
	SkipListPoint    = 0.25
	MaxRandThreshold = 10000
	RandThreshold    = int(SkipListPoint * float64(MaxRandThreshold))
)

// New build SkipList
func New(indexLevel int) *SkipList {
	if indexLevel > LevelLast || indexLevel < Level1 {
		log.Fatalf("Max support %d index level, or less than %d", indexLevel, Level1)
	}
	list := &SkipList{
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

type SkipList struct {
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

func (list *SkipList) Query(index int64) (value interface{}, exist bool) {
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

func (list *SkipList) Set(index int64, value interface{}) {
	level := list.randomLevel()
	if level == Level0 {
		// set to data
		list.directSetData(index, value)
		return
	}

	// build index
	list.setWithIndex(level, index, value)
}

func (list *SkipList) directSetData(index int64, value interface{}) {
}

func (list *SkipList) setWithIndex(level int, index int64, value interface{}) {
	// create current level index
}

func (list *SkipList) randomLevel() int {
	var level int
	for list.Rand.Intn(MaxRandThreshold) < RandThreshold {
		level++
	}

	if level > list.MaxLevel {
		return list.MaxLevel
	}
	return level
}
