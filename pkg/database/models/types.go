package models

type Protocol string
type PoolType string
type AssetType string
type AssetInterface string
type AssetStatus string
type PoolStatus string
type RelationType string

const (
	// Protocols
	ProtocolRaydium  Protocol = "RAYDIUM"
	ProtocolJupiter  Protocol = "JUPITER"
	ProtocolMeteora  Protocol = "METEORA"
	ProtocolMoonshot Protocol = "MOONSHOT"
	ProtocolPumpFun  Protocol = "PUMPFUN"
	ProtocolOrca     Protocol = "ORCA"

	// Pool Types
	PoolTypeAMM          PoolType = "AMM"
	PoolTypeCLMM         PoolType = "CLMM"
	PoolTypeWhirlpool    PoolType = "WHIRLPOOL"
	PoolTypeBondingCurve PoolType = "BONDING_CURVE"

	// Asset Types
	AssetTypeFungible    AssetType = "FUNGIBLE"
	AssetTypeNonFungible AssetType = "NON_FUNGIBLE"
	AssetTypeCompressed  AssetType = "COMPRESSED"

	// Asset Interfaces
	InterfaceV1NFT         AssetInterface = "V1_NFT"
	InterfaceFungibleToken AssetInterface = "FUNGIBLE_TOKEN"
	InterfaceCompressedNFT AssetInterface = "COMPRESSED_NFT"

	// Asset Statuses
	AssetStatusActive   AssetStatus = "ACTIVE"
	AssetStatusInactive AssetStatus = "INACTIVE"
	AssetStatusBurned   AssetStatus = "BURNED"
	AssetStatusMigrated AssetStatus = "MIGRATED"

	// Pool Statuses
	PoolStatusActive   PoolStatus = "ACTIVE"
	PoolStatusInactive PoolStatus = "INACTIVE"
	PoolStatusMigrated PoolStatus = "MIGRATED"

	// Relationship Types
	RelationTypeMigration RelationType = "MIGRATION"
	RelationTypeWrap      RelationType = "WRAP"
	RelationTypeVersion   RelationType = "VERSION"
)
