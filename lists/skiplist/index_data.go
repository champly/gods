package skiplist

type IndexList struct {
	DataList *SortDubboLinkedList
	Level    int
}

type IndexData struct {
	Index      int64
	Reference  *Node
	IndexLevel int
}

func findIndexNode(il *IndexList, index int64) (*Node, bool) {
	if il == nil || il.DataList == nil {
		return nil, false
	}
	indexNode := il.DataList.FindLargestNodeNotLargerThanIndex(index)
	if indexNode == nil {
		return nil, false
	}
	return indexNode, true
}

func findDataWithIndexData(indexNode *Node, index int64) (foundNode *Node, found bool) {
	if indexNode == nil {
		return nil, false
	}

	switch node := indexNode.Value.(type) {
	case *IndexData:
		if indexNode.Index == index {
			return findDataWithIndexData(node.Reference, index)
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
		return findDataWithIndexData(n.Reference, index)

	default:
		// currentNode is Data SortDubboLinkedList
		return indexNode, true
	}
}
