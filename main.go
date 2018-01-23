// Copyright © 2018 Christian Müller <cmueller.dev@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	stdoutBackend := logging.NewLogBackend(os.Stdout, "", 0)

	stdoutBackendFormatter := logging.NewBackendFormatter(stdoutBackend, format)
	stderrLeveled := logging.AddModuleLevel(stderrBackend)
	stderrLeveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(stderrLeveled, stdoutBackendFormatter)

	log.Debug("Initialized Logger")

	cmd.Execute()
}
