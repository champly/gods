package skiplist

type Iterator struct {
	currentNode *SortDoublyLinkedListNode
}

type IIterator interface {
	GetValue() (index int64, value interface{}, ok bool)
	Previous()
	Next()
}

func (iterator *Iterator) GetValue() (index int64, value interface{}, ok bool) {
	if iterator.currentNode == nil {
		return 0, nil, false
	}
	return iterator.currentNode.Index, iterator.currentNode.Value, true
}

func (iterator *Iterator) Previous() {
	if iterator.currentNode == nil {
		return
	}
	iterator.currentNode = iterator.currentNode.Previous
}

func (iterator *Iterator) Next() {
	if iterator.currentNode == nil {
		return
	}
	iterator.currentNode = iterator.currentNode.Next
}
