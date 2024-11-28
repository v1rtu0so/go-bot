package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	LastUpdated time.Time `gorm:"index;not null"`
}

type BaseAsset struct {
	BaseModel
	ID        string         `gorm:"primaryKey;type:varchar(64)"`
	Address   string         `gorm:"type:varchar(64);uniqueIndex"`
	Interface AssetInterface `gorm:"type:varchar(32);index"`
	Type      AssetType      `gorm:"type:varchar(20);index"`
	Status    AssetStatus    `gorm:"type:varchar(20);index"`
	Name      string         `gorm:"type:varchar(128)"`
	Symbol    string         `gorm:"type:varchar(32)"`
	Mint      string         `gorm:"type:varchar(64);index"`
	Decimals  uint8          `gorm:"type:smallint"`
	Supply    *uint64
}

type BasePool struct {
	BaseModel
	ID        string     `gorm:"primaryKey;type:varchar(64)"`
	Protocol  Protocol   `gorm:"type:varchar(20);index"`
	Type      PoolType   `gorm:"type:varchar(20);index"`
	Status    PoolStatus `gorm:"type:varchar(20);index"`
	ProgramID string     `gorm:"type:varchar(64);index"`
	BaseMint  string     `gorm:"type:varchar(64);index"`
	QuoteMint string     `gorm:"type:varchar(64);index"`
	Version   uint8      `gorm:"type:smallint"`
}

type BaseMetric struct {
	BaseModel
	PoolID      string    `gorm:"type:varchar(64);index"`
	Timestamp   time.Time `gorm:"index;not null"`
	LastFetched time.Time `gorm:"index"`
}
