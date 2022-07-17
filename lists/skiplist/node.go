package skiplist

type Node struct {
	Index int64
	Value interface{}

	Previous *Node
	Next     *Node
}

func (node *Node) Set(sldList *SortDubboLinkedList, index int64, value interface{}) {
	if sldList == nil {
		return
	}

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
		setForward(sldList, node.Previous, newNode)
	default:
		// -> later set
		setLater(sldList, node.Next, newNode)
	}
}

func setForward(sldList *SortDubboLinkedList, currentNode, newNode *Node) {
	if sldList == nil || newNode == nil {
		return
	}

	for currentNode != nil {
		if currentNode.Index == newNode.Index {
			// udpate
			currentNode.Value = newNode.Value
			return
		}

		// not match condition, find next
		if currentNode.Index > newNode.Index {
			currentNode = currentNode.Previous
			continue
		}

		// currentNode.Index < newNode.Index
		// should insert current previous
		if currentNode.Next == nil {
			// currentNode is Tail node, create chain
			//   Tail -> newNode
			// 1. Tail -> newNode
			sldList.Tail = newNode
		} else {
			// currentNode.Previous.Index < newNode.Index
			// should insert between currentNode.Next and currentNode
			//    newNode <-> currentNode.Next
			// 1. newNode <- currentNode.Next
			currentNode.Next.Previous = newNode
			// 2. newNode -> currentNode.Previous
			newNode.Next = currentNode.Next
		}

		// currentNode <-> newNode
		// 1. newNode -> currentNode
		newNode.Previous = currentNode
		// 2. newNode <- node
		currentNode.Next = newNode
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

func setLater(sldList *SortDubboLinkedList, currentNode, newNode *Node) {
	if sldList == nil || newNode == nil {
		return
	}

	for currentNode != nil {
		if currentNode.Index == newNode.Index {
			// udpate
			currentNode.Value = newNode.Value
			return
		}

		// not match condition, find next
		if currentNode.Index < newNode.Index {
			currentNode = currentNode.Next
			continue
		}

		// currentNode.Index > newNode.Index
		// should insert current previous
		if currentNode.Previous == nil {
			// currentNode is Head node, create chain
			//   Head -> newNode
			// 1. Head -> newNode
			sldList.Head = newNode
		} else {
			// currentNode.Previous.Index < newNode.Index
			// should insert between currentNode.Previous and currentNode
			//    currentNode.Previous <-> newNode
			// 1. newNode <- currentNode.Previous
			currentNode.Previous.Next = newNode
			// 2. newNode -> currentNode.Previous
			newNode.Previous = currentNode.Previous
		}

		// newNode <-> currentNode
		// 1. newNode -> currentNode
		newNode.Next = currentNode
		// 2. newNode <- node
		currentNode.Previous = newNode
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
		return findForward(node.Previous, index)
	default:
		return findLaster(node.Next, index)
	}
}

func findForward(currentNode *Node, index int64) (n *Node, exist bool) {
	for currentNode != nil {
		if currentNode.Index == index {
			return currentNode, true
		}
		if currentNode.Index < index {
			return nil, false
		}
		currentNode = currentNode.Previous
	}
	return nil, false
}

func findLaster(currentNode *Node, index int64) (n *Node, exist bool) {
	for currentNode != nil {
		if currentNode.Index == index {
			return currentNode, true
		}
		if currentNode.Index > index {
			return nil, false
		}
		currentNode = currentNode.Next
	}
	return nil, false
}

func (node *Node) Remove(sdlList *SortDubboLinkedList, index int64) {
	if node == nil {
		return
	}

	switch {
	case node.Index == index:
		removeCurrentNode(sdlList, node)
	case node.Index > index:
		removeForward(sdlList, node.Previous, index)
	default:
		removeLaster(sdlList, node.Next, index)
	}
}

func removeForward(sdlList *SortDubboLinkedList, currentNode *Node, index int64) {
	for currentNode != nil {
		if currentNode.Index > index {
			currentNode = currentNode.Previous
			continue
		}

		if currentNode.Index < index {
			// not found return
			return
		}

		// currentNode.Index == index
		removeCurrentNode(sdlList, currentNode)
		return
	}
}

func removeLaster(sdlList *SortDubboLinkedList, currentNode *Node, index int64) {
	for currentNode != nil {
		if currentNode.Index < index {
			currentNode = currentNode.Next
			continue
		}
		if currentNode.Index > index {
			// not found return
			return
		}

		// currentNode.Index == index
		removeCurrentNode(sdlList, currentNode)
		return
	}
}

func removeCurrentNode(sdlList *SortDubboLinkedList, currentNode *Node) {
	if sdlList == nil {
		return
	}

	if currentNode.Previous == nil {
		// currentNode is Head
		// Head -> nextNode
		sdlList.Head = currentNode.Next
	} else {
		// previousNode -> nextNode
		currentNode.Previous.Next = currentNode.Next
	}

	if currentNode.Next == nil {
		// currentNode is Tail
		// Tail -> previousNode
		sdlList.Tail = currentNode.Previous
	} else {
		// previousNode <- nextNode
		currentNode.Next.Previous = currentNode.Previous
	}
}

// FindLargestNodeNotLargerThanIndex find the largest node not larger than the index
// eg: 1 2 3 4 5 6 9
//    if index = 3, return node which index is 3
//    if index = 7, return node which index is 6
// eg: 2 (just one node)
//    if index = 1, return nil
func (node *Node) FindLargestNodeNotLargerThanIndex(index int64) *Node {
	if node == nil {
		return nil
	}

	switch {
	case node.Index == index:
		return node
	case node.Index > index:
		return findLargestNodeNotLargerThanIndexForward(node.Previous, index)
	default:
		if node.Next == nil {
			// node is Tail
			return node
		}
		return findLargestNodeNotLargerThanIndexLater(node.Next, index)
	}
}

func findLargestNodeNotLargerThanIndexForward(currentNode *Node, index int64) *Node {
	for currentNode != nil {
		if currentNode.Index <= index {
			return currentNode
		}

		currentNode = currentNode.Previous
	}
	return nil
}

func findLargestNodeNotLargerThanIndexLater(currentNode *Node, index int64) *Node {
	for currentNode != nil {
		if currentNode.Index == index {
			return currentNode
		}

		if currentNode.Index > index {
			if currentNode.Previous == nil {
				// !import: just one node
				return nil
			}
			return currentNode.Previous
		}

		if currentNode.Next == nil {
			// node is Tail
			return currentNode
		}
		currentNode = currentNode.Next
	}
	return nil
}

// FindSmallestNodeNotSmallerThanIndex find the smallest node not smaller than the index
// eg: 1 2 3 4 5 6 9
//    if index = 3, return node which index is 3
//    if index = 7, return node which index is 9
// eg: 2 (just one node)
//    if index = 3, return nil
func (node *Node) FindSmallestNodeNotSmallerThanIndex(index int64) *Node {
	if node == nil {
		return nil
	}

	switch {
	case node.Index == index:
		return node
	case node.Index > index:
		if node.Previous == nil {
			// node is Tail
			return node
		}
		return findSmallestNodeNotSmallerThanIndexForward(node.Previous, index)
	default:
		return findSmallestNodeNotSmallerThanIndexLater(node.Next, index)
	}
}

func findSmallestNodeNotSmallerThanIndexForward(currentNode *Node, index int64) *Node {
	for currentNode != nil {
		if currentNode.Index == index {
			return currentNode
		}

		if currentNode.Index < index {
			if currentNode.Next == nil {
				// !import: just one node
				return nil
			}
			return currentNode.Next
		}

		if currentNode.Previous == nil {
			// node is Head
			return currentNode
		}
		currentNode = currentNode.Previous
	}
	return nil
}

func findSmallestNodeNotSmallerThanIndexLater(currentNode *Node, index int64) *Node {
	for currentNode != nil {
		if currentNode.Index >= index {
			return currentNode
		}

		currentNode = currentNode.Next
	}
	return nil
}
