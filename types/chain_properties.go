package types

import (
	"github.com/klyed/hive-go/encoding/transaction"
)

//ChainProperties is an additional structure used by other structures.
type ChainProperties struct {
	AccountCreationFee *Asset `json:"account_creation_fee"`
	MaximumBlockSize   uint32 `json:"maximum_block_size"`
	HbdInterestRate    uint16 `json:"hbd_interest_rate"`
}

//MarshalTransaction is a function of converting type ChainProperties to bytes.
func (cp *ChainProperties) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(cp.AccountCreationFee)
	enc.Encode(cp.MaximumBlockSize)
	enc.Encode(cp.HbdInterestRate)
	return enc.Err()
}
