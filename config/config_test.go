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
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Handling(t *testing.T) {
	t.Log("Getting Tempdir")
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	t.Log("Reading empty file")
	configPath := filepath.Join(tmpdir, "config.yml")
	cfg, _ := ReadOrCreate(configPath)

	t.Log("Updating Data")
	cfg.DatabasePath = "TEST"
	cfg.UpdateInterval = 1337
	cfg.Credentials = RouterCredentials{
		Endpoint: "1.1.1.1",
		Username: "Test",
		Password: "123",
	}

	t.Log("Writing data")
	cfg.Write()

	t.Log("Reading Config again")
	cfg2, _ := ReadOrCreate(configPath)

	t.Log("Checking if both instances are equal")
	if cfg.DatabasePath != cfg2.DatabasePath {
		t.FailNow()
	} else if cfg.UpdateInterval != cfg2.UpdateInterval {
		t.FailNow()
	} else if cfg.Credentials != cfg2.Credentials {
		t.FailNow()
	}
}

func Test_LoadFromEnv(t *testing.T) {
	const dbPathTestValue = "myspectra.db"
	const fritzCredsUsernameTestValue = "caaarl"
	const fritzCredsPasswordTestValue = "passw0rd"
	const endpointDefault = "192.168.178.1"

	os.Setenv("DB_PATH", dbPathTestValue)
	os.Setenv("AUTOLAUNCH", "true")
	os.Setenv("FRITZ_USERNAME", fritzCredsUsernameTestValue)
	os.Setenv("FRITZ_PASSWORD", fritzCredsPasswordTestValue)

	cfg, err := FromEnvironment()
	assert.NoError(t, err)

	assert.True(t, cfg.Autolaunch)
	assert.Equal(t, cfg.DatabasePath, dbPathTestValue)
	assert.Equal(t, cfg.Credentials.Username, fritzCredsUsernameTestValue)
	assert.Equal(t, cfg.Credentials.Password, fritzCredsPasswordTestValue)
	assert.Equal(t, cfg.Credentials.Endpoint, endpointDefault)
}
