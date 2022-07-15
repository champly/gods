package skiplist

type Node struct {
	Index int64
	Value interface{}

	Previous *Node
	Next     *Node
}

func (node *Node) Set(sldList *SortDubboLinkedList, index int64, value interface{}) {
	newNode := &Node{Index: index, Value: value}
	if node == nil {
		// first node
		sldList.Head = newNode
		sldList.Tail = newNode
		return
	}

	switch {
	case node.Index == index:
		// udpate
		node.Value = value
	case node.Index > index:
		// <- forward set
		node.setForward(sldList, newNode)
	default:
		// -> later set
		node.setLater(sldList, newNode)
	}
}

func (node *Node) setForward(sldList *SortDubboLinkedList, newNode *Node) {
	for node != nil {
		if node.Index == newNode.Index {
			// udpate
			node.Value = newNode.Value
			return
		}

		// not match condition, find next
		if node.Index > newNode.Index {
			node = node.Previous
			continue
		}

		// currentNode.Index < newNode.Index
		// should insert current previous
		if node.Next == nil {
			// currentNode is Tail node, create chain
			//   Tail -> newNode
			// 1. Tail -> newNode
			sldList.Tail = newNode
		} else {
			// currentNode.Previous.Index < newNode.Index
			// should insert between currentNode.Next and currentNode
			//    newNode <-> currentNode.Next
			// 1. newNode <- currentNode.Next
			node.Next.Previous = newNode
			// 2. newNode -> currentNode.Previous
			newNode.Next = node.Next
		}

		// currentNode <-> newNode
		// 1. newNode -> currentNode
		newNode.Previous = node
		// 2. newNode <- node
		node.Next = newNode
		return
	}

	// loop end, not match, should insert head.
	if sldList.Head == nil {
		// first node
		sldList.Head = newNode
		sldList.Tail = newNode
		return
	}

	// Head -> headNode <-> newNode
	// 1. newNode -> headNode
	newNode.Next = sldList.Head
	// 2. newNode <- headNode
	sldList.Head.Previous = newNode
	// 3. Head -> newNode
	sldList.Head = newNode
}

func (node *Node) setLater(sldList *SortDubboLinkedList, newNode *Node) {
	for node != nil {
		if node.Index == newNode.Index {
			// udpate
			node.Value = newNode.Value
			return
		}

		// not match condition, find next
		if node.Index < newNode.Index {
			node = node.Next
			continue
		}

		// currentNode.Index > newNode.Index
		// should insert current previous
		if node.Previous == nil {
			// currentNode is Head node, create chain
			//   Head -> newNode
			// 1. Head -> newNode
			sldList.Head = newNode
		} else {
			// currentNode.Previous.Index < newNode.Index
			// should insert between currentNode.Previous and currentNode
			//    currentNode.Previous <-> newNode
			// 1. newNode <- currentNode.Previous
			node.Previous.Next = newNode
			// 2. newNode -> currentNode.Previous
			newNode.Previous = node.Previous
		}

		// newNode <-> currentNode
		// 1. newNode -> currentNode
		newNode.Next = node
		// 2. newNode <- node
		node.Previous = newNode
		return
	}

	// loop end, not match, should insert tail.
	if sldList.Tail == nil {
		// first node
		sldList.Head = newNode
		sldList.Tail = newNode
		return
	}

	// tailNode <-> newNode <- Tail
	// 1. newNode -> tailNode
	newNode.Previous = sldList.Tail
	// 2. newNode <- tailNode
	sldList.Tail.Next = newNode
	// 3. Tail -> newNode
	sldList.Tail = newNode
}

func (node *Node) FindNode(index int64) (n *Node, exist bool) {
	if node == nil {
		return nil, false
	}

	switch {
	case node.Index == index:
		return node, true
	case node.Index > index:
		return node.findForward(index)
	default:
		return node.findLaster(index)
	}
}

func (node *Node) findForward(index int64) (n *Node, exist bool) {
	for node != nil {
		if node.Index == index {
			return node, true
		}
		if node.Index < index {
			return nil, false
		}
		node = node.Previous
	}
	return nil, false
}

func (node *Node) findLaster(index int64) (n *Node, exist bool) {
	for node != nil {
		if node.Index == index {
			return node, true
		}
		if node.Index > index {
			return nil, false
		}
		node = node.Next
	}
	return nil, false
}

func (node *Node) Remove(sdlList *SortDubboLinkedList, index int64) {
	if node == nil {
		return
	}

	switch {
	case node.Index == index:
		node.removeCurrentNode(sdlList)
	case node.Index > index:
		node.removeForward(sdlList, index)
	default:
		node.removeLaster(sdlList, index)
	}
}

func (node *Node) removeForward(sdlList *SortDubboLinkedList, index int64) {
	for node != nil {
		if node.Index > index {
			node = node.Previous
			continue
		}

		if node.Index < index {
			// not found return
			return
		}

		// node.Index == index
		node.removeCurrentNode(sdlList)
		return
	}
}

func (node *Node) removeLaster(sdlList *SortDubboLinkedList, index int64) {
	for node != nil {
		if node.Index < index {
			node = node.Next
			continue
		}
		if node.Index > index {
			// not found return
			return
		}

		// node.Index == index
		node.removeCurrentNode(sdlList)
		return
	}
}

func (node *Node) removeCurrentNode(sdlList *SortDubboLinkedList) {
	if node.Previous == nil {
		// currentNode is Head
		// Head -> nextNode
		sdlList.Head = node.Next
	} else {
		// previousNode -> nextNode
		node.Previous.Next = node.Next
	}

	if node.Next == nil {
		// currentNode is Tail
		// Tail -> previousNode
		sdlList.Tail = node.Previous
	} else {
		// previousNode <- nextNode
		node.Next.Previous = node.Previous
	}
}
