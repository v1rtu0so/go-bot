// pkg/database/db.go

/* EXAMPLE USAGE IN APP:
func main() {
    cfg, err := config.LoadConfig("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    db, err := database.NewDatabase(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    // Check database health
    if err := db.Health(); err != nil {
        log.Fatalf("Database health check failed: %v", err)
    }

    // Your application logic here...
}
*/

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

// Database wraps the GORM DB connection and provides additional functionality
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection and initializes the schema
func NewDatabase(cfg *config.Config) (*Database, error) {
	// Create custom logger for GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(cfg.Database.CorvusGoDb), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	database := &Database{db}

	// Initialize database schema
	if err := database.InitializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

// InitializeSchema sets up the database schema and required extensions
func (db *Database) InitializeSchema() error {
	// Enable required PostgreSQL extensions
	extensions := []string{
		"uuid-ossp",          // For UUID generation
		"timescaledb",        // For time-series functionality
		"pg_stat_statements", // For query performance monitoring
	}

	for _, ext := range extensions {
		if err := db.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", ext)).Error; err != nil {
			return fmt.Errorf("failed to create extension %s: %w", ext, err)
		}
	}

	// Initialize tables and types
	if err := db.InitializeTables(); err != nil {
		return fmt.Errorf("failed to initialize tables: %w", err)
	}

	return nil
}

// InitializeTables creates all necessary tables and types
func (db *Database) InitializeTables() error {
	// Create protocol and pool type enums
	if err := db.Exec(`DO $$ BEGIN
        CREATE TYPE protocol AS ENUM (
            'RAYDIUM', 
            'JUPITER', 
            'METEORA', 
            'MOONSHOT', 
            'PUMPFUN',
            'ORCA'
        );
        CREATE TYPE pool_type AS ENUM (
            'AMM', 
            'CLMM', 
            'BONDING_CURVE',
            'WHIRLPOOL',
            'AGGREGATOR'
        );
        EXCEPTION WHEN duplicate_object THEN NULL;
    END $$;`).Error; err != nil {
		return fmt.Errorf("failed to create enums: %w", err)
	}

	// Auto-migrate all pool types
	if err := db.AutoMigrate(
		&models.RaydiumAMMPool{},
		&models.RaydiumCLMMPool{},
		&models.PumpFunPool{},
		&models.JupiterPool{},
		&models.MeteoraPool{},
		&models.MoonshotPool{},
		&models.OrcaWhirlpool{},
		&models.RaydiumAMMMetric{},
		&models.RaydiumCLMMMetric{},
		&models.PumpFunMetric{},
		&models.JupiterMetric{},
		&models.MeteoraMetric{},
		&models.MoonshotMetric{},
		&models.OrcaMetric{},
	); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	// Create hypertables for all metric tables
	tables := []string{
		"raydium_amm_metrics",
		"raydium_clmm_metrics",
		"pump_fun_metrics",
		"jupiter_metrics",
		"meteora_metrics",
		"moonshot_metrics",
		"orca_metrics",
	}

	// Create hypertables and indexes
	for _, table := range tables {
		// Create hypertable
		if err := db.Exec(
			`SELECT create_hypertable(?, 'timestamp', if_not_exists => TRUE)`,
			table,
		).Error; err != nil {
			return fmt.Errorf("failed to create hypertable for %s: %w", table, err)
		}
	}

	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Health checks if the database connection is healthy
func (db *Database) Health() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}
	return sqlDB.Ping()
}

// Truncate truncates all tables (useful for testing)
func (db *Database) Truncate() error {
	tables := []string{
		"raydium_amm_pools",
		"raydium_clmm_pools",
		"pump_fun_pools",
		"raydium_amm_metrics",
		"raydium_clmm_metrics",
		"pump_fun_metrics",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}
