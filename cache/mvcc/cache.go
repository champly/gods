package mvcc

func (mc *MVCCCache) PutValue(txid int64, tableName string, rowName string, value interface{}) {
	tableID := int(fnv32(tableName) % mc.BucketNum)
	tableBucket := mc.TableBuckets[tableID]

	tableBucket.L.Lock()
	defer tableBucket.L.Unlock()

	if _, ok := tableBucket.Tables[tableName]; !ok {
		tableBucket.Tables[tableName] = &Table{}
	}
	table := tableBucket.Tables[tableName]
	if _, ok := table.Rows[rowName]; !ok {
		table.Rows[rowName] = &Row{MultiVersionValue: []SingleVersionData{}}
	}

	table.Rows[rowName].Put(txid, value)
}

func (mc *MVCCCache) ForeachRows(txid int64, tableName string, h func(v interface{}) bool) {
	tableID := int(fnv32(tableName) % mc.BucketNum)
	tableBucket := mc.TableBuckets[tableID]

	tableBucket.L.RLock()
	defer tableBucket.L.RUnlock()

	if _, ok := tableBucket.Tables[tableName]; !ok {
		return
	}

	for _, row := range tableBucket.Tables[tableName].Rows {
		v, ok := row.Read(txid)
		if !ok {
			continue
		}

		if !h(v) {
			break
		}
	}
}
