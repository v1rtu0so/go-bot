package models

type Asset struct {
	BaseAsset
	Content     Content      `gorm:"type:jsonb"`
	Authorities []Authority  `gorm:"type:jsonb"`
	Compression *Compression `gorm:"type:jsonb"`
	Grouping    []Grouping   `gorm:"type:jsonb"`
	Royalty     *Royalty     `gorm:"type:jsonb"`
	Creators    []Creator    `gorm:"type:jsonb"`
	Ownership   Ownership    `gorm:"type:jsonb"`
	TokenInfo   TokenInfo    `gorm:"type:jsonb"`
	MarketData  MarketData   `gorm:"type:jsonb"`

	// Token specific data
	TokenData *TokenData `gorm:"type:jsonb"`
	NFTData   *NFTData   `gorm:"type:jsonb"`

	// Relationships
	Pools   []Pool        `gorm:"many2many:asset_pools;"`
	Metrics []AssetMetric `gorm:"foreignKey:AssetID"`
}

type Content struct {
	Schema   string   `json:"$schema"`
	JsonUri  string   `json:"json_uri"`
	Files    []File   `json:"files"`
	Metadata Metadata `json:"metadata"`
}

type File struct {
	Uri         string  `json:"uri"`
	FileType    string  `json:"file_type"`
	ContentType string  `json:"content_type"`
	CDN         *string `json:"cdn,omitempty"`
}

type Metadata struct {
	Name        string      `json:"name"`
	Symbol      string      `json:"symbol"`
	Description string      `json:"description"`
	Attributes  []Attribute `json:"attributes"`
	External    *string     `json:"external_url,omitempty"`
}

type Authority struct {
	Address string   `json:"address"`
	Scopes  []string `json:"scopes"`
}

type Compression struct {
	Eligible    bool   `json:"eligible"`
	Compressed  bool   `json:"compressed"`
	DataHash    string `json:"data_hash"`
	CreatorHash string `json:"creator_hash"`
	AssetHash   string `json:"asset_hash"`
	Tree        string `json:"tree"`
	Seq         uint64 `json:"seq,omitempty"`
	LeafId      uint64 `json:"leaf_id,omitempty"`
}

type Grouping struct {
	GroupKey   string `json:"group_key"`
	GroupValue string `json:"group_value"`
}

type Royalty struct {
	RoyaltyModel        string  `json:"royalty_model"`
	Target              string  `json:"target"`
	Percent             float64 `json:"percent"`
	PrimarySaleHappened bool    `json:"primary_sale_happened"`
	Locked              bool    `json:"locked"`
}

type Creator struct {
	Address  string `json:"address"`
	Verified bool   `json:"verified"`
	Share    uint8  `json:"share"`
}

type Ownership struct {
	Frozen    bool   `json:"frozen"`
	Delegated bool   `json:"delegated"`
	Delegate  string `json:"delegate"`
	Owner     string `json:"owner"`
}

type TokenInfo struct {
	TokenProgram    string `json:"token_program"`
	MintAuthority   string `json:"mint_authority"`
	FreezeAuthority string `json:"freeze_authority"`
	Decimals        uint8  `json:"decimals"`
}

type TokenData struct {
	CurrentSupply     uint64  `json:"current_supply"`
	CirculatingSupply uint64  `json:"circulating_supply"`
	MaxSupply         *uint64 `json:"max_supply,omitempty"`
	Price             float64 `json:"price"`
	MarketCap         float64 `json:"market_cap"`
	Volume24h         float64 `json:"volume_24h"`
	PriceChange24h    float64 `json:"price_change_24h"`
	HolderCount       uint32  `json:"holder_count"`
}

type NFTData struct {
	Collection       string   `json:"collection"`
	CollectionFamily string   `json:"collection_family"`
	RarityRank       *uint32  `json:"rarity_rank,omitempty"`
	RarityScore      *float64 `json:"rarity_score,omitempty"`
	EditionNumber    *uint32  `json:"edition_number,omitempty"`
	MaxEditions      *uint32  `json:"max_editions,omitempty"`
}

type MarketData struct {
	Price          float64 `json:"price"`
	MarketCap      float64 `json:"market_cap"`
	Volume24h      float64 `json:"volume_24h"`
	PriceChange24h float64 `json:"price_change_24h"`
	HolderCount    uint32  `json:"holder_count"`
	ListingCount   uint32  `json:"listing_count"`
	Liquidity      float64 `json:"liquidity"`
}

type Attribute struct {
	TraitType string   `json:"trait_type"`
	Value     string   `json:"value"`
	MaxValue  *string  `json:"max_value,omitempty"`
	MinValue  *string  `json:"min_value,omitempty"`
	Count     *uint32  `json:"count,omitempty"`
	Frequency *float64 `json:"frequency,omitempty"`
}
