// pkg/database/models/jupiter.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// JupiterPool represents Jupiter's aggregator pool structure
type JupiterPool struct {
	BasePool
	// Jupiter specific fields
	Label          string  `gorm:"type:varchar(128)"` // Human readable label
	InputToken     string  `gorm:"type:varchar(64);index"`
	OutputToken    string  `gorm:"type:varchar(64);index"`
	InAmount       uint64  `gorm:"not null"`
	OutAmount      uint64  `gorm:"not null"`
	FeeStructure   string  `gorm:"type:jsonb"`       // Fee structure as JSON
	RouteType      string  `gorm:"type:varchar(32)"` // Direct, Transitive, etc.
	PlatformFee    float64 `gorm:"type:decimal(6,4)"`
	InputDecimals  uint8   `gorm:"type:smallint"`
	OutputDecimals uint8   `gorm:"type:smallint"`

	Metrics []JupiterMetric `gorm:"foreignKey:PoolID"`
}

// JupiterMetric represents time-series metrics for Jupiter pools
type JupiterMetric struct {
	gorm.Model
	PoolID         string    `gorm:"type:varchar(64);not null;index"`
	Timestamp      time.Time `gorm:"index;not null"`
	Volume24h      float64
	TVL            float64
	PriceImpact    float64
	RouteScore     float64 // Jupiter's route scoring metric
	SuccessRate24h float64 // Success rate of swaps in last 24h
	AverageTxTime  float64 // Average transaction time in ms
}
