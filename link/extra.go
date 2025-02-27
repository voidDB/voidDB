package link

import (
	"time"

	"github.com/voidDB/voidDB/common"
)

type Metadata []byte

func (metadata Metadata) Timestamp() Timestamp {
	return common.Field(metadata, 0, common.WordSize)
}

func (metadata Metadata) SetTimestamp(timestamp Timestamp) {
	copy(
		metadata.Timestamp(),
		timestamp,
	)

	return
}

func (metadata Metadata) TxnSerial() TxnSerial {
	return common.Field(metadata, common.WordSize, common.WordSize)
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
