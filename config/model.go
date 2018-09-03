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

package config

import (
	"errors"
	"fmt"
)

const DatabaseModeBolt = "bolt"
const DatabaseModeSQLite = "sqlite"

var supportedDatabaseModes = []string{DatabaseModeBolt, DatabaseModeSQLite}

var InvalidDbModeError = errors.New(fmt.Sprintf("config: invalid database mode. Supported modes: %s",
	fmt.Sprint(supportedDatabaseModes)))

const defaultBindAddress = ":8080"
const defaultInterval = 60
const defaultSessionRefreshInterval = 3600
const defaultSessionRefreshAttempts = 5
const defaultMaxDownloadFails = 5
const defaultAutoLaunch = false
const defaultDbPath = "spectra.db"
const defaultEndpoint = "192.168.178.1"
const defaultDbMode = "bolt"

type Configuration struct {
	Credentials                RouterCredentials `yaml:"credentials" json:"credentials"`
	DatabaseMode               string            `yaml:"database_mode" json:"database_mode" env:"DB_MODE" envDefault:"bolt"`
	DatabasePath               string            `yaml:"database_path" json:"database_path" env:"DB_PATH" envDefault:"spectra.db"`
	UpdateInterval             int               `yaml:"update_interval" json:"update_interval" env:"UPDATE_INTERVAL" envDefault:"60"`
	Autolaunch                 bool              `yaml:"autolaunch" json:"autolaunch" env:"AUTOLAUNCH" envDefault:"false"`
	BindAddress                string            `yaml:"bind_address" json:"bind_address" env:"ENDPOINT_URL" envDefault:":8080"`
	SessionRefreshInterval     int               `yaml:"session_refresh_interval" json:"session_refresh_interval" env:"SESSION_REFRESH_INTERVAL" envDefault:"3600"`
	SessionRenewalAttemptCount int               `yaml:"session_refresh_attempts" json:"session_refresh_attempts" env:"SESSION_REFRESH_ATTEMPTS" envDefault:"5"`
	MaxDownloadFails           int               `yaml:"max_download_fails" json:"max_download_fails" env:"MAX_DOWNLOAD_FAILS" envDefault:"5"`
	cfgPath                    string
}

type RouterCredentials struct {
	Endpoint string `yaml:"endpoint" json:"endpoint" env:"FRITZ_ENDPOINT" envDefault:"192.168.178.1"`
	Username string `yaml:"username" json:"username" env:"FRITZ_USERNAME"`
	Password string `yaml:"password" json:"password" env:"FRITZ_PASSWORD"`
}
