package mvcc

import (
	"sync"
)

type TableBucket struct {
	L      sync.RWMutex
	Tables map[string]*Table
}

func (tbb *TableBucket) putValue(txid int64, tableName string, rowName string, value interface{}) {
	if tbb == nil {
		return
	}

	tbb.L.Lock()
	defer tbb.L.Unlock()

	table, ok := tbb.Tables[tableName]
	if !ok {
		table = NewTable(tableName)
		tbb.Tables[tableName] = table
	}
	table.putValue(txid, rowName, value)
}

func (tbb *TableBucket) foreachRows(txid int64, tableName string, h Handler) {
	if tbb == nil {
		return
	}

	tbb.L.Lock()
	defer tbb.L.Unlock()

	table, ok := tbb.Tables[tableName]
	if !ok {
		return
	}
	table.foreachRows(txid, h)
}

type Table struct {
	Name string
	Rows map[string]*Row
}

func NewTable(name string) *Table {
	return &Table{
		Name: name,
		Rows: map[string]*Row{},
	}
}

func (tb *Table) putValue(txid int64, rowName string, value interface{}) {
	if tb == nil {
		return
	}

	if len(tb.Rows) == 0 {
		tb.Rows = map[string]*Row{}
	}
	row, ok := tb.Rows[rowName]
	if !ok {
		row = NewRows(rowName)
		tb.Rows[rowName] = row
	}

	row.putValue(txid, value)
}

func (tb *Table) foreachRows(txid int64, h Handler) {
	if tb == nil {
		return
	}

	for _, row := range tb.Rows {
		iterator := row.findLargestNodeNotLargerThanTxID(txid)
		for {
			_, value, ok := iterator.GetValue()
			if !ok {
				break
			}
			singleVersionData, ok := value.(*SingleVersionData)
			if !ok {
				break
			}

			// judge this txid can read this value.
			if currentTxIDCanReadTxID(txid, singleVersionData.TransactionID) {
				// if this version is delete, break it.
				if !singleVersionData.IsDelete {
					h(singleVersionData.Value)
				}
				break
			}

			// can't read this version, read previous.
			iterator.Previous()
		}
	}
}
