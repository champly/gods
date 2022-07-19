package mvcc

var (
	defaultBucketCount uint32 = 32
)

type Handler func(v interface{})

type IMVCCCache interface {
	PutValueOnce(tableName string, rowName string, value interface{})
	PutValueWithTxID(txid int64, tableName string, rowName string, value interface{})
	ForeachRowsOnce(tableName string, h Handler)
	ForeachRows(txid int64, tableName string, h Handler)
	RemoveRowOnce(tableName, rowName string)
	RemoveRow(txid int64, tableName string, rowName string)
}

func New(name string) IMVCCCache {
	return newMVCCCache(name)
}

type MVCCCache struct {
	Name         string
	BucketNum    uint32
	TableBuckets []*TableBucket
}

func newMVCCCache(name string) *MVCCCache {
	cache := &MVCCCache{
		Name:         name,
		BucketNum:    defaultBucketCount,
		TableBuckets: make([]*TableBucket, 0, defaultBucketCount),
	}
	for i := 0; i < int(defaultBucketCount); i++ {
		cache.TableBuckets = append(cache.TableBuckets, &TableBucket{Tables: make(map[string]*Table)})
	}

	return cache
}

func (mc *MVCCCache) PutValueOnce(tableName string, rowName string, value interface{}) {
	txid := GetTxID()
	defer func() {
		PutTxID(txid)
	}()
	mc.PutValueWithTxID(txid, tableName, rowName, value)
}

func (mc *MVCCCache) PutValueWithTxID(txid int64, tableName string, rowName string, value interface{}) {
	tableBucket := mc.getTableBucket(tableName)
	tableBucket.putValue(txid, tableName, rowName, value)
}

func (mc *MVCCCache) ForeachRowsOnce(tableName string, h Handler) {
	txid := GetTxID()
	defer func() {
		PutTxID(txid)
	}()
	mc.ForeachRows(txid, tableName, h)
}

func (mc *MVCCCache) ForeachRows(txid int64, tableName string, h Handler) {
	tableBucket := mc.getTableBucket(tableName)
	tableBucket.foreachRows(txid, tableName, h)
}

func (mc *MVCCCache) RemoveRowOnce(tableName, rowName string) {
	txid := GetTxID()
	defer func() {
		PutTxID(txid)
	}()

	mc.RemoveRow(txid, tableName, rowName)
}

func (mc *MVCCCache) RemoveRow(txid int64, tableName string, rowName string) {
	tableBucket := mc.getTableBucket(tableName)
	tableBucket.remove(txid, tableName, rowName)
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
