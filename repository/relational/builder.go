package relational

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository/reporegistry"
)

func init() {
	reporegistry.Register(sqliteBuilder{})
	reporegistry.Register(mysqlBuilder{})
	reporegistry.Register(postgresBuilder{})
}

type sqliteBuilder struct {
}

func (sqliteBuilder) BuildRepository(compress bool, path string) (repository.Repository, error) {
	return NewSQLiteRepository(path, compress)
}

func (sqliteBuilder) GetName() string {
	return "sqlite"
}

type mysqlBuilder struct {
}

func (mysqlBuilder) BuildRepository(compress bool, path string) (repository.Repository, error) {
	return NewRelationalRepository("mysql", path, compress)
}

func (mysqlBuilder) GetName() string {
	return "mysql"
}

type postgresBuilder struct {
}

func (postgresBuilder) BuildRepository(compress bool, path string) (repository.Repository, error) {
	return NewRelationalRepository("postgres", path, compress)
}

func (postgresBuilder) GetName() string {
	return "postgresql"
}
