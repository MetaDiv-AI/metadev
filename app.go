package metadev

import (
	"fmt"
	"os"

	"github.com/MetaDiv-AI/metadev/types"
)

// NewApp creates a new app builder
// name is the name of the app
func NewApp(name string) *appBuilder {
	return &appBuilder{name: name}
}

type appBuilder struct {
	name       string
	migrations []types.Model
}

// Database sets the database of the app (optional)
func (b *appBuilder) SQLMigrates(migrations ...types.Model) *appBuilder {
	if len(migrations) > 0 {
		b.migrations = migrations
	}
	return b
}

func (b *appBuilder) Build() types.App {
	_, ok := Apps[b.name]
	if ok {
		panic(fmt.Sprintf("app %s already registered", b.name))
	}
	Apps[b.name] = types.NewApp(b.name, os.Getenv("SQL_DATABASE"), os.Getenv("MONGO_DATABASE"), b.migrations)
	return Apps[b.name]
}
