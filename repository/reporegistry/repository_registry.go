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

package reporegistry

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/op/go-logging"
)

var repos = make(map[string]repository.RepoBuilder, 0)

var log = logging.MustGetLogger("repository_registry")

func Register(builder repository.RepoBuilder) {
	log.Infof("Registered %q repository", builder.GetName())
	repos[builder.GetName()] = builder
}

func GetRepositories() []repository.RepoBuilder {
	builders := make([]repository.RepoBuilder, 0)
	for _, v := range repos {
		builders = append(builders, v)
	}
	return builders
}

func GetForName(name string) repository.RepoBuilder {
	return repos[name]
}
