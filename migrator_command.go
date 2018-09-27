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

package main

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository/bolt"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository/migrator"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository/relational"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	migratorCommand = kingpin.Command("merge", "Merge two databases into one")

	targetPath = migratorCommand.Arg("target", "Path to target database").Required().ExistingFile()
	sourcePath = migratorCommand.Arg("source", "Path to source database").Required().ExistingFile()

	compressFlag = migratorCommand.Flag("disable-compression", "Turn of Gzip Compression of spectra data").Default("false").Bool()
)

func handleDatabaseMerge() {
	target := getRepository(*targetPath)
	source := getRepository(*sourcePath)
	if target == nil || source == nil {
		log.Error("Failed to initialize repositories")
		os.Exit(1)
	}

	defer target.Close()
	defer source.Close()

	mig, err := migrator.NewMigrator(target, *verbose, source)
	if err != nil {
		log.Errorf("Creating migrator failed. Error: %s", err.Error())
	}
	mig.Migrate()
}

func getRepository(path string) repository.Repository {
	//Attempt creating a BoltDB Repo
	var repo repository.Repository

	repo, err := bolt.NewBoltRepository(path, !*compressFlag)
	if err == nil {
		return repo
	}
	//Attempt creating a SQLite Repo
	repo, err = relational.NewSQLiteRepository(path, !*compressFlag)
	if err == nil {
		return repo
	}
	return nil
}
