package models

type AssetRelationship struct {
	BaseModel
	SourceID     string       `gorm:"type:varchar(64);index"`
	TargetID     string       `gorm:"type:varchar(64);index"`
	RelationType RelationType `gorm:"type:varchar(32)"`
	Protocol     Protocol     `gorm:"type:varchar(20)"`
	Metadata     interface{}  `gorm:"type:jsonb"`

	Source Asset `gorm:"foreignKey:SourceID"`
	Target Asset `gorm:"foreignKey:TargetID"`
}

type PoolRelationship struct {
	BaseModel
	SourcePoolID   string       `gorm:"type:varchar(64);index"`
	TargetPoolID   string       `gorm:"type:varchar(64);index"`
	RelationType   RelationType `gorm:"type:varchar(32)"`
	SourceProtocol Protocol     `gorm:"type:varchar(20)"`
	TargetProtocol Protocol     `gorm:"type:varchar(20)"`
	Metadata       interface{}  `gorm:"type:jsonb"`

	SourcePool Pool `gorm:"foreignKey:SourcePoolID"`
	TargetPool Pool `gorm:"foreignKey:TargetPoolID"`
}

type AssetPool struct {
	BaseModel
	AssetID string  `gorm:"type:varchar(64);index"`
	PoolID  string  `gorm:"type:varchar(64);index"`
	Role    string  `gorm:"type:varchar(20)"` // BASE, QUOTE, LP
	Weight  float64 // For weighted pools
	Limit   *uint64 // For capped pools

	Asset Asset `gorm:"foreignKey:AssetID"`
	Pool  Pool  `gorm:"foreignKey:PoolID"`
}

type Migration struct {
	BaseModel
	SourceAssetID string `gorm:"type:varchar(64);index"`
	TargetAssetID string `gorm:"type:varchar(64);index"`
	SourcePoolID  string `gorm:"type:varchar(64);index"`
	TargetPoolID  string `gorm:"type:varchar(64);index"`
	Type          string `gorm:"type:varchar(32)"` // PUMP_TO_RAYDIUM, VERSION_UPGRADE, etc.
	Status        string `gorm:"type:varchar(20)"`
	TxSignature   string `gorm:"type:varchar(128)"`
	BlockTime     int64
	Slot          uint64

	SourceAsset Asset `gorm:"foreignKey:SourceAssetID"`
	TargetAsset Asset `gorm:"foreignKey:TargetAssetID"`
	SourcePool  Pool  `gorm:"foreignKey:SourcePoolID"`
	TargetPool  Pool  `gorm:"foreignKey:TargetPoolID"`
}
