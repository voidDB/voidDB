package link

import (
	"time"

	"github.com/voidDB/voidDB/common"
)

type Metadata []byte

func NewMetadata(timestamp Timestamp, txnSerial int) (metadata Metadata) {
	metadata = make([]byte, 2*common.WordSize)

	metadata.setTimestamp(timestamp)

	metadata.setTxnSerial(txnSerial)

	return
}

func (metadata Metadata) Timestamp() Timestamp {
	return common.Field(metadata, 0, common.WordSize)
}

func (metadata Metadata) setTimestamp(timestamp Timestamp) {
	copy(
		metadata.Timestamp(),
		timestamp,
	)

	return
}

func (metadata Metadata) TxnSerial() TxnSerial {
	return common.Field(metadata, common.WordSize, common.WordSize)
}

func (metadata Metadata) setTxnSerial(i int) {
	common.PutInt(
		metadata.TxnSerial(),
		i,
	)

	return
}

type Timestamp []byte

func (timestamp Timestamp) Nanoseconds() int {
	return common.GetInt(timestamp)
}

func (timestamp Timestamp) Time() time.Time {
	return time.Unix(0,
		int64(timestamp.Nanoseconds()),
	)
}

type TxnSerial []byte

func (txnSerial TxnSerial) Int() int {
	return common.GetInt(txnSerial)
}
