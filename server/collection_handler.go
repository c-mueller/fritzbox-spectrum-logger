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
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/gin-gonic/gin"
	"time"
)

func (a *Application) startCollecting(ctx *gin.Context) {
	if a.state != LOGGING {
		a.startLogging()
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
		a.stopLogging()
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

func (a *Application) startLogging() {
	log.Info("Starting Spectrum Logging")
	go a.collectionHandler()
}

func (a *Application) stopLogging() {
	log.Info("Stopping Collector")
	a.updateTicker.Stop()
	log.Info("Collection Stopped!")
	a.state = IDLE
}

func (a *Application) collectionHandler() {
	log.Info("Launching Collection Handler...")
	updateInterval := time.Duration(a.config.UpdateInterval) * time.Second
	a.updateTicker = time.NewTicker(updateInterval)
	a.state = LOGGING

	if !a.renewSession(true) {
		return
	}

	renewalAttempts := 0
	spectrumLoadErrors := 0

	for range a.updateTicker.C {

		if a.session.TokenTimedOut(int64(a.config.SessionRefreshInterval)) {
			log.Info("Renewing Session...")
			renewalSuccessful := a.renewSession(false)
			if !renewalSuccessful && renewalAttempts < a.config.SessionRenewalAttemptCount {
				renewalAttempts++
				continue
			} else if !renewalSuccessful && renewalAttempts >= a.config.SessionRenewalAttemptCount {
				log.Errorf("Logging Stopped because %d Login attempts in a row have failed!", spectrumLoadErrors)
				a.updateTicker.Stop()
				a.state = ERROR
				return
			}
		}
		renewalAttempts = 0

		log.Debug("Downloading Spectrum...")
		err := a.collect()
		if err != nil {
			spectrumLoadErrors++

			failCountVec.WithLabelValues().Inc()

			log.Errorf("Could not download Spectrum. Aborting. Error: %v", err)
			if spectrumLoadErrors >= a.config.MaxDownloadFails {
				log.Errorf("Logging Stopped because %d Download attempts in a row have failed!", spectrumLoadErrors)
				a.updateTicker.Stop()
				a.state = ERROR
				return
			}
		} else {
			spectrumLoadErrors = 0
		}
	}
}

func (a *Application) renewSession(failOnError bool) bool {
	log.Info("Logging into Fritz!Box")
	cred := a.config.Credentials
	session := fritz.NewClient(cred.Endpoint, cred.Username, cred.Password)
	err := session.Login()
	if err != nil {
		log.Error("Login failed: ", err)
		if failOnError {
			a.updateTicker.Stop()
			a.state = ERROR
		}
		return false
	}
	a.session = session
	log.Info("Logged In!")
	return true
}

func (a *Application) collect() error {
	spec, err := a.session.GetSpectrum()
	if err != nil {
		return err
	}

	go a.updatePrometheus(spec)

	sk := repository.GetFromTimestamp(spec.Timestamp)
	a.latest = &LatestSpectrumResponse{
		Key:       sk,
		Timestamp: spec.Timestamp,
	}

	err = a.repo.Insert(spec)
	if err != nil {
		return err
	}

	a.sessionLogCounter++

	return nil
}

func (a *Application) updatePrometheus(spectrum *fritz.Spectrum) {
	conInfo, err := spectrum.GetConnectionInformation()
	if err != nil {
		return
	}

	streams := []fritz.ConnectionTransmissionDirection{conInfo.Downstream, conInfo.Upstream}
	name := []string{"DOWNSTREAM", "UPSTREAM"}

	for index, streamInfo := range streams {
		streamName := name[index]

		maxDataRateVec.WithLabelValues(streamName).Set(float64(streamInfo.MaximumDataRate))
		minDataRateVec.WithLabelValues(streamName).Set(float64(streamInfo.MinimumDataRate))
		capacityVec.WithLabelValues(streamName).Set(float64(streamInfo.Capacity))
		currentDataRateVec.WithLabelValues(streamName).Set(float64(streamInfo.CurrentDataRate))

		lineLatencyVec.WithLabelValues(streamName).Set(float64(streamInfo.Latency))
		inpValueVec.WithLabelValues(streamName).Set(float64(streamInfo.INPValue))
		snrVec.WithLabelValues(streamName).Set(streamInfo.SNMargin)
		attenuationVec.WithLabelValues(streamName).Set(streamInfo.LineAttenuation)

		errorVec.WithLabelValues(streamName).Set(streamInfo.Errors.SecondsWithErrors)
		manyErrorVec.WithLabelValues(streamName).Set(streamInfo.Errors.SecondsWithManyErrors)
		errorsPerMinVec.WithLabelValues(streamName).Set(streamInfo.Errors.ErrorsPerMinute)
		errorsLast15MinVec.WithLabelValues(streamName).Set(streamInfo.Errors.ErrorsLast15Min)
	}
}
