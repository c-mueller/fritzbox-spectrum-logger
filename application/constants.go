package application

import (
	"errors"
)

type ApplicationState int

const IDLE ApplicationState = iota
const LOGGING ApplicationState = IDLE + 1
const ERROR ApplicationState = LOGGING + 1

var InvalidBodyError = errors.New("application: Could not deserialize request body. The body has to be JSON")
var JSONParsingError = errors.New("application: Could not parse JSON")
var FileSystemError = errors.New("application: Fileaccess has failed")

func (s ApplicationState) String() string {
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
