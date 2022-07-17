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
	DataLevelList []*IndexList
}

// New build SkipList
func New(indexLevel int) *SkipList {
	if indexLevel > IndexLevelMax || indexLevel < IndexLevel1 {
		log.Fatalf("Max support %d index level, or less than %d", indexLevel, IndexLevel1)
	}
	list := &SkipList{
		Rand:          rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		MaxIndexLevel: indexLevel,
		DataLevelList: make([]*IndexList, indexLevel),
	}

	for i := range list.DataLevelList {
		list.DataLevelList[i] = &IndexList{
			Level: i,
		}
	}

	return list
}

func (skipl *SkipList) FindValue(index int64) (value interface{}, exist bool) {
	node, ok := skipl.FindNode(index)
	if !ok {
		return nil, false
	}
	return node.Value, true
}

func (skipl *SkipList) FindNode(index int64) (*Node, bool) {
	// find node by index
	for level := skipl.MaxIndexLevel; level > DataLevel0; level-- {
		indexNode, ok := findIndexNode(skipl.getDataListWithLevel(level), index)
		if !ok {
			// not found index node, find next index
			continue
		}

		foundNode, ok := findDataWithIndexData(indexNode, index)
		if !ok {
			return nil, false
		}
		if foundNode.Index == index {
			return foundNode, true
		}
		foundNode, ok = foundNode.FindNode(index)
		if !ok {
			return nil, false
		}
		return foundNode, true
	}
	return nil, false
}

func (skipl *SkipList) Set(index int64, value interface{}) {
	foundNode, ok := skipl.FindNode(index)
	if ok {
		// update
		foundNode.Value = value
		return
	}

	// create new node
	level := skipl.randomLevel()
	if level == DataLevel0 {
		// set to data
		skipl.directSetData(index, value)
		return
	}

	// build index
	skipl.setWithIndex(level, index, value)
}

func (skipl *SkipList) getDataListWithLevel(level int) *IndexList {
	switch {
	case level > skipl.MaxIndexLevel:
		return skipl.DataLevelList[skipl.MaxIndexLevel]
	// case level > DataLevel:
	//     return list.LevelList[Level0]
	default:
		return skipl.DataLevelList[level]
	}
}

func (skipli *SkipList) directSetData(index int64, value interface{}) {
	indexList := skipli.getDataListWithLevel(DataLevel0)
	indexList.DataList.Set(index, value)
}

func (skipl *SkipList) setWithIndex(level int, index int64, value interface{}) {
	// create current level index
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
