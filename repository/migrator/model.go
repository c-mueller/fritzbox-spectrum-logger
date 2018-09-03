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

package migrator

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("migrator")

type DatabaseMigrator struct {
	TargetRepository   repository.Repository
	SourceRepositories []repository.Repository
	Verbose            bool
}