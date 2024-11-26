// pkg/database/models/moonshot.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// MoonshotPool represents Moonshot's market making pool
type MoonshotPool struct {
	BasePool
	Authority     string `gorm:"type:varchar(64);not null"`
	BaseVault     string `gorm:"type:varchar(64);not null"`
	QuoteVault    string `gorm:"type:varchar(64);not null"`
	BaseDecimals  uint8  `gorm:"type:smallint"`
	QuoteDecimals uint8  `gorm:"type:smallint"`
	MMFeeBps      uint16 `gorm:"not null"` // Market maker fee in basis points
	VolumeFeeBps  uint16 `gorm:"not null"` // Volume-based fee in basis points
	BaseLotSize   uint64 `gorm:"not null"`
	QuoteLotSize  uint64 `gorm:"not null"`
	MinBaseRate   uint64 `gorm:"not null"`
	MaxBaseRate   uint64 `gorm:"not null"`

	Metrics []MoonshotMetric `gorm:"foreignKey:PoolID"`
}

// MoonshotMetric represents time-series metrics for Moonshot pools
type MoonshotMetric struct {
	gorm.Model
	PoolID        string    `gorm:"type:varchar(64);not null;index"`
	Timestamp     time.Time `gorm:"index;not null"`
	BaseReserve   float64
	QuoteReserve  float64
	Volume24h     float64
	TVL           float64
	APR           float64
	TotalFees     float64
	ActiveMakers  int32   // Number of active market makers
	AverageSpread float64 // Average bid-ask spread
}
