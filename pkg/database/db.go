package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"corvus_bot/pkg/config"
	"corvus_bot/pkg/database/models"
)

type Database struct {
	*gorm.DB
}

// Add the ConnectDatabase function
func ConnectDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.CorvusGoDb
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Use organization's logging preferences
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return db, nil
}

// NewDatabase creates a new database connection and initializes the schema
func NewDatabase(cfg *config.Config) (*Database, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.Database.CorvusGoDb), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	database := &Database{db}

	if err := database.InitializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

// InitializeSchema sets up the database schema and required extensions
func (db *Database) InitializeSchema() error {
	// Enable PostgreSQL extensions
	extensions := []string{
		"uuid-ossp",          // For UUID generation
		"timescaledb",        // For time-series functionality
		"pg_stat_statements", // For query monitoring
	}

	for _, ext := range extensions {
		if err := db.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", ext)).Error; err != nil {
			return fmt.Errorf("failed to create extension %s: %w", ext, err)
		}
	}

	if err := db.InitializeEnums(); err != nil {
		return fmt.Errorf("failed to initialize enums: %w", err)
	}

	if err := db.InitializeTables(); err != nil {
		return fmt.Errorf("failed to initialize tables: %w", err)
	}

	return nil
}

// InitializeEnums creates custom enum types
func (db *Database) InitializeEnums() error {
	enumDefinitions := `DO $$ BEGIN
		CREATE TYPE protocol AS ENUM (
			'RAYDIUM', 'JUPITER', 'METEORA', 'MOONSHOT', 'PUMPFUN', 'ORCA'
		);
		CREATE TYPE pool_type AS ENUM (
			'AMM', 'CLMM', 'WHIRLPOOL', 'BONDING_CURVE'
		);
		CREATE TYPE asset_type AS ENUM (
			'FUNGIBLE', 'NON_FUNGIBLE', 'COMPRESSED'
		);
		CREATE TYPE asset_interface AS ENUM (
			'V1_NFT', 'FUNGIBLE_TOKEN', 'COMPRESSED_NFT'
		);
		CREATE TYPE asset_status AS ENUM (
			'ACTIVE', 'INACTIVE', 'BURNED', 'MIGRATED'
		);
		CREATE TYPE pool_status AS ENUM (
			'ACTIVE', 'INACTIVE', 'MIGRATED'
		);
		CREATE TYPE relation_type AS ENUM (
			'MIGRATION', 'WRAP', 'VERSION'
		);
		EXCEPTION WHEN duplicate_object THEN NULL;
	END $$;`

	return db.Exec(enumDefinitions).Error
}

// InitializeTables creates and migrates all tables
func (db *Database) InitializeTables() error {
	// Auto-migrate all models
	if err := db.AutoMigrate(
		&models.Asset{},
		&models.Pool{},
		&models.AssetMetric{},
		&models.PoolMetric{},
		&models.TokenMetric{},
		&models.MarketMetric{},
		&models.AssetRelationship{},
		&models.PoolRelationship{},
		&models.AssetPool{},
		&models.Migration{},
	); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	// Create hypertables for time-series data
	tables := []string{
		"asset_metrics",
		"pool_metrics",
		"token_metrics",
		"market_metrics",
	}

	for _, table := range tables {
		if err := db.Exec(
			`SELECT create_hypertable(?, 'timestamp', if_not_exists => TRUE)`,
			table,
		).Error; err != nil {
			return fmt.Errorf("failed to create hypertable for %s: %w", table, err)
		}
	}

	return nil
}

// Asset Operations
func (db *Database) UpsertAsset(asset *models.Asset) error {
	return db.Save(asset).Error
}

func (db *Database) GetAssetByAddress(address string) (*models.Asset, error) {
	var asset models.Asset
	err := db.Where("address = ?", address).First(&asset).Error
	return &asset, err
}

// Pool Operations
func (db *Database) UpsertPool(pool *models.Pool) error {
	return db.Save(pool).Error
}

func (db *Database) GetPoolByID(id string) (*models.Pool, error) {
	var pool models.Pool
	err := db.Where("id = ?", id).First(&pool).Error
	return &pool, err
}

// Metric Operations
func (db *Database) InsertMetrics(metrics interface{}) error {
	return db.Create(metrics).Error
}

func (db *Database) GetLatestMetrics(assetID string, limit int) ([]models.AssetMetric, error) {
	var metrics []models.AssetMetric
	err := db.Where("asset_id = ?", assetID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&metrics).Error
	return metrics, err
}

// Relationship Operations
func (db *Database) CreateAssetRelationship(rel *models.AssetRelationship) error {
	return db.Create(rel).Error
}

func (db *Database) GetAssetRelationships(assetID string) ([]models.AssetRelationship, error) {
	var relationships []models.AssetRelationship
	err := db.Where("source_id = ? OR target_id = ?", assetID, assetID).
		Find(&relationships).Error
	return relationships, err
}

// Migration Operations
func (db *Database) CreateMigration(migration *models.Migration) error {
	return db.Create(migration).Error
}

func (db *Database) GetMigrationsByAsset(assetID string) ([]models.Migration, error) {
	var migrations []models.Migration
	err := db.Where("source_asset_id = ? OR target_asset_id = ?", assetID, assetID).
		Order("created_at DESC").
		Find(&migrations).Error
	return migrations, err
}

func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}
	return sqlDB.Close()
}

func (db *Database) Health() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}
	return sqlDB.Ping()
}

// Transaction wrapper
func (db *Database) WithTx(fn func(*gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
