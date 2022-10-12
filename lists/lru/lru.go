package lru

import "fmt"

type LRUCache struct {
	Filter map[int]struct{}
	Len    int
	Cap    int
	Tail   *Node
	Head   *Node
}

type Node struct {
	Previous *Node
	Next     *Node
	Data     *Data
}

type Data struct {
	Key   int
	Value int
}

func New(capacity int) *LRUCache {
	return &LRUCache{
		Filter: make(map[int]struct{}, capacity),
		Cap:    capacity,
	}
}

func (lc *LRUCache) Get(key int) int {
	// empty cache
	if lc.Len == 0 {
		return 0
	}

	// not Put this key
	if !lc.keyIsExist(key) {
		return 0
	}

	// exist
	pHead := lc.Head
	var foundData int
	for pHead != nil {
		if pHead.Data.Key == key {
			// found Node
			foundData = pHead.Data.Value

			// move currentNode
			lc.removeCurrentNode(pHead)
			lc.putNewValue(key, foundData)
			return foundData
		}
		pHead = pHead.Next
	}

	// will not do this
	// if do this, check filter logic
	return 0
}

func (lc *LRUCache) Put(key int, value int) {
	defer func() {
		lc.putKeyToFilter(key)
	}()

	if !lc.keyIsExist(key) {
		// key not exist
		lc.putNewValue(key, value)
		if lc.Len == lc.Cap {
			// cache full, remove tail node
			lc.removeCurrentNode(lc.Tail)
		} else {
			lc.Len++
		}
		return
	}

	// found old data
	pHead := lc.Head
	for pHead != nil {
		if pHead.Data.Key == key {
			// found, remove currentNode
			lc.removeCurrentNode(pHead)
			lc.putNewValue(key, value)
			return
		}

		pHead = pHead.Next
	}

	// will not do this
	// if do this, check filter logic
}

func (lc *LRUCache) keyIsExist(key int) bool {
	_, ok := lc.Filter[key]
	return ok
}

func (lc *LRUCache) putKeyToFilter(key int) {
	lc.Filter[key] = struct{}{}
}

func (lc *LRUCache) removeExistKey(key int) {
	delete(lc.Filter, key)
}

func (lc *LRUCache) putNewValue(key int, value int) {
	newNode := &Node{Data: &Data{Key: key, Value: value}}
	if lc.Head == nil {
		// first node
		lc.Head = newNode
		lc.Tail = newNode
		return
	}

	lc.Head.Previous = newNode
	newNode.Next = lc.Head
	lc.Head = newNode
}

func (lc *LRUCache) removeCurrentNode(curr *Node) {
	if curr.Previous == nil {
		// currNode is head
		lc.Head = curr.Next
	} else {
		curr.Previous.Next = curr.Next
	}

	if curr.Next == nil {
		// currnode is tail
		lc.Tail = curr.Previous
	} else {
		curr.Next.Previous = curr.Previous
	}
}

func (lc *LRUCache) Print() []string {
	result := []string{}
	pHead := lc.Head
	for pHead != nil {
		result = append(result, fmt.Sprintf("%d->%d", pHead.Data.Key, pHead.Data.Value))
		pHead = pHead.Next
	}
	return result
}
