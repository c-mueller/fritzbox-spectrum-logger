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

import (
	"github.com/Flaque/filet"
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
