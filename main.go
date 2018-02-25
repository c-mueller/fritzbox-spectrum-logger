// Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger).
// Copyright (c) 2018 Christian Müller<cmueller.dev@gmail.com>.
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
	"github.com/c-mueller/fritzbox-spectrum-logger/cmd"
	"github.com/op/go-logging"
	"os"
)

var format = logging.MustStringFormatter(
	`%{color}[%{time:15:04:05} - %{level}] - %{module}:%{color:reset} %{message}`,
)

var log = logging.MustGetLogger("main")

func main() {
	stdoutBackend := logging.NewLogBackend(os.Stdout, "", 0)

	stdoutBackendFormatter := logging.NewBackendFormatter(stdoutBackend, format)

	// Set the backends to be used.
	logging.SetBackend(stdoutBackendFormatter)

	log.Debug("Initialized Logger")

	cmd.Execute()
}
