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
	"os"
	"github.com/op/go-logging"
)

var format = logging.MustStringFormatter(
	`%{color}[%{time:15:04:05} - %{level}] - %{module}:%{color:reset} %{message}`,
)

var log = logging.MustGetLogger("main")

func initializeLogger() {
	stdoutBackend := logging.NewLogBackend(os.Stdout, "", 0)

	backendFormatter := logging.NewBackendFormatter(stdoutBackend, format)

	leveledBackend := logging.AddModuleLevel(backendFormatter)

	if *debug {
		leveledBackend.SetLevel(logging.DEBUG, "")
	} else if *verbose {
		leveledBackend.SetLevel(logging.INFO, "")
	} else {
		leveledBackend.SetLevel(logging.ERROR, "")
	}

	logging.SetBackend(leveledBackend)
	log.Debug("Initialized Logger")
}
