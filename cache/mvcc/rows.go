package mvcc

import (
	"github.com/champly/gods/lists/skiplist"
)

var skiplistIndexLevel = 2

type Row struct {
	Name              string
	MultiVersionValue skiplist.ISkipList
}

type SingleVersionData struct {
	Value         interface{}
	TransactionID int64
	IsDelete      bool
}

func NewRows(name string) *Row {
	return &Row{
		Name:              name,
		MultiVersionValue: skiplist.NewSkipList(skiplistIndexLevel),
	}
}

func (r *Row) putValue(txid int64, value interface{}) {
	r.MultiVersionValue.Set(txid, &SingleVersionData{
		Value:         value,
		TransactionID: txid,
		IsDelete:      false,
	})
}

func (r *Row) findLargestNodeNotLargerThanTxID(txid int64) skiplist.IIterator {
	return r.MultiVersionValue.FindLargestNodeNotLargerThanIndex(txid)
}

func (r *Row) remove(txid int64) {
	r.MultiVersionValue.Set(txid, &SingleVersionData{
		TransactionID: txid,
		IsDelete:      true,
	})
}
