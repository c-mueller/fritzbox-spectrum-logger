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

package application

import (
	"errors"
)

type APIState int

const IDLE APIState = iota
const LOGGING APIState = IDLE + 1
const ERROR APIState = LOGGING + 1

var InvalidBodyError = errors.New("application: Could not deserialize request body. The body has to be JSON")
var JSONParsingError = errors.New("application: Could not parse JSON")
var FileSystemError = errors.New("application: Fileaccess has failed")

func (s APIState) String() string {
	if s == IDLE {
		return "IDLE"
	} else if s == LOGGING {
		return "LOGGING"
	} else if s == ERROR {
		return "ERROR"
	} else {
		return "ILLEGAL"
	}
}
