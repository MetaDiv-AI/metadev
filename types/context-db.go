package types

import (
	"os"
	"strconv"
	"sync"

	pkglogger "github.com/MetaDiv-AI/logger"
	"github.com/MetaDiv-AI/metaorm"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// dbPoolManager manages database connections for all apps
type dbPoolManager struct {
	connections map[string]metaorm.Database
	mutex       sync.RWMutex
}

var dbPool = &dbPoolManager{
	connections: make(map[string]metaorm.Database),
}

func NewDbContext(database metaorm.Database) *dbContext {
	return &dbContext{db: database}
}

// GetConnection returns a shared database connection for the given database name
func (dpm *dbPoolManager) GetConnection(database string) metaorm.Database {
	dpm.mutex.RLock()
	if conn, exists := dpm.connections[database]; exists {
		dpm.mutex.RUnlock()
		return conn
	}
	dpm.mutex.RUnlock()

	dpm.mutex.Lock()
	defer dpm.mutex.Unlock()

	// Double-check pattern
	if conn, exists := dpm.connections[database]; exists {
		return conn
	}

	// Create new connection with proper pooling configuration
	conn := createSharedConnection(database)
	if conn == nil {
		return nil
	}
	dpm.connections[database] = conn
	return conn
}

// getDBConfig retrieves database configuration from system config
func getDBConfig() (host string, port int, username string, password string) {
	host = os.Getenv("SQL_HOST")
	port, _ = strconv.Atoi(os.Getenv("SQL_PORT"))
	username = os.Getenv("SQL_USER")
	password = os.Getenv("SQL_PASSWORD")
	return host, port, username, password
}

// createSharedConnection creates a new database connection with proper pooling
func createSharedConnection(database string) metaorm.Database {
	loggerInstance := pkglogger.New().Build()
	ginMode := os.Getenv("GIN_MODE")
	connector := metaorm.NewConnector()

	switch ginMode {
	case gin.ReleaseMode:
		host, port, username, password := getDBConfig()
		if port == 0 {
			loggerInstance.Error("Failed to create database connection: SQL configuration is missing or invalid",
				zap.String("database", database),
				zap.String("mode", "release"),
			)
			return nil
		}
		db, err := connector.MySQL().
			Host(host).
			Port(port).
			Username(username).
			Password(password).
			Database(database).
			Silent().
			Connect()
		if err != nil {
			loggerInstance.Error("Failed to connect to MySQL database",
				zap.String("database", database),
				zap.String("host", host),
				zap.Int("port", port),
				zap.String("mode", "release"),
				zap.Error(err),
			)
			return nil
		}

		// Configure connection pool settings
		configureConnectionPool(db.Gorm())
		return db

	case gin.DebugMode:
		host, port, username, password := getDBConfig()
		if port == 0 {
			loggerInstance.Error("Failed to create database connection: SQL configuration is missing or invalid",
				zap.String("database", database),
				zap.String("mode", "debug"),
			)
			return nil
		}
		db, err := connector.MySQL().
			Host(host).
			Port(port).
			Username(username).
			Password(password).
			Database(database).
			Connect()
		if err != nil {
			loggerInstance.Error("Failed to connect to MySQL database",
				zap.String("database", database),
				zap.String("host", host),
				zap.Int("port", port),
				zap.String("mode", "debug"),
				zap.Error(err),
			)
			return nil
		}

		// Configure connection pool settings
		configureConnectionPool(db.Gorm())
		return db

	default:
		db, err := connector.Sqlite().Path(database + ".db").Connect()
		if err != nil {
			loggerInstance.Error("Failed to connect to SQLite database",
				zap.String("database", database),
				zap.String("path", database+".db"),
				zap.String("mode", "default"),
				zap.Error(err),
			)
			return nil
		}

		// Configure connection pool settings
		configureConnectionPool(db.Gorm())
		return db
	}
}

// configureConnectionPool configures GORM connection pool settings
func configureConnectionPool(gormDB *gorm.DB) {
	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)   // Maximum number of open connections
	sqlDB.SetMaxIdleConns(10)   // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(0) // Connection lifetime (0 = no limit)
	sqlDB.SetConnMaxIdleTime(0) // Idle connection timeout (0 = no limit)
}

// CloseAllConnections closes all database connections
func CloseAllConnections() {
	dbPool.mutex.Lock()
	defer dbPool.mutex.Unlock()

	for _, conn := range dbPool.connections {
		if sqlDB, err := conn.Gorm().DB(); err == nil {
			sqlDB.Close()
		}
	}

	// Clear the connections map
	dbPool.connections = make(map[string]metaorm.Database)
}

type DbContext interface {
	DB() metaorm.Database
}

type dbContext struct {
	db metaorm.Database
}

func (c *dbContext) DB() metaorm.Database {
	return c.db
}
