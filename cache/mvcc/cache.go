package mvcc

type Handler func(v interface{})

type MVCCCache struct {
	Name         string
	BucketNum    uint32
	TableBuckets []*TableBucket
}

func (mc *MVCCCache) PutValue(txid int64, tableName string, rowName string, value interface{}) {
	tableBucket := mc.getTableBucket(tableName)
	tableBucket.putValue(txid, tableName, rowName, value)
}

func (mc *MVCCCache) ForeachRows(txid int64, tableName string, h Handler) {
	tableBucket := mc.getTableBucket(tableName)
	tableBucket.foreachRows(txid, tableName, h)
}

func (mc *MVCCCache) getTableBucket(tableName string) *TableBucket {
	index := int(fnv32(tableName) % mc.BucketNum)
	return mc.TableBuckets[index]
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
