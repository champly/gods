package skiplist

import (
	"fmt"
)

type IndexList struct {
	DataList *SortDoublyLinkedList
}

type IndexData struct {
	Index     int64
	Reference *SortDoublyLinkedListNode
}

func findIndexNode(sdll *SortDoublyLinkedList, index int64) (*SortDoublyLinkedListNode, bool) {
	if sdll == nil {
		return nil, false
	}
	indexNode := sdll.FindLargestNodeNotLargerThanIndex(index)
	if indexNode == nil {
		return nil, false
	}
	return indexNode, true
}

func findDataIndexNodeWithIndexNode(indexNode *SortDoublyLinkedListNode, index int64) (foundNode *SortDoublyLinkedListNode, found bool) {
	if indexNode == nil {
		return nil, false
	}

	switch node := indexNode.Value.(type) {
	case *IndexData:
		if indexNode.Index == index {
			return findDataIndexNodeWithIndexNode(node.Reference, index)
		}
		newIndexNode := indexNode.FindLargestNodeNotLargerThanIndex(index)
		if newIndexNode == nil {
			// not do this normal!
			return nil, false
		}
		n, ok := newIndexNode.Value.(*IndexData)
		if !ok {
			// not do this
			return nil, false
		}
		return findDataIndexNodeWithIndexNode(n.Reference, index)

	default:
		// currentNode is Data SortDubboLinkedList
		return indexNode, true
	}
}

func (id *IndexData) Println() {
	fmt.Printf("%d", id.Index)
	nextIndexData, ok := id.Reference.Value.(*IndexData)
	if ok {
		fmt.Print("->")
		nextIndexData.Println()
		return
	}
	fmt.Printf("->%d\n", id.Index)
}
