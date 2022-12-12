package models

import "time"

type StakeAndReward struct {
	Rewards []apiRewards `json:"rewards"`
	Stake   apiStake     `json:"delegation_response"`
}

type apiRewards struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type apiStake struct {
	Delegation apiDelegation `json:"delegation"`
	Balance    apiBalance    `json:"balance"`
}

type apiDelegation struct {
	DAddr  string `json:"delegator_address"`
	VAddr  string `json:"validator_address"`
	Shares string `json:"shares"`
}

type apiBalance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type TxFetch struct {
	Tx          Tx          `json:"tx"`
	TxFetchTime TxFetchTime `json:"tx_response"`
}

type TxFetchTime struct {
	Timestamp time.Time `json:"timestamp"`
}

type Tx struct {
	Body TxBody `json:"body"`
}

type msg struct {
	DelegatorAddr string     `json:"delegator_address"`
	Amount        apiBalance `json:"amount"`
}

type TxBody struct {
	Body []msg `json:"messages"`
}
