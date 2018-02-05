package application

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/config"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"time"
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
		config:            *cfg,
		bindAdr:           cfg.BindAddress,
		state:             IDLE,
		startTime:         time.Now(),
		sessionLogCounter: 0,
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
	engine.POST("/api/control/start", a.startCollecting)
	engine.POST("/api/control/stop", a.stopCollecting)
}
