// pkg/database/models/pumpfun.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// PumpFunPool represents PumpFun's bonding curve pool structure
type PumpFunPool struct {
	BasePool
	BondingCurveType string    `gorm:"type:varchar(32);not null"`
	SpotPrice        float64   `gorm:"not null"`
	Delta            float64   `gorm:"not null"`
	Fee              float64   `gorm:"not null"`
	Supply           uint64    `gorm:"not null"`
	BaseDecimal      uint8     `gorm:"type:smallint"`
	QuoteDecimal     uint8     `gorm:"type:smallint"`
	LaunchTime       time.Time `gorm:"not null"`

	Metrics []PumpFunMetric `gorm:"foreignKey:PoolID"`
}

type PumpFunMetric struct {
	gorm.Model
	PoolID       string    `gorm:"type:varchar(64);not null;index"`
	Timestamp    time.Time `gorm:"index;not null"`
	CurrentPrice float64
	TotalSupply  uint64
	Volume24h    float64
	MarketCap    float64
}
