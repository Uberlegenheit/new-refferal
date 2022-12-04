package models

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