package skiplist

import "fmt"

type SortDubboLinkedList struct {
	Head *Node
	Tail *Node
}

func (sdlList *SortDubboLinkedList) Set(index int64, value interface{}) {
	sdlList.Head.Set(sdlList, index, value)
}

func (sdlList *SortDubboLinkedList) FindNode(index int64) (*Node, bool) {
	return sdlList.Head.FindNode(index)
}

// just for debug
func (sdlList *SortDubboLinkedList) Print() {
	pHead := sdlList.Head
	for pHead != nil {
		fmt.Printf("%d(%v)\t", pHead.Index, pHead.Value)
		pHead = pHead.Next
	}
	fmt.Println()
}

func (sdlList *SortDubboLinkedList) Remove(index int64) {
	sdlList.Head.Remove(sdlList, index)
}

func (sdlList *SortDubboLinkedList) FindLargestNodeNotLargerThanIndex(index int64) *Node {
	return sdlList.Head.FindLargestNodeNotLargerThanIndex(index)
}

func (sdlList *SortDubboLinkedList) FindSmallestNodeNotSmallerThanIndex(index int64) *Node {
	return sdlList.Head.FindSmallestNodeNotSmallerThanIndex(index)
}
