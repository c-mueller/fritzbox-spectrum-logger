package application

import (
    "github.com/op/go-logging"
    "github.com/c-mueller/fritzbox-spectrum-logger/config"
    "github.com/gin-gonic/gin"
    "github.com/c-mueller/fritzbox-spectrum-logger/repository"
)

var log = logging.MustGetLogger("server")

func LaunchApplication(configPath string) *Application {
    log.Debugf("Loading configuration from '%s'", configPath)
    cfg, err := config.ReadOrCreate(configPath)
    if err != nil {
        log.Error("Opening of the Config failed")
        log.Error(err.Error())
        return nil
    }
    return &Application{
        config:  *cfg,
        bindAdr: cfg.BindAddress,
    }
}

func (a *Application) Listen() error {
    log.Debug("Launching server...")
    log.Debug("Initilializing repository (datastore)")
    repo, err  := repository.NewRepository(a.config.DatabasePath)

    if err != nil {
        return err
    }
    defer repo.Close()
    a.repo = repo

    log.Debug("Registering HTTP mappings")
    gin.SetMode(gin.ReleaseMode)
    engine := gin.Default()

    log.Debug("Launching HTTP Server")
    return engine.Run(a.bindAdr)
}
