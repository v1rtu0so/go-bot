// pkg/database/models/base.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Protocol represents different DEX protocols
type Protocol string

const (
	ProtocolRaydium  Protocol = "RAYDIUM"
	ProtocolJupiter  Protocol = "JUPITER"
	ProtocolMeteora  Protocol = "METEORA"
	ProtocolMoonshot Protocol = "MOONSHOT"
	ProtocolPumpFun  Protocol = "PUMPFUN"
)

// PoolType represents different types of liquidity pools
type PoolType string

const (
	PoolTypeAMM          PoolType = "AMM"
	PoolTypeCLMM         PoolType = "CLMM"
	PoolTypeBondingCurve PoolType = "BONDING_CURVE"
)

// BasePool contains common fields for all pool types
type BasePool struct {
	gorm.Model
	ID          string    `gorm:"primaryKey;type:varchar(64)"`
	Protocol    Protocol  `gorm:"type:varchar(20);not null;index"`
	PoolType    PoolType  `gorm:"type:varchar(20);not null;index"`
	ProgramID   string    `gorm:"type:varchar(64);not null;index"`
	BaseMint    string    `gorm:"type:varchar(64);not null;index"`
	QuoteMint   string    `gorm:"type:varchar(64);not null;index"`
	LastUpdated time.Time `gorm:"index;not null"`
}
