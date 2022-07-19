package mvcc

import (
	"log"
	"sync"
	"sync/atomic"
)

var (
	txIDLock                         sync.Mutex
	lastIndexID                      int64
	processingTransactionInformation []TransactionIDInfo
)

type TransactionIDInfo struct {
	ID       int64
	SnapShot []int64
	Expired  bool
}

func GetTxID() int64 {
	for {
		oldid := atomic.LoadInt64(&lastIndexID)
		newid := oldid + 1
		if atomic.CompareAndSwapInt64(&lastIndexID, oldid, newid) {
			putTxIDToCache(newid)
			return newid
		}
	}
}

func PutTxID(txid int64) {
	removeTxIDFromCache(txid)
}

func currentTxIDCanReadTxID(currentTxID, txID int64) bool {
	if currentTxID == txID {
		return true
	}

	txidInfo, ok := getTxInfo(currentTxID)
	if !ok {
		return false
	}

	for _, snapshotTxID := range txidInfo.SnapShot {
		if snapshotTxID == txID {
			return false
		}
	}
	return true
}

func getTxInfo(txid int64) (*TransactionIDInfo, bool) {
	txIDLock.Lock()
	defer txIDLock.Unlock()

	for i, info := range processingTransactionInformation {
		if !info.Expired && info.ID == txid {
			return &processingTransactionInformation[i], true
		}
	}
	return nil, false
}

func putTxIDToCache(txid int64) {
	txIDLock.Lock()
	defer txIDLock.Unlock()

	foundOldSample := false
	oldSampleIndex := 0
	for i, info := range processingTransactionInformation {
		if info.Expired {
			foundOldSample = true
			oldSampleIndex = i
			break
		}
	}

	if foundOldSample {
		processingTransactionInformation[oldSampleIndex] = TransactionIDInfo{
			ID:       txid,
			SnapShot: getProcessingTxIDSnapshot(),
			Expired:  false,
		}
	} else {
		// processingTransactionInformation is full
		processingTransactionInformation = append(processingTransactionInformation, TransactionIDInfo{
			ID:       txid,
			SnapShot: getProcessingTxIDSnapshot(),
			Expired:  false,
		})
	}
}

func removeTxIDFromCache(txid int64) {
	txIDLock.Lock()
	defer txIDLock.Unlock()

	for i := range processingTransactionInformation {
		if processingTransactionInformation[i].ID == txid {
			processingTransactionInformation[i].Expired = true
			return
		}
	}
	log.Printf("Put an error txid %d, or repeat put.", txid)
}

func getProcessingTxIDSnapshot() []int64 {
	snapshot := []int64{}
	for _, info := range processingTransactionInformation {
		if !info.Expired {
			snapshot = append(snapshot, info.ID)
		}
	}
	return snapshot
}
