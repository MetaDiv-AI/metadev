package types

import (
	"context"
	"os"
	"sync"

	"github.com/MetaDiv-AI/metamongo"
)

// mongoPoolManager manages MongoDB connections for all apps
type mongoPoolManager struct {
	connections map[string]metamongo.Database
	mutex       sync.RWMutex
}

var mongoPool = &mongoPoolManager{
	connections: make(map[string]metamongo.Database),
}

func NewMongoContext(database metamongo.Database) *mongoContext {
	return &mongoContext{db: database}
}

// GetConnection returns a shared MongoDB connection for the given database name
func (mpm *mongoPoolManager) GetConnection(database string) metamongo.Database {
	mpm.mutex.RLock()
	if conn, exists := mpm.connections[database]; exists {
		mpm.mutex.RUnlock()
		return conn
	}
	mpm.mutex.RUnlock()

	mpm.mutex.Lock()
	defer mpm.mutex.Unlock()

	// Double-check pattern
	if conn, exists := mpm.connections[database]; exists {
		return conn
	}

	// Create new connection
	conn := createSharedMongoConnection(database)
	if conn == nil {
		return nil
	}
	mpm.connections[database] = conn
	return conn
}

// mongoConfig represents the MongoDB configuration structure
type mongoConfig struct {
	MongoURI      string `json:"mongo_uri"`
	MongoUsername string `json:"mongo_username"`
	MongoPassword string `json:"mongo_password"`
}

// getMongoConfig retrieves MongoDB configuration from system config
func getMongoConfig() (uri string, username string, password string) {
	uri = os.Getenv("MONGO_URI")
	username = os.Getenv("MONGO_USERNAME")
	password = os.Getenv("MONGO_PASSWORD")
	return uri, username, password
}

// createSharedMongoConnection creates a new MongoDB connection
func createSharedMongoConnection(database string) metamongo.Database {
	uri, username, password := getMongoConfig()
	if uri == "" {
		return nil
	}

	connector := metamongo.NewConnector()
	client, err := connector.
		URI(uri).
		Username(username).
		Password(password).
		Database("admin"). // default database
		Connect()
	if err != nil {
		return nil
	}
	return metamongo.NewDatabase(client, database)
}

// CloseAllMongoConnections closes all MongoDB connections
func CloseAllMongoConnections() {
	mongoPool.mutex.Lock()
	defer mongoPool.mutex.Unlock()

	for _, conn := range mongoPool.connections {
		if client := conn.Client(); client != nil {
			client.Disconnect(context.TODO())
		}
	}

	// Clear the connections map
	mongoPool.connections = make(map[string]metamongo.Database)
}

type MongoContext interface {
	MongoDB() metamongo.Database
}

type mongoContext struct {
	db metamongo.Database
}

func (c *mongoContext) MongoDB() metamongo.Database {
	return c.db
}
