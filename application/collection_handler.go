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
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/gin-gonic/gin"
	"time"
)

func (a *Application) startCollecting(ctx *gin.Context) {
	if a.state != LOGGING {
		log.Info("Starting Spectrum Logging")
		go a.collectionHandler()
		ctx.JSON(200, InfoResponse{
			State:   LOGGING.String(),
			Message: "Collection Started!",
		})
	} else {
		ctx.JSON(400, InfoResponse{
			State:   a.state.String(),
			Message: "Already Running!",
		})
	}
}

func (a *Application) stopCollecting(ctx *gin.Context) {
	if a.state == LOGGING {
		log.Info("Stopping Collector")
		a.updateTicker.Stop()
		log.Info("Collection Stopped!")
		a.state = IDLE
		ctx.JSON(200, InfoResponse{
			State:   a.state.String(),
			Message: "Collection Stopped! State Change will occur soon!",
		})
	} else {
		ctx.JSON(400, InfoResponse{
			State:   a.state.String(),
			Message: "Not Logging. Cannot Stop.",
		})
	}
}

func (a *Application) collectionHandler() {
	log.Info("Launching Collection Handler...")
	updateInterval := time.Duration(a.config.UpdateInterval) * time.Second
	a.updateTicker = time.NewTicker(updateInterval)
	a.state = LOGGING

	log.Info("Logging into Fritz!Box")
	cred := a.config.Credentials
	a.session = fritz.NewClient(cred.Endpoint, cred.Username, cred.Password)
	err := a.session.Login()
	if err != nil {
		log.Error("Login failed: ", err)
		a.updateTicker.Stop()
		a.state = ERROR
		return
	}
	log.Info("Logged In!")

	for range a.updateTicker.C {
		log.Info("Collecting...")
		err := a.collect()
		if err != nil {
			log.Errorf("Could not download Spectrum. Aborting. Error: %v", err)
			a.state = ERROR
			a.updateTicker.Stop()
			return
		}
	}
}

func (a *Application) collect() error {
	spec, err := a.session.GetSpectrum()
	if err != nil {
		return err
	}

	err = a.repo.Insert(spec)

	if err != nil {
		return err
	}

	a.sessionLogCounter++

	return nil
}
