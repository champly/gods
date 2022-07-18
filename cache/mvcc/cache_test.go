package mvcc

import "testing"

func TestCache(t *testing.T) {
	cache := newMVCCCache("demo")
	tableName1 := "t1"
	rowName1 := "r1"
	cache.PutValueOnce(tableName1, rowName1, 1)
	cache.PutValueOnce(tableName1, rowName1, 2)
	cache.PutValueOnce(tableName1, rowName1, 3)

	count := 0
	cache.ForeachRowsOnce(tableName1, func(v interface{}) {
		data := v.(int)
		count += data
	})
	if count != 3 {
		t.Errorf("count should be %d, but got %d", 3, count)
	}

	// multi version control
	txid := GetTxID()
	cache.PutValueOnce(tableName1, rowName1, 5)

	// old txid
	count = 0
	cache.ForeachRows(txid, tableName1, func(v interface{}) {
		data := v.(int)
		count += data
	})
	if count != 3 {
		t.Errorf("count should be %d, but got %d", 3, count)
	}
	PutTxID(txid)

	// last txid
	count = 0
	cache.ForeachRowsOnce(tableName1, func(v interface{}) {
		data := v.(int)
		count += data
	})
	if count != 5 {
		t.Errorf("count should be %d, but got %d", 5, count)
	}

	// multi row and multi version, r1=>5, r2 is empty
	rowName2 := "r2"
	txid = GetTxID()
	// r1=>6, r2=>2
	cache.PutValueOnce(tableName1, rowName1, 6)
	cache.PutValueOnce(tableName1, rowName2, 2)
	// old txid
	count = 0
	time := 0
	cache.ForeachRows(txid, tableName1, func(v interface{}) {
		data := v.(int)
		count += data
		time++
	})
	if count != 5 || time != 1 {
		t.Errorf("count should be %d, but got %d", 5, count)
		t.Errorf("should found %d row, but got %d", 1, time)
	}

	// last txid
	count = 0
	time = 0
	cache.ForeachRowsOnce(tableName1, func(v interface{}) {
		data := v.(int)
		count += data
		time++
	})
	if count != 8 || time != 2 {
		t.Errorf("count should be %d, but got %d", 8, count)
		t.Errorf("should found %d row, but got %d", 2, time)
	}

	// remove r1
	cache.RemoveRowOnce(tableName1, rowName1)
	count = 0
	time = 0
	cache.ForeachRowsOnce(tableName1, func(v interface{}) {
		data := v.(int)
		count += data
		time++
	})
	if count != 2 || time != 1 {
		t.Errorf("count should be %d, but got %d", 2, count)
		t.Errorf("should found %d row, but got %d", 1, time)
	}
	// old txid
	count = 0
	time = 0
	cache.ForeachRows(txid, tableName1, func(v interface{}) {
		data := v.(int)
		count += data
		time++
	})
	if count != 5 || time != 1 {
		t.Errorf("count should be %d, but got %d", 5, count)
		t.Errorf("should found %d row, but got %d", 1, time)
	}
	// remove r2
	cache.RemoveRowOnce(tableName1, rowName2)
	count = 0
	time = 0
	cache.ForeachRowsOnce(tableName1, func(v interface{}) {
		data := v.(int)
		count += data
		time++
	})
	if count != 0 || time != 0 {
		t.Errorf("count should be %d, but got %d", 0, count)
		t.Errorf("should found %d row, but got %d", 0, time)
	}
	// old txid
	count = 0
	time = 0
	cache.ForeachRows(txid, tableName1, func(v interface{}) {
		data := v.(int)
		count += data
		time++
	})
	if count != 5 || time != 1 {
		t.Errorf("count should be %d, but got %d", 5, count)
		t.Errorf("should found %d row, but got %d", 1, time)
	}
}
