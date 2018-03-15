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

package application

import (
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/c-mueller/fritzbox-spectrum-logger/config"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
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
	log.Debug("Initializing repository (Datastore)")
	repo, err := repository.NewRepository(a.config.DatabasePath)

	if err != nil {
		return err
	}
	defer repo.Close()
	a.repo = repo

	log.Debug("Registering HTTP mappings")
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.Use(cors.Default())

	a.registerHTTPMappings(engine)

	if a.config.Autolaunch {
		log.Debug("Autolaunching Spectrum Logging")
		a.startLogging()
	}

	log.Debug("Launched HTTP Server")
	return engine.Run(a.bindAdr)
}

func (a *Application) registerHTTPMappings(engine *gin.Engine) {
	ui, err := rice.FindBox("ui-dist")
	if err == nil {
		engine.StaticFS("/ui", ui.HTTPBox())

		engine.GET("/", a.redirectToUi)
	} else {
		log.Warning("This is a Development Binary. This Means the WebApplication is not available on <URL>/ui")
	}
	//Status Informations
	engine.GET("/api/status", a.getStatus)
	engine.GET("/api/stats", a.getStats)
	engine.GET("/api/config", a.getConfiguration)

	//Spectra Retrieval
	engine.GET("/api/spectra", a.getValidDates)
	engine.GET("/api/spectra/:year/:month/:day", a.listSpectraForDay)
	engine.GET("/api/spectrum/:timestamp", a.getJsonSpectrum)
	engine.GET("/api/spectrum/:timestamp/img", a.getRenderedSpectrum)

	//Configuration Operations
	engine.POST("/api/config", a.updateConfig)
	engine.POST("/api/control/start", a.startCollecting)
	engine.POST("/api/control/stop", a.stopCollecting)
}

func (a *Application) redirectToUi(ctx *gin.Context) {
	ctx.Redirect(301, "/ui")
}
