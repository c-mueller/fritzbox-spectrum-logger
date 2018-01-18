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

package cmd

import (
    "fmt"
    "os"

    "github.com/mitchellh/go-homedir"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/c-mueller/fritzbox-spectrum-logger/fritz"
    "time"
    "github.com/c-mueller/fritzbox-spectrum-logger/repository"
)

var cfgFile string

var usernameFlag string
var passwordFlag string
var endpointFlag string
var dbPathFlag string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
    Use:   "fritzbox-spectrum-logger",
    Short: "This application requests the DSL Spectrum of a Fritz!Box in a given interval and Stores it for later review.",
    Run:   mainCommand,
}

func mainCommand(cmd *cobra.Command, args []string) {
    fmt.Print("Logging in...")
    client := fritz.NewClient(endpointFlag, usernameFlag, passwordFlag)
    err := client.Login()
    failOnError(err)
    fmt.Println("Done!")

    fmt.Print("Opening DB...")
    repo, err := repository.NewRepository(dbPathFlag)
    failOnError(err)
    defer repo.Close()
    fmt.Println("Done!")

    ticker := time.NewTicker(time.Second * 5)
    for range ticker.C {
        go func() {
            currentTime := time.Now()
            fmt.Printf("[%d:%d:%d]: Downloading Spectrum...",
                currentTime.Hour(), currentTime.Minute(), currentTime.Second())
            spectrum, err := client.GetSpectrum()
            if err != nil {
                fmt.Println("Fail!")
                fmt.Println(err)
                return
            }
            err = repo.Insert(spectrum)
            if err != nil {
                fmt.Println("Fail!")
                fmt.Println(err)
                return
            }
            fmt.Println("Done!")
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Find home directory.
        home, err := homedir.Dir()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // Search config in home directory with name ".fritzbox-spectrum-logger" (without extension).
        viper.AddConfigPath(home)
        viper.SetConfigName(".fritzbox-spectrum-logger")
    }

    viper.AutomaticEnv() // read in environment variables that match

    viper.SetDefault("endpoint", "fritz.box")
    viper.SetDefault("use-tls", "false")
    viper.SetDefault("username", "")
    viper.SetDefault("password", "password")
    viper.SetDefault("application-endpoint", ":8080")

    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    }
}
