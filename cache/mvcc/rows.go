package mvcc

import "sync/atomic"

type Row struct {
	MultiVersionValue []SingleVersionData
	LastTransactionID int64
}

type SingleVersionData struct {
	Value         interface{}
	TransactionID int64
	IsDelete      bool
}

func (r *Row) Read(txid int64) (value interface{}, ok bool) {
	// if r.LastTransactionID <= txid {
	//     for _, sd := range r.MultiVersionValue {
	//         if sd.TransactionID == r.LastTransactionID {
	//             if !sd.IsDelete {
	//                 return nil, false
	//             }
	//         }
	//     }
	// }

	// filter txid with currentSnapshot
	return nil, false
}

func (r *Row) Put(txid int64, value interface{}) bool {
	r.MultiVersionValue = append(r.MultiVersionValue, SingleVersionData{
		Value:         value,
		TransactionID: txid,
	})

	if atomic.LoadInt64(&r.LastTransactionID) >= txid {
		return true
	}

	for {
		if atomic.CompareAndSwapInt64(&r.LastTransactionID, r.LastTransactionID, txid) {
			return true
		}
	}
}
