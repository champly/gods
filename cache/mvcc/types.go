package mvcc

import (
	"sync"
)

type MVCCCache struct {
	Name         string
	BucketNum    uint32
	TableBuckets []*TableBucket
}

type TableBucket struct {
	L      sync.RWMutex
	Tables map[string]*Table
}

type Table struct {
	Rows map[string]*Row
}
