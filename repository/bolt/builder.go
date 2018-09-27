package bolt

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository/reporegistry"
)

func init() {
	reporegistry.Register(sqliteBuilder{})
}

type sqliteBuilder struct {
}

func (sqliteBuilder) BuildRepository(compress bool, path string) (repository.Repository, error) {
	return NewBoltRepository(path, compress)
}

func (sqliteBuilder) GetName() string {
	return "bolt"
}
