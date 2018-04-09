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
	"github.com/c-mueller/fritzbox-spectrum-logger/server"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	serverCmd    = kingpin.Command("server", "Run server")
	serverConfig = serverCmd.Flag("config",
		"Path to the config file").Default("config.yml").ExistingFile()

	serverEnvCmd = serverCmd.Command("env", "Launch from environment configuration")
)

func launchServerWithConfig() {
	srv := server.LaunchApplication(*serverConfig)
	err := srv.Listen()
	if err != nil {
		log.Error(err.Error())
	}
}

func launchServerFromEnvironment() {
	srv := server.LaunchFromEnvironment()
	err := srv.Listen()
	if err != nil {
		log.Error(err.Error())
	}
}
