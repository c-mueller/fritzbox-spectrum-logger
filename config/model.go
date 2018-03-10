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

package config

const defaultBindAddress = ":8080"
const defaultInterval = 60
const defaultSessionRefreshInterval = 3600
const defaultAskForPassword = false
const defaultAutoLaunch = false
const defaultDbPath = "spectra.db"
const defaultEndpoint = "192.168.178.1"

type Configuration struct {
	Credentials            RouterCredentials `yaml:"credentials" json:"credentials"`
	DatabasePath           string            `yaml:"database_path" json:"database_path"`
	UpdateInterval         int               `yaml:"update_interval" json:"update_interval"`
	AskForPassword         bool              `yaml:"ask_for_password" json:"ask_for_password"`
	Autolaunch             bool              `yaml:"autolaunch" json:"autolaunch"`
	BindAddress            string            `yaml:"bind_address" json:"bind_address"`
	SessionRefreshInterval int               `yaml:"session_refresh_interval" json:"session_refresh_interval" `
	cfgPath                string
}

type RouterCredentials struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}
