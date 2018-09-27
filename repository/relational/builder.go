// Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger).
// Copyright (c) 2018 Christian Müller <cmueller.dev@gmail.com>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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
