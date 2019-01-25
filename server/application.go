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

package server

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/repository/reporegistry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"regexp"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/c-mueller/fritzbox-spectrum-logger/config"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"

	_ "github.com/c-mueller/fritzbox-spectrum-logger/repository/bolt"
	_ "github.com/c-mueller/fritzbox-spectrum-logger/repository/relational"
)

var log = logging.MustGetLogger("server")

func counterMiddleware() gin.HandlerFunc {

	timestampRegex, _ := regexp.Compile("/[0-9]{9,}")
	dateRegex, _ := regexp.Compile("/[0-9]{4}/1?[0-9]/[1-3]?[0-9]")

	return func(context *gin.Context) {
		requestPath := string(timestampRegex.ReplaceAll([]byte(context.Request.URL.Path), []byte("/TIMESTAMP")))
		requestPath = string(dateRegex.ReplaceAll([]byte(requestPath), []byte("/YYYY/MM/DD")))

		start := time.Now()

		context.Next()

		requestCounter.WithLabelValues(requestPath).Inc()
		processingTimeHistogram.WithLabelValues(requestPath).Observe(float64(time.Now().Sub(start).Nanoseconds()) / 1000000)
	}
}

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

func LaunchFromEnvironment() *Application {
	log.Debug("Loading configuration from the application environment")
	cfg, err := config.FromEnvironment()
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

func (a *Application) initRepository() error {
	log.Debug("Launching server...")
	log.Debug("Initializing repository (Datastore)")

	var repo repository.Repository
	var err error

	switch a.config.DatabaseMode {
	case config.DatabaseModeBolt:
		log.Debugf("Using BoltDB based datastore located at %q...", a.config.DatabasePath)
		repo, err = reporegistry.GetForName("bolt").BuildRepository(a.config.DatabaseCompression, a.config.DatabasePath)
		break
	case config.DatabaseModeSQLite:
		log.Debugf("Using SQLite based datastore located at %q...", a.config.DatabasePath)
		repo, err = reporegistry.GetForName("sqlite").BuildRepository(a.config.DatabaseCompression, a.config.DatabasePath)
		break
	}

	if err != nil {
		return err
	}
	a.repo = repo
	return nil
}

func (a *Application) Listen() error {
	if err := a.initRepository(); err != nil {
		return err
	}

	defer a.repo.Close()

	log.Debug("Registering HTTP mappings")
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(counterMiddleware())
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
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
	// Register Ui Mappings (if present)
	ui, err := rice.FindBox("ui-dist")
	if err == nil {
		engine.StaticFS("/ui", ui.HTTPBox())

		engine.GET("/", a.redirectToUi)
	} else {
		log.Warning("This is a Development Binary. This Means the WebApplication is not available on <URL>/ui")
	}

	//Prometheus Metrics
	prometheusHandler := promhttp.Handler()
	engine.GET("/metrics", func(c *gin.Context) {
		prometheusHandler.ServeHTTP(c.Writer, c.Request)
	})

	//Status Informations
	engine.GET("/api/status", a.getStatus)
	engine.GET("/api/stats", a.getStats)
	engine.GET("/api/config", a.getConfiguration)

	//Spectra Listing
	engine.GET("/api/spectra", a.getValidDates)
	engine.GET("/api/spectra/:year/:month/:day", a.listSpectraForDay)

	//Spectrum Retrieval
	engine.GET("/api/spectrum/:timestamp", a.getJsonSpectrum)
	engine.GET("/api/spectrum/:timestamp/img", a.getRenderedSpectrum)
	engine.GET("/api/spectrum/:timestamp/neighbours", a.getNeighbours)
	engine.GET("/api/spectrum/:timestamp/info", a.getConnectionInformation)
	engine.GET("/api/spectrum/:timestamp/info.json", a.getParsedConnectionInformation)

	engine.GET("/api/sd", a.getSupportTimestamps)
	engine.GET("/api/sd/:timestamp", a.getSupportData)

	//Spectrum Comparsion
	engine.POST("/api/comparison", a.getSpectraComparison)

	//Configuration Operations
	engine.POST("/api/config", a.updateConfig)
	engine.POST("/api/control/start", a.startCollecting)
	engine.POST("/api/control/stop", a.stopCollecting)

	engine.GET("/api/db/backup", a.repo.Backup())
}

func (a *Application) redirectToUi(ctx *gin.Context) {
	ctx.Redirect(301, "/ui")
}
