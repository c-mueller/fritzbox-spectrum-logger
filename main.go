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
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose",
		"Run command with verbose output").Short('v').Default("false").Bool()
	debug = kingpin.Flag("debug",
		"Run command in debug mode (includes verbose mode)").Short('d').Default("false").Bool()
)

func main() {
	initializeLogger()

	switch kingpin.Parse() {
	case "server":
		launchServerWithConfig()
	case "server env":
		launchServerFromEnvironment()
	case "version":
		versionInfo()
	case "generate-config yaml":
		generateYaml()
	case "generate-config docker":
		generateDockerfileCommands()
	case "generate-config bash":
		generateBash()
	case "generate-config fish":
		generateFish()
	}
}
