// pkg/database/models/orca.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// OrcaWhirlpool represents Orca's concentrated liquidity pool (Whirlpool)
type OrcaWhirlpool struct {
	BasePool
	WhirlpoolsConfig string `gorm:"type:varchar(64);not null"`
	TokenVaultA      string `gorm:"type:varchar(64);not null"`
	TokenVaultB      string `gorm:"type:varchar(64);not null"`
	TickSpacing      int32  `gorm:"not null"`
	TickArrayBits    uint16 `gorm:"not null"`
	FeeRate          uint16 `gorm:"not null"`
	ProtocolFeeRate  uint16 `gorm:"not null"`
	TokenADecimals   uint8  `gorm:"type:smallint"`
	TokenBDecimals   uint8  `gorm:"type:smallint"`
	SqrtPrice        string `gorm:"type:varchar(64)"`
	TickCurrentIndex int32  `gorm:"not null"`
	LiquidityState   string `gorm:"type:jsonb"` // Current liquidity state
	FeeGrowthGlobal  string `gorm:"type:varchar(128)"`

	Metrics []OrcaMetric `gorm:"foreignKey:PoolID"`
}

// OrcaMetric represents time-series metrics for Orca Whirlpools
type OrcaMetric struct {
	gorm.Model
	PoolID          string    `gorm:"type:varchar(64);not null;index"`
	Timestamp       time.Time `gorm:"index;not null"`
	TokenAReserve   float64
	TokenBReserve   float64
	Volume24h       float64
	TVL             float64
	APR             float64
	FeesGenerated   float64
	ActivePositions int32
	Liquidity       string
	PriceRange      string  `gorm:"type:jsonb"` // Current price range
	UtilizationRate float64 // Pool utilization rate
}
