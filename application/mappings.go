package application

import (
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "time"
    "github.com/c-mueller/fritzbox-spectrum-logger/config"
    "encoding/json"
)

func (a *Application) getStatus(ctx *gin.Context) {
    ctx.JSON(200, StatusResponse{
        State:         a.state.String(),
        Uptime:        int64(time.Since(a.startTime).Seconds()),
        SpectrumCount: a.sessionLogCounter,
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
