package link

import (
	"time"

	"github.com/voidDB/voidDB/common"
)

type Metadata []byte

func NewMetadata(timestamp Timestamp, txnSerial int) (metadata Metadata) {
	metadata = common.NewTwin()

	metadata.setTimestamp(timestamp)

	metadata.setTxnSerial(txnSerial)

	return
}

func (metadata Metadata) Timestamp() Timestamp {
	return common.WordN(metadata, 0)
}

func (metadata Metadata) setTimestamp(timestamp Timestamp) {
	copy(
		metadata.Timestamp(),
		timestamp,
	)

	return
}

func (metadata Metadata) TxnSerial() TxnSerial {
	return common.WordN(metadata, 1)
}

func (metadata Metadata) setTxnSerial(i int) {
	common.PutIntIntoWord(
		metadata.TxnSerial(),
		i,
	)

	return
}

type Timestamp []byte

func (timestamp Timestamp) Nanoseconds() int {
	return common.GetIntFromWord(timestamp)
}

func (timestamp Timestamp) Time() time.Time {
	return time.Unix(0,
		int64(timestamp.Nanoseconds()),
	)
}

type TxnSerial []byte

func (txnSerial TxnSerial) Int() int {
	return common.GetIntFromWord(txnSerial)
}
