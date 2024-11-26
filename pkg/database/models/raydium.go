// pkg/database/models/raydium.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// RaydiumAMMPool represents Raydium's AMM pool structure
type RaydiumAMMPool struct {
	BasePool
	Version      uint8  `gorm:"type:smallint"`
	LpMint       string `gorm:"type:varchar(64);not null"`
	BaseVault    string `gorm:"type:varchar(64);not null"`
	QuoteVault   string `gorm:"type:varchar(64);not null"`
	LpVault      string `gorm:"type:varchar(64);not null"`
	OpenOrders   string `gorm:"type:varchar(64)"`
	TargetOrders string `gorm:"type:varchar(64)"`
	Authority    string `gorm:"type:varchar(64);not null"`
	BaseDecimal  uint8  `gorm:"type:smallint"`
	QuoteDecimal uint8  `gorm:"type:smallint"`
	LpDecimal    uint8  `gorm:"type:smallint"`

	// Relationship with metrics
	Metrics []RaydiumAMMMetric `gorm:"foreignKey:PoolID"`
}

// RaydiumCLMMPool represents Raydium's Concentrated Liquidity pool
type RaydiumCLMMPool struct {
	BasePool
	Tick         int64  `gorm:"not null"`
	SqrtPrice    string `gorm:"type:varchar(64);not null"`
	FeeRate      uint32 `gorm:"not null"`
	BaseDecimal  uint8  `gorm:"type:smallint"`
	QuoteDecimal uint8  `gorm:"type:smallint"`
	TickSpacing  uint16 `gorm:"not null"`
	BaseReserve  string `gorm:"type:varchar(64)"`
	QuoteReserve string `gorm:"type:varchar(64)"`

	Metrics []RaydiumCLMMMetric `gorm:"foreignKey:PoolID"`
}

// RaydiumAMMMetric represents time-series metrics for AMM pools
type RaydiumAMMMetric struct {
	gorm.Model
	PoolID       string    `gorm:"type:varchar(64);not null;index"`
	Timestamp    time.Time `gorm:"index;not null"`
	BaseReserve  float64
	QuoteReserve float64
	LpSupply     float64
	Volume24h    float64
	TVL          float64
	APR          float64
}

// RaydiumCLMMMetric represents time-series metrics for CLMM pools
type RaydiumCLMMMetric struct {
	gorm.Model
	PoolID      string    `gorm:"type:varchar(64);not null;index"`
	Timestamp   time.Time `gorm:"index;not null"`
	CurrentTick int64
	SqrtPrice   string
	Liquidity   string
	Volume24h   float64
	TVL         float64
	APR         float64
}
