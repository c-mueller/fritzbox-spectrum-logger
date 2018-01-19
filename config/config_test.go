package config

import (
    "testing"
    "github.com/Flaque/filet"
    "path/filepath"
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
