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
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func (a *Application) getValidDates(ctx *gin.Context) {
	keys, err := a.repo.GetAllSpectrumKeys()
	if err != nil {
		sendError(ctx, 500, "Failed to Retrieve Keys: %s", err.Error())
		return
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
		sendError(ctx, 404, "A Invalid Key was requested: %s", key.String())
		return
	}
	timestamps, err := a.repo.GetTimestampsForSpectrumKey(key)
	if err != nil {
		sendError(ctx, 404, "Spectra Retrieval Failed: %s", err.Error())
		return
	}
	ctx.JSON(200, TimestampResponse{
		Timestamps:       timestamps,
		Key:              key,
		RequestTimestamp: time.Now().Unix(),
	})
}

func (a *Application) getNeighbours(ctx *gin.Context) {
	timestampString := ctx.Param("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	key := repository.GetFromTimestamp(timestamp)
	if !key.IsValid() || err != nil {
		sendError(ctx, 404, "A Invalid Key was requested: %s - Timestamp: %s", key.String(), timestampString)
		return
	}

	keys, err := a.repo.GetTimestampsForSpectrumKey(key)
	if err != nil {
		sendError(ctx, 404, "Retrieving the Neighbours for the Timestamp %d has failed", timestamp)
		return
	}

	index := keys.Search(timestamp)
	if keys[index] != timestamp {
		sendError(ctx, 404, "No spectrum found with timestamp %d", timestamp)
		return
	}
	previous, next := int64(index-1), int64(index+1)
	if previous < 0 {
		previous = -1
	} else {
		previous = keys[previous]
	}
	if next >= int64(keys.Len()) {
		next = -1
	} else {
		next = keys[next]
	}
	ctx.JSON(200, NeighboursResponse{
		PreviousTimestamp: previous,
		NextTimestamp:     next,
		RequestTimestamp:  time.Now().Unix(),
	})
}

func (a *Application) getJsonSpectrum(ctx *gin.Context) {
	timestampString := ctx.Param("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	key := repository.GetFromTimestamp(timestamp)
	if !key.IsValid() || err != nil {
		sendError(ctx, 404, "A Invalid Key was requested: %s - Timestamp: %s", key.String(), timestampString)
		return
	}
	spectrum, err := a.repo.GetSpectrumForTimestamp(timestamp)
	if err != nil {
		sendError(ctx, 404, "Spectra Retrieval Failed: %s", err.Error())
		return
	}
	ctx.JSON(200, spectrum)
}

func (a *Application) getRenderedSpectrum(ctx *gin.Context) {
	scaled := ctx.Query("scaled")
	timestampString := ctx.Param("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	key := repository.GetFromTimestamp(timestamp)
	if !key.IsValid() || err != nil {
		sendError(ctx, 404, "A Invalid Key was requested: %s - Timestamp: %s", key.String(), timestampString)
		return
	}
	spectrum, err := a.repo.GetSpectrumForTimestamp(timestamp)
	if err != nil {
		sendError(ctx, 404, "Spectra Retrieval Failed: %s", err.Error())
		return
	}
	image, err := spectrum.Render(scaled == "true")
	if err != nil {
		sendError(ctx, 500, "Spectra Rendering has Failed: %s", err.Error())
		return
	}
	ctx.Data(200, "image/png", image)
}
