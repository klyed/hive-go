package types

import (
	"github.com/klyed/hive-go/encoding/transaction"
)

//ClaimRewardBalanceOperation represents claim_reward_balance operation data.
type ClaimRewardBalanceOperation struct {
	Account     string `json:"account"`
	RewardHive  *Asset `json:"reward_hive"`
	RewardHbd   *Asset `json:"reward_hbd"`
	RewardVests *Asset `json:"reward_vests"`
}

//Type function that defines the type of operation ClaimRewardBalanceOperation.
func (op *ClaimRewardBalanceOperation) Type() OpType {
	return TypeClaimRewardBalance
}

//Data returns the operation data ClaimRewardBalanceOperation.
func (op *ClaimRewardBalanceOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type ClaimRewardBalanceOperation to bytes.
func (op *ClaimRewardBalanceOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeClaimRewardBalance.Code()))
	enc.Encode(op.Account)
	enc.Encode(op.RewardHive)
	enc.Encode(op.RewardHbd)
	enc.Encode(op.RewardVests)
	return enc.Err()
}
