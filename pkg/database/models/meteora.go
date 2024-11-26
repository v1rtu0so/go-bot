// pkg/database/models/meteora.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// MeteoraPool represents Meteora's concentrated liquidity pool
type MeteoraPool struct {
	BasePool
	TickSpacing      int32  `gorm:"not null"`
	FeeRate          uint32 `gorm:"not null"`
	LiquidityToken   string `gorm:"type:varchar(64);not null"`
	BaseDecimals     uint8  `gorm:"type:smallint"`
	QuoteDecimals    uint8  `gorm:"type:smallint"`
	CurrentSqrtPrice string `gorm:"type:varchar(64)"`
	CurrentTick      int32  `gorm:"not null"`
	FeeGrowthGlobal  string `gorm:"type:varchar(64)"`
	Protocol         string `gorm:"type:varchar(32)"`
	SwapEnabled      bool   `gorm:"default:true"`

	Metrics []MeteoraMetric `gorm:"foreignKey:PoolID"`
}

// MeteoraMetric represents time-series metrics for Meteora pools
type MeteoraMetric struct {
	gorm.Model
	PoolID         string    `gorm:"type:varchar(64);not null;index"`
	Timestamp      time.Time `gorm:"index;not null"`
	BaseReserve    float64
	QuoteReserve   float64
	Volume24h      float64
	TVL            float64
	APR            float64
	Liquidity      string
	ActivePosition int32   // Number of active positions
	FeesGenerated  float64 // Fees generated in last 24h
}
