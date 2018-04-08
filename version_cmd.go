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
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

const versionDetailsFormat = "%-16s %-64s\n"

var (
	versionCmd = kingpin.Command("version", "Show 'statusbar' version information")

	version        = "master"
	revision       = "dev"
	buildNumber    = "dev"
	buildTimestamp = "dev"
)

type dependency struct {
	ProjectRoot  string `json:"ProjectRoot"`
	Constraint   string `json:"Constraint"`
	Version      string `json:"Version"`
	Revision     string `json:"Revision"`
	Latest       string `json:"Latest"`
	PackageCount int    `json:"PackageCount"`
}

type dependencies []dependency

func versionInfo() {
	var deps *dependencies
	maxNameLen := 0
	json.Unmarshal([]byte(depInfo), &deps)

	for _, v := range *deps {
		if maxNameLen < len(v.ProjectRoot) {
			maxNameLen = len(v.ProjectRoot)
		}
	}

	fmt.Printf("fritzbox-spectrum-logger - version %s.%s\n", version, revision)
	fmt.Println("Copyright (c) 2018 Christian Müller <cmueller.dev@gmail.com>")
	fmt.Println()

	fmt.Printf(versionDetailsFormat, "GITHUB", "https://github.com/c-mueller/fritzbox-spectrum-logger")
	fmt.Printf(versionDetailsFormat, "VERSION", version)
	fmt.Printf(versionDetailsFormat, "REVISION", revision)
	fmt.Printf(versionDetailsFormat, "BUILD NUMBER", buildNumber)
	fmt.Printf(versionDetailsFormat, "BUILD TIMESTAMP", buildTimestamp)
	fmt.Printf(versionDetailsFormat, "LIB COUNT", fmt.Sprintf("%d", len(*deps)))
	fmt.Println()
	fmt.Println("Dependendies:")

	format := "%-" + fmt.Sprint(maxNameLen) + "s %-20s %-40s\n"
	fmt.Printf(format, "ROOT PATH", "VERSION", "REVISION")
	for _, v := range *deps {
		version := "----"
		if len(v.Version) != 0 {
			version = strings.Replace(v.Version, "branch ", "", -1)
		}
		fmt.Printf(format, v.ProjectRoot, version, v.Revision)
	}
}
