package config

import (
    "os"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

func ReadOrCreate(path string) (*Configuration, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        cfg := Configuration{
            cfgPath: path,
        }
        err = cfg.Write()
        if err != nil {
            return nil, err
        }
        return &cfg, nil
    }
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
