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

	log.Debug("Launched HTTP Server")
	return engine.Run(a.bindAdr)
}

func (a *Application) registerHTTPMappings(engine *gin.Engine) {
	//Status Informations
	engine.GET("/api/status", a.getStatus)
	engine.GET("/api/config", a.getConfiguration)

	//Spectra Retrieval
	engine.GET("/api/spectra", a.getValidDates)
	engine.GET("/api/spectra/:year/:month/:day", a.listSpectraForDay)
	engine.GET("/api/spectra/:year/:month/:day/:timestamp", a.getRawSpectrum)
	engine.GET("/api/spectra/:year/:month/:day/:timestamp/img", a.getRenderedSpectrum)

	//Configuration Operations
	engine.POST("/api/config", a.updateConfig)
	engine.POST("/api/control/start", a.startCollecting)
	engine.POST("/api/control/stop", a.stopCollecting)
}
