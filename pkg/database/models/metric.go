package models

type AssetMetric struct {
	BaseMetric
	AssetID        string `gorm:"type:varchar(64);index"`
	Price          float64
	MarketCap      float64
	Volume24h      float64
	TotalSupply    uint64
	HolderCount    uint32
	Liquidity      float64
	TradeCount     uint32
	TradedAmount   float64
	AveragePrice   float64
	MaxPrice24h    float64
	MinPrice24h    float64
	PriceChange24h float64

	Asset Asset `gorm:"foreignKey:AssetID"`
}

type PoolMetric struct {
	BaseMetric
	PoolID           string `gorm:"type:varchar(64);index"`
	BaseReserve      float64
	QuoteReserve     float64
	Price            float64
	Volume24h        float64
	TVL              float64
	APR              float64
	TradeCount       uint32
	SwapEvents       uint32
	PriceImpact      float64
	Liquidity        *string
	FeesGenerated    float64
	UtilizationRate  float64
	UniqueTraders    uint32
	AverageTradeSize float64
	MaxTradeSize24h  float64
	MinTradeSize24h  float64
	Volatility24h    float64

	Pool Pool `gorm:"foreignKey:PoolID"`
}

type TokenMetric struct {
	BaseMetric
	TokenID            string `gorm:"type:varchar(64);index"`
	CirculatingSupply  uint64
	TotalSupply        uint64
	HolderCount        uint32
	ActiveHolders      uint32
	TransferCount      uint32
	UniqueReceivers    uint32
	AverageHoldTime    float64
	VelocityDaily      float64
	ConcentrationTop10 float64

	Token Asset `gorm:"foreignKey:TokenID"`
}

type MarketMetric struct {
	BaseMetric
	TotalVolume     float64
	TotalTVL        float64
	UniqueUsers     uint32
	ActivePools     uint32
	FailedTxCount   uint32
	GasUsed         float64
	AverageSlippage float64
	RoutingVolume   float64
}
