package store

import (
	"time"

	"github.com/elastos/Elastos.ELA.SPV/database"
	"github.com/elastos/Elastos.ELA.SPV/sdk"
	"github.com/elastos/Elastos.ELA.SPV/util"
	"github.com/elastos/Elastos.ELA/common"

	"github.com/syndtr/goleveldb/leveldb"
)

type HeaderStore interface {
	database.Headers
	GetByHeight(height uint32) (header *util.Header, err error)
}

type DataStore interface {
	database.DB
	Addrs() Addrs
	Txs() Txs
	Ops() Ops
	Que() Que
	Arbiters() Arbiters
	CID() CustomID
	Batch() DataBatch
}

type DataBatch interface {
	batch
	Txs() TxsBatch
	Ops() OpsBatch
	Que() QueBatch
	GetNakedBatch() *leveldb.Batch
	// Delete all transactions, ops, queued items on
	// the given height.
	DelAll(height uint32) error
}

type batch interface {
	Rollback() error
	Commit() error
}

type Addrs interface {
	database.DB
	GetFilter() *sdk.AddrFilter
	Put(addr *common.Uint168) error
	GetAll() []*common.Uint168
}

type Txs interface {
	database.DB
	Put(tx *util.Tx) error
	Get(txId *common.Uint256) (*util.Tx, error)
	GetAll() ([]*util.Tx, error)
	GetIds(height uint32) ([]*common.Uint256, error)
	PutForkTxs(txs []*util.Tx, hash *common.Uint256) error
	GetForkTxs(hash *common.Uint256) ([]*util.Tx, error)
	Del(txId *common.Uint256) error
	Batch() TxsBatch
}

type TxsBatch interface {
	batch
	Put(tx *util.Tx) error
	Del(txId *common.Uint256) error
	DelAll(height uint32) error
}

type Ops interface {
	database.DB
	Put(*util.OutPoint, common.Uint168) error
	HaveOp(*util.OutPoint) *common.Uint168
	GetAll() ([]*util.OutPoint, error)
	Batch() OpsBatch
}

type OpsBatch interface {
	batch
	Put(*util.OutPoint, common.Uint168) error
	Del(*util.OutPoint) error
}

type Que interface {
	database.DB

	// Put a queue item to database
	Put(item *QueItem) error

	// Get all items in queue
	GetAll() ([]*QueItem, error)

	// Delete confirmed item in queue
	Del(notifyId, txHash *common.Uint256) error

	// Batch returns a queue batch instance.
	Batch() QueBatch
}

type QueBatch interface {
	batch

	// Put a queue item to database
	Put(item *QueItem) error

	// Delete confirmed item in queue
	Del(notifyId, txHash *common.Uint256) error

	// Delete all items on the given height.
	DelAll(height uint32) error
}

type QueItem struct {
	NotifyId   common.Uint256
	TxId       common.Uint256
	Height     uint32
	LastNotify time.Time
}

type Arbiters interface {
	database.DB
	Put(height uint32, crcArbiters [][]byte, normalArbiters [][]byte) error
	BatchPut(height uint32, crcArbiters [][]byte, normalArbiters [][]byte, batch *leveldb.Batch) error
	Get() (crcArbiters [][]byte, normalArbiters [][]byte, err error)
	GetByHeight(height uint32) (crcArbiters [][]byte, normalArbiters [][]byte, err error)
}

type CustomID interface {
	database.DB
	PutControversialReservedCustomIDs(reservedCustomIDs []string) error
	BatchPutControversialReservedCustomIDs(reservedCustomIDs []string, batch *leveldb.Batch) error

	PutControversialReceivedCustomIDs(reservedCustomIDs []string, did common.Uint168) error
	BatchPutControversialReceivedCustomIDs(receivedCustomIDs []string, did common.Uint168, batch *leveldb.Batch) error

	PutRChangeCustomIDFee(rate common.Fixed64) error
	BatchPutChangeCustomIDFee(rate common.Fixed64, batch *leveldb.Batch) error

	GetControversialReservedCustomIDs() (map[string]struct{}, error)
	GetControversialReceivedCustomIDs() (map[string]common.Uint168, error)
	GetCustomIDFeeRate() (common.Fixed64, error)
}
