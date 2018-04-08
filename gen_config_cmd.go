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
	"fmt"
	"github.com/c-mueller/fritzbox-spectrum-logger/config"
	"github.com/c-mueller/fritzbox-spectrum-logger/util"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	genConfigCmd = kingpin.Command("generate-config", "Create a configuration file")

	// Generator modes
	yamlGenCmd   = genConfigCmd.Command("yaml", "Print a yaml configuration to stdout")
	dockerGenCmd = genConfigCmd.Command("docker", "Print dockerfile commands to setup the environment to stdout")
	bashGenCmd   = genConfigCmd.Command("bash", "Print environment setup commands for bash to stdout")
	fishGenCmd   = genConfigCmd.Command("fish", "Print environment setup commands for fish to stdout")

	// Login Flags
	genEndpointFlag = genConfigCmd.Flag("endpoint",
		"The endpoint url of the Fritz!Box to monitor").Default("192.168.178.1").Short('e').URL()
	genUsernameFlag = genConfigCmd.Flag("username",
		"The username used to identify on the Fritz!Box").Default("").Short('u').String()
	genPasswordFlag = genConfigCmd.Flag("password",
		"The password used to authenticate with the Fritz!Box").Default("").Short('p').String()

	// Repository Flags
	genRepoPathFlag = genConfigCmd.Flag("db-path", "Path to the database").Default("spectra.db").File()

	// Server Flags
	genServerPortFlag = genConfigCmd.Flag("bind-url",
		"The endpoint url for the server to listen on (example: ':8080' will listen on 0.0.0.0:8080)").Short('P').Default(":8080").String()
	genUpdateInterval = genConfigCmd.Flag("download-interval",
		"The interval (in seconds) in which the spectrum should be downloaded").Short('U').Default("60").Int()
	genSessionRefreshAttempts = genConfigCmd.Flag("session-renewal-attempts",
		"How often should fsl try to renew the session if a attempt failed?").Default("5").Int()
	genSessionRefreshInterval = genConfigCmd.Flag("session-renewal-interval",
		"After how many seconds should the current session get renewed?").Default("3600").Int()
	genMaxDownloadFails = genConfigCmd.Flag("max-download-fails",
		"How many spectrum download failures are allowed before logging stops?").Default("5").Int()
	genAutolaunch = genConfigCmd.Flag("autolaunch",
		"Automatically start logging once the server is ready").Default("false").Short('a').Bool()
)

func generateYaml() {
	printOutputString(func(c *config.Configuration) string {
		data, err := yaml.Marshal(c)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		return string(data)
	})
}

func generateDockerfileCommands() {
	printOutputString(func(c *config.Configuration) string {
		a, _ := util.MapStructToDockerfileCommands(*c)
		return a
	})
}

func generateBash() {
	printOutputString(func(c *config.Configuration) string {
		a, _ := util.MapStructToBashCommands(*c)
		return a
	})
}

func generateFish() {
	printOutputString(func(c *config.Configuration) string {
		a, _ := util.MapStructToFishCommands(*c)
		return a
	})
}

func getConfigFromFlags() *config.Configuration {
	return &config.Configuration{
		Credentials: config.RouterCredentials{
			Endpoint: (*genEndpointFlag).String(),
			Username: *genUsernameFlag,
			Password: *genPasswordFlag,
		},
		DatabasePath:               (*genRepoPathFlag).Name(),
		UpdateInterval:             *genUpdateInterval,
		Autolaunch:                 *genAutolaunch,
		BindAddress:                *genServerPortFlag,
		SessionRenewalAttemptCount: *genSessionRefreshAttempts,
		SessionRefreshInterval:     *genSessionRefreshInterval,
		MaxDownloadFails:           *genMaxDownloadFails,
	}
}

func printOutputString(f func(c *config.Configuration) string) {
	cfg := getConfigFromFlags()
	fmt.Print(f(cfg))
}
