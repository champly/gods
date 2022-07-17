package skiplist

import (
	"log"
	"math/rand"
	"time"
)

const (
	DataLevel0    = 0
	IndexLevel1   = 1
	IndexLevelMax = 5
)

var (
	// https://github.com/redis/redis/blob/unstable/src/t_zset.c#L122
	SkipListPoint    = 0.25
	MaxRandThreshold = 10000
	RandThreshold    = int(SkipListPoint * float64(MaxRandThreshold))
)

type SkipList struct {
	Rand          *rand.Rand
	MaxIndexLevel int
	DataLevelList []*SortDoublyLinkedList
}

// New build SkipList
func New(indexLevel int) *SkipList {
	if indexLevel > IndexLevelMax || indexLevel < IndexLevel1 {
		log.Fatalf("Max support %d index level, or less than %d", indexLevel, IndexLevel1)
	}
	list := &SkipList{
		Rand:          rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		MaxIndexLevel: indexLevel,
		DataLevelList: make([]*SortDoublyLinkedList, indexLevel+1),
	}

	for i := range list.DataLevelList {
		list.DataLevelList[i] = &SortDoublyLinkedList{}
	}

	return list
}
func (skipl *SkipList) Set(index int64, value interface{}) {
	foundDataNode, ok := skipl.findNode(index)
	if ok {
		// update
		foundDataNode.Value = value
		return
	}

	// create new node
	level := skipl.randomLevel()
	if level == DataLevel0 {
		// set to data
		skipl.createNodeDirectly(index, value)
		return
	}

	// build index
	skipl.createNodeWithIndex(level, index, value)
}

// createNodeDirectly set Data without create index
func (skipl *SkipList) createNodeDirectly(index int64, value interface{}) {
	// indexList := skipl.getDataListWithLevel(DataLevel0)
	// indexList.DataList.Set(index, value)
	dataIndexNode := skipl.findDataIndexNode(index)
	dataIndexNode.Set(skipl.getDataListWithLevel(DataLevel0), index, value)
}

func (skipl *SkipList) createNodeWithIndex(level int, index int64, value interface{}) {
	// create current level index
	dataIndexNode := skipl.findDataIndexNode(index)
	newDataNode := dataIndexNode.Set(skipl.getDataListWithLevel(DataLevel0), index, value)

	// create index node
	tmpIndexNode := newDataNode
	for l := 1; l <= level; l++ {
		tmpIndexNode = skipl.createIndex(l, index, tmpIndexNode)
	}
}

func (skipl *SkipList) createIndex(level int, index int64, reference *SortDoublyLinkedListNode) *SortDoublyLinkedListNode {
	currentLevelIndexList := skipl.getDataListWithLevel(level)
	return currentLevelIndexList.Set(index, &IndexData{Index: index, Reference: reference})
}

func (skipl *SkipList) FindValue(index int64) (value interface{}, exist bool) {
	node, ok := skipl.findNode(index)
	if !ok {
		return nil, false
	}
	return node.Value, true
}

func (skipl *SkipList) findNode(index int64) (*SortDoublyLinkedListNode, bool) {
	dataIndexNode := skipl.findDataIndexNode(index)
	return dataIndexNode.FindNode(index)
}

func (skipl *SkipList) findDataIndexNode(index int64) *SortDoublyLinkedListNode {
	// find node by index
	for level := skipl.MaxIndexLevel; level > DataLevel0; level-- {
		indexNode, ok := findIndexNode(skipl.getDataListWithLevel(level), index)
		if !ok {
			// not found index node, find next index
			continue
		}

		dataIndexNode, ok := findDataIndexNodeWithIndexNode(indexNode, index)
		if ok {
			return dataIndexNode
		}
	}

	// not found, use Level0 to found
	return skipl.getDataListWithLevel(DataLevel0).Head
}

func (skipl *SkipList) getDataListWithLevel(level int) *SortDoublyLinkedList {
	switch {
	case level > skipl.MaxIndexLevel:
		return skipl.DataLevelList[skipl.MaxIndexLevel]
	// case level > DataLevel:
	//     return list.LevelList[Level0]
	default:
		return skipl.DataLevelList[level]
	}
}

func (skipl *SkipList) randomLevel() int {
	var level int
	for skipl.Rand.Intn(MaxRandThreshold) < RandThreshold {
		level++
	}

	if level > skipl.MaxIndexLevel {
		return skipl.MaxIndexLevel
	}
	return level
}
