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
