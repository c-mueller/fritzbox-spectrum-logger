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
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func (a *Application) getValidDates(ctx *gin.Context) {
	keys, err := a.repo.GetAllSpectrumKeys()
	if err != nil {
		log.Errorf("Failed to Retrieve Keys: %v", err)
		ctx.String(500, "")
	}
	response := KeysResponse{
		Keys:             keys,
		RequestTimestamp: time.Now().Unix(),
	}
	ctx.JSON(200, response)
}

func (a *Application) listSpectraForDay(ctx *gin.Context) {
	key := getSpectrumKeyFormContext(ctx)
	if !key.IsValid() {
		log.Errorf("A Invalid Key was requested: %s", key.String())
		ctx.String(404, "")
		return
	}
	spectra, err := a.repo.GetSpectraForSpectrumKey(key)
	if err != nil {
		log.Errorf("Spectra Retrieval failed: %s", err)
		ctx.String(404, "")
		return
	}
	timestamps := make([]int64, 0)
	for _, v := range spectra {
		timestamps = append(timestamps, v.Timestamp)
	}
	ctx.JSON(200, TimestampResponse{
		Timestamps:       timestamps,
		Key:              key,
		RequestTimestamp: time.Now().Unix(),
	})
}

func (a *Application) getJsonSpectrum(ctx *gin.Context) {
	key := getSpectrumKeyFormContext(ctx)
	timestampString := ctx.Param("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	if !key.IsValid() || err != nil {
		log.Errorf("A Invalid Key was requested: %s - Timestamp: %s", key.String(), timestampString)
		ctx.String(404, "")
		return
	}
	spectrum, err := a.repo.GetSpectrumBySpectrumKey(&key, timestamp)
	if err != nil {
		log.Errorf("Spectra Retrieval failed: %s", err)
		ctx.String(404, "")
		return
	}
	ctx.JSON(200, spectrum)
}

func (a *Application) getRenderedSpectrum(ctx *gin.Context) {
	key := getSpectrumKeyFormContext(ctx)
	timestampString := ctx.Param("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	if !key.IsValid() || err != nil {
		log.Errorf("A Invalid Key was requested: %s - Timestamp: %s", key.String(), timestampString)
		ctx.String(404, "")
		return
	}
	spectrum, err := a.repo.GetSpectrumBySpectrumKey(&key, timestamp)
	if err != nil {
		log.Errorf("Spectra Retrieval failed: %s", err)
		ctx.String(404, "")
		return
	}
	image, err := spectrum.Render()
	if err != nil {
		log.Errorf("Rendering The spectrum failed", err)
		ctx.String(500, "")
		return
	}
	ctx.Data(200, "image/png", image)
}

func getSpectrumKeyFormContext(ctx *gin.Context) repository.SpectrumKey {
	year := ctx.Param("year")
	month := ctx.Param("month")
	day := ctx.Param("day")
	key := repository.SpectrumKey{
		Year:  year,
		Month: month,
		Day:   day,
	}
	return key
}
