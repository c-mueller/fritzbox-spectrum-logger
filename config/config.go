package config

import (
    "os"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "github.com/op/go-logging"
)

var log = logging.MustGetLogger("config")

func ReadOrCreate(path string) (*Configuration, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        log.Debug("Creating new Configuration")
        cfg := Configuration{
            cfgPath:        path,
            BindAddress:    defaultBindAddress,
            DatabasePath:   defaultDbPath,
            AskForPassword: defaultAskForPassword,
            Autolaunch:     defaultAutoLaunch,
            UpdateInterval: defaultInterval,
            Credentials: RouterCredentials{
                Endpoint: defaultEndpoint,
            },
        }
        log.Debug("Writing new Configuration")
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
    var data []byte
    var err error

    if !c.AskForPassword {
        data, err = yaml.Marshal(c)
    } else {
        c2 := *c
        c2.Credentials.Password = ""
        data, err = yaml.Marshal(c2)
    }
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
