package models

// Pool represents a unified pool model for all DEX types
type Pool struct {
	BasePool
	Config      PoolConfig  `gorm:"type:jsonb"`
	MarketState MarketState `gorm:"type:jsonb"`
	Authority   string      `gorm:"type:varchar(64)"`
	BaseVault   string      `gorm:"type:varchar(64)"`
	QuoteVault  string      `gorm:"type:varchar(64)"`
	LpMint      *string     `gorm:"type:varchar(64)"`
	LpVault     *string     `gorm:"type:varchar(64)"`

	// Protocol-specific extensions
	RaydiumData   *RaydiumPoolData `gorm:"type:jsonb"`
	WhirlpoolData *WhirlpoolData   `gorm:"type:jsonb"`
	JupiterData   *JupiterPoolData `gorm:"type:jsonb"`
	MeteoraData   *MeteoraPoolData `gorm:"type:jsonb"`
	PumpFunData   *PumpFunPoolData `gorm:"type:jsonb"`

	// Relationships
	Assets  []Asset      `gorm:"many2many:asset_pools;"`
	Metrics []PoolMetric `gorm:"foreignKey:PoolID"`
}

type PoolConfig struct {
	Decimals       uint8    `json:"decimals"`
	FeeRate        uint32   `json:"fee_rate"`
	TickSpacing    *int32   `json:"tick_spacing,omitempty"`
	TradeFeeBps    *uint16  `json:"trade_fee_bps,omitempty"`
	ProtocolFeeBps *uint16  `json:"protocol_fee_bps,omitempty"`
	InitSqrtPrice  *string  `json:"init_sqrt_price,omitempty"`
	InitTick       *int64   `json:"init_tick,omitempty"`
	TokenMintA     *string  `json:"token_mint_a,omitempty"`
	TokenMintB     *string  `json:"token_mint_b,omitempty"`
	Delta          *float64 `json:"delta,omitempty"`
	SpotPrice      *float64 `json:"spot_price,omitempty"`
	FeatureFlags   []string `json:"feature_flags,omitempty"`
	Permissions    []string `json:"permissions,omitempty"`
}

type MarketState struct {
	BaseReserve     float64  `json:"base_reserve"`
	QuoteReserve    float64  `json:"quote_reserve"`
	LpSupply        *float64 `json:"lp_supply,omitempty"`
	Price           float64  `json:"price"`
	Volume24h       float64  `json:"volume_24h"`
	TVL             float64  `json:"tvl"`
	APR             float64  `json:"apr"`
	Liquidity       *string  `json:"liquidity,omitempty"`
	CurrentTick     *int64   `json:"current_tick,omitempty"`
	SqrtPrice       *string  `json:"sqrt_price,omitempty"`
	UtilizationRate float64  `json:"utilization_rate"`
	SwapCount       uint64   `json:"swap_count"`
}

type RaydiumPoolData struct {
	OpenOrders    string `json:"open_orders"`
	TargetOrders  string `json:"target_orders"`
	WithdrawQueue string `json:"withdraw_queue"`
	BaseDecimal   uint8  `json:"base_decimal"`
	QuoteDecimal  uint8  `json:"quote_decimal"`
	LpDecimal     uint8  `json:"lp_decimal"`
}

type WhirlpoolData struct {
	RewardVaultA   string `json:"reward_vault_a"`
	RewardVaultB   string `json:"reward_vault_b"`
	Protocol       string `json:"protocol"`
	TokenADecimals uint8  `json:"token_a_decimals"`
	TokenBDecimals uint8  `json:"token_b_decimals"`
}

type JupiterPoolData struct {
	RouteType      string `json:"route_type"`
	InputMint      string `json:"input_mint"`
	OutputMint     string `json:"output_mint"`
	InAmount       uint64 `json:"in_amount"`
	OutAmount      uint64 `json:"out_amount"`
	InputDecimals  uint8  `json:"input_decimals"`
	OutputDecimals uint8  `json:"output_decimals"`
}

type MeteoraPoolData struct {
	SwapEnabled    bool   `json:"swap_enabled"`
	LiquidityFee   uint32 `json:"liquidity_fee"`
	ProtocolFee    uint32 `json:"protocol_fee"`
	FeeGrowthBase  string `json:"fee_growth_base"`
	FeeGrowthQuote string `json:"fee_growth_quote"`
	BaseDecimals   uint8  `json:"base_decimals"`
	QuoteDecimals  uint8  `json:"quote_decimals"`
}

type PumpFunPoolData struct {
	BondingCurveType string  `json:"bonding_curve_type"`
	Delta            float64 `json:"delta"`
	Fee              float64 `json:"fee"`
	Supply           uint64  `json:"supply"`
	BaseDecimal      uint8   `json:"base_decimal"`
	QuoteDecimal     uint8   `json:"quote_decimal"`
}
