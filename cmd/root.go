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

package cmd

import (
	"fmt"
	"os"

	"github.com/c-mueller/fritzbox-spectrum-logger/application"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var log = logging.MustGetLogger("cli")

var cfgFile string

var appCfgFlag string

var usernameFlag string
var passwordFlag string
var endpointFlag string
var dbPathFlag string

var RootCmd = &cobra.Command{
	Use:   "fritzbox-spectrum-logger",
	Short: "This application requests the DSL Spectrum of a Fritz!Box in a given interval and Stores it for later review.",
	Run:   launchServer,
}

func launchServer(cmd *cobra.Command, args []string) {
	app := application.LaunchApplication(appCfgFlag)
	if app != nil {
		app.Listen()
	}
}

func mainCommand(cmd *cobra.Command, args []string) {
	log.Info("Logging in...")
	client := fritz.NewClient(endpointFlag, usernameFlag, passwordFlag)
	err := client.Login()
	failOnError(err)
	log.Info("Login Done!")

	log.Info("Opening DB...")
	repo, err := repository.NewRepository(dbPathFlag)
	failOnError(err)
	defer repo.Close()
	log.Info("Done!")

	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		go func() {
			currentTime := time.Now()
			log.Infof("[%d:%d:%d]: Downloading Spectrum...",
				currentTime.Hour(), currentTime.Minute(), currentTime.Second())
			spectrum, err := client.GetSpectrum()
			if err != nil {
				log.Error("Fail!")
				log.Error(err.Error())
				return
			}
			err = repo.Insert(spectrum)
			if err != nil {
				log.Error("Fail!")
				log.Error(err.Error())
				return
			}
			log.Info("Download Done!")
		}()
	}
}

func failOnError(err error) {
	if err != nil {
		fmt.Println("Failed!")
		fmt.Println(err)
		os.Exit(1)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.fritzbox-spectrum-logger)")

	RootCmd.Flags().StringVarP(&endpointFlag, "endpoint", "e",
		"192.168.178.1", "The Endpoint of the Fritz!Box (IP or Hostname)")
	RootCmd.Flags().StringVarP(&usernameFlag, "username", "u",
		"", "The Username to login, empty if none")
	RootCmd.Flags().StringVarP(&passwordFlag, "password", "p",
		"", "The password used to login, empty if none")
	RootCmd.Flags().StringVarP(&dbPathFlag, "db-path", "d",
		"spectra.db", "The path to store the spectrums at")

	RootCmd.Flags().StringVarP(&appCfgFlag, "app-config", "c", "config.yml", "The path to the config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".fritzbox-spectrum-logger")
	}

	viper.AutomaticEnv()

	viper.SetDefault("endpoint", "fritz.box")
	viper.SetDefault("use-tls", "false")
	viper.SetDefault("username", "")
	viper.SetDefault("password", "password")
	viper.SetDefault("application-endpoint", ":8080")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
