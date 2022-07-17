package skiplist

import "fmt"

type SortDoublyLinkedList struct {
	Head *SortDoublyLinkedListNode
	Tail *SortDoublyLinkedListNode
}

func (sdlList *SortDoublyLinkedList) Set(index int64, value interface{}) *SortDoublyLinkedListNode {
	return sdlList.Head.Set(sdlList, index, value)
}

func (sdlList *SortDoublyLinkedList) FindNode(index int64) (*SortDoublyLinkedListNode, bool) {
	return sdlList.Head.FindNode(index)
}

// just for debug
func (sdlList *SortDoublyLinkedList) Print() {
	pHead := sdlList.Head
	for pHead != nil {
		fmt.Printf("%d(%v)\t", pHead.Index, pHead.Value)
		pHead = pHead.Next
	}
	fmt.Println()
}

func (sdlList *SortDoublyLinkedList) Remove(index int64) {
	sdlList.Head.Remove(sdlList, index)
}

func (sdlList *SortDoublyLinkedList) FindLargestNodeNotLargerThanIndex(index int64) *SortDoublyLinkedListNode {
	return sdlList.Head.FindLargestNodeNotLargerThanIndex(index)
}

func (sdlList *SortDoublyLinkedList) FindSmallestNodeNotSmallerThanIndex(index int64) *SortDoublyLinkedListNode {
	return sdlList.Head.FindSmallestNodeNotSmallerThanIndex(index)
}
