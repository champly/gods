package skiplist

type Iterator struct {
	SDLList    *SortDoublyLinkedList
	StartIndex int64
	Desc       bool
}

func (iterator *Iterator) Feach(f func(index int64, v interface{}) (breakoff bool)) {
	// find start node
	startNode := iterator.SDLList.Head

	for {
		if f(startNode.Index, startNode.Value) {
			return
		}

		if iterator.Desc {
			if startNode.Previous != nil {
				startNode = startNode.Previous
			}
		} else {
			if startNode.Next != nil {
				startNode = startNode.Next
			}
		}
	}
}
