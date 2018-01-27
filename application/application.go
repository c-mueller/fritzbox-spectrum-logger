package application

import (
    "github.com/op/go-logging"
    "github.com/c-mueller/fritzbox-spectrum-logger/config"
    "github.com/gin-gonic/gin"
    "github.com/c-mueller/fritzbox-spectrum-logger/repository"
    "time"
    "io/ioutil"
    "encoding/json"
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
        config:    *cfg,
        bindAdr:   cfg.BindAddress,
        state:     IDLE,
        startTime: time.Now(),
    }
}

func (a *Application) Listen() error {
    log.Debug("Launching server...")
    log.Debug("Initilializing repository (datastore)")
    repo, err := repository.NewRepository(a.config.DatabasePath)

    if err != nil {
        return err
    }
    defer repo.Close()
    a.repo = repo

    log.Debug("Registering HTTP mappings")
    gin.SetMode(gin.ReleaseMode)
    engine := gin.Default()
    a.registerHTTPMappings(engine)

    log.Debug("Launching HTTP Server")
    return engine.Run(a.bindAdr)
}

func (a *Application) registerHTTPMappings(engine *gin.Engine) {
    engine.GET("/api/status", a.getStatus)
    engine.GET("/api/config", a.getConfiguration)

    engine.POST("/api/config", a.updateConfig)
}

func (a *Application) getStatus(ctx *gin.Context) {
    ctx.JSON(200, StatusResponse{
        State:  a.state.String(),
        Uptime: int64(time.Since(a.startTime).Seconds()),
    })
}

func (a *Application) getConfiguration(ctx *gin.Context) {
    var cfg config.Configuration = a.config
    cfg.Credentials.Password = "HIDDEN"
    ctx.JSON(200, cfg)
}

func (a *Application) updateConfig(ctx *gin.Context) {
    data, err := ioutil.ReadAll(ctx.Request.Body)
    if err != nil {
        ctx.AbortWithError(400, InvalidBodyError)
        return
    }
    var parsedBody *config.Configuration
    err = json.Unmarshal(data, &parsedBody)
    if err != nil {
        ctx.AbortWithError(400, JSONParsingError)
        return
    }
    a.config.Update(parsedBody)
    err = a.config.Write()
    if err != nil {
        ctx.AbortWithError(500, FileSystemError)
    }
    ctx.String(200, "")
}
