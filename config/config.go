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
	"github.com/caarlos0/env"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var log = logging.MustGetLogger("config")

var DefaultConfig = Configuration{
	BindAddress:                defaultBindAddress,
	DatabasePath:               defaultDbPath,
	Autolaunch:                 defaultAutoLaunch,
	SessionRenewalAttemptCount: defaultSessionRefreshAttempts,
	MaxDownloadFails:           defaultMaxDownloadFails,
	UpdateInterval:             defaultInterval,
	SessionRefreshInterval:     defaultSessionRefreshInterval,
	Credentials: RouterCredentials{
		Endpoint: defaultEndpoint,
	},
}

func FromEnvironment() (*Configuration, error) {
	data := &Configuration{}
	creds := RouterCredentials{}

	err := env.Parse(&creds)
	if err != nil {
		return nil, err
	}

	err = env.Parse(data)
	if err != nil {
		return nil, err
	}

	data.Credentials = creds

	return data, nil
}

func ReadOrCreate(path string) (*Configuration, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Debug("Creating new Configuration")

		log.Debug("Writing new Configuration")

		var cfg Configuration
		cfg = DefaultConfig
		cfg.cfgPath = path

		err = cfg.Write()
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	log.Debug("Loading Configuration")
	cfgFile, err := os.Open(path)
	defer cfgFile.Close()
	if err != nil {
		return nil, err
	}

	cfgData, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return nil, err
	}

	var conf Configuration
	err = yaml.Unmarshal(cfgData, &conf)
	if err != nil {
		return nil, err
	}
	conf.cfgPath = path

	return &conf, nil
}

func (c *Configuration) Update(cfg *Configuration) {
	//Only allow the updating of the Credentials, Intervall and autolaunch property using the webinterface
	c.Credentials = cfg.Credentials
	c.Autolaunch = cfg.Autolaunch
	c.UpdateInterval = cfg.UpdateInterval
}

func (c *Configuration) Write() error {
	// Close write if the object comes from the environment
	if c.cfgPath == "" {
		return nil
	}

	var data []byte
	var err error
	data, err = yaml.Marshal(c)
	if err != nil {
		return err
	}

	file, err := os.Create(c.cfgPath)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
