package skiplist

import "testing"

func TestNodeSet(t *testing.T) {
	var node *SortDoublyLinkedListNode
	node.Set(nil, 1, 1)
	setForward(nil, node, nil)
	setLater(nil, node, nil)
}

func TestNodeFind(t *testing.T) {
	var node *SortDoublyLinkedListNode
	node.FindNode(1)
	findForward(node, 1)
	findLaster(node, 1)
}

func TestNodeRemove(t *testing.T) {
	var node *SortDoublyLinkedListNode
	node.Remove(nil, 1)
	removeForward(nil, node, 1)
	removeLaster(nil, node, 1)
	removeCurrentNode(nil, node)
}

func TestNodeFindLargestNodeNotLargerThanIndex(t *testing.T) {
	var node *SortDoublyLinkedListNode
	node.FindLargestNodeNotLargerThanIndex(1)
	findLargestNodeNotLargerThanIndexForward(node, 1)
	findLargestNodeNotLargerThanIndexLater(node, 1)
}

func TestNodeFindSmallestNodeNotSmallerThanIndexForward(t *testing.T) {
	var node *SortDoublyLinkedListNode
	node.FindSmallestNodeNotSmallerThanIndex(1)
	findSmallestNodeNotSmallerThanIndexForward(node, 1)
	findSmallestNodeNotSmallerThanIndexLater(node, 1)
}
