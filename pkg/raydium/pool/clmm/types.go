package clmm

// RaydiumClmmPool represents the structure of a CLMM pool.
type RaydiumClmmPool struct {
	ID                 string                     `json:"id"`
	MintProgramIDA     string                     `json:"mintProgramIdA"`
	MintProgramIDB     string                     `json:"mintProgramIdB"`
	MintA              string                     `json:"mintA"`
	MintB              string                     `json:"mintB"`
	VaultA             string                     `json:"vaultA"`
	VaultB             string                     `json:"vaultB"`
	MintDecimalsA      int                        `json:"mintDecimalsA"`
	MintDecimalsB      int                        `json:"mintDecimalsB"`
	AmmConfig          ApiClmmConfigurationItem   `json:"ammConfig"`
	RewardInfos        []RewardInfo               `json:"rewardInfos"`
	TVL                int                        `json:"tvl"`
	Day                *ApiClmmPoolsItemStatistic `json:"day,omitempty"`
	Week               *ApiClmmPoolsItemStatistic `json:"week,omitempty"`
	Month              *ApiClmmPoolsItemStatistic `json:"month,omitempty"`
	LookupTableAccount string                     `json:"lookupTableAccount"`
}

// ApiClmmConfigurationItem represents configuration details of a CLMM pool.
type ApiClmmConfigurationItem struct {
	ID              string `json:"id"`
	Index           int    `json:"index"`
	ProtocolFeeRate int    `json:"protocolFeeRate"`
	TradeFeeRate    int    `json:"tradeFeeRate"`
	TickSpacing     int    `json:"tickSpacing"`
	FundFeeRate     int    `json:"fundFeeRate"`
	FundOwner       string `json:"fundOwner"`
}

// ApiClmmPoolsItemStatistic represents aggregated statistics of a CLMM pool.
type ApiClmmPoolsItemStatistic struct {
	ID         int `json:"id"`
	Volume     int `json:"volume"`
	VolumeFee  int `json:"volumeFee"`
	FeeA       int `json:"feeA"`
	FeeB       int `json:"feeB"`
	FeeApr     int `json:"feeApr"`
	RewardAprA int `json:"rewardAprA"`
	RewardAprB int `json:"rewardAprB"`
	RewardAprC int `json:"rewardAprC"`
	Apr        int `json:"apr"`
	PriceMin   int `json:"priceMin"`
	PriceMax   int `json:"priceMax"`
}

// RewardInfo represents reward information for a CLMM pool.
type RewardInfo struct {
	ID        int    `json:"id"`
	Mint      string `json:"mint"`
	ProgramID string `json:"programId"`
}
