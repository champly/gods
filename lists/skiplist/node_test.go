package skiplist

import "testing"

func TestNodeSet(t *testing.T) {
	var node *Node
	node.Set(nil, 1, 1)
	node.setForward(nil, nil)
	node.setLater(nil, nil)
}

func TestNodeFind(t *testing.T) {
	var node *Node
	node.FindNode(1)
	node.findForward(1)
	node.findLaster(1)
}

func TestNodeRemove(t *testing.T) {
	var node *Node
	node.Remove(nil, 1)
	node.removeForward(nil, 1)
	node.removeLaster(nil, 1)
	node.removeCurrentNode(nil)

	// list := buildSortDubboLinkedList(1, 2, 3, 4, 5, 6, 7)
}
