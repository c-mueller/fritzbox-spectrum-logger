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
	"strconv"
)

func (a *Application) getJsonSpectrum(ctx *gin.Context) {
	spectrum := a.getSpectrumFromParameters(ctx)
	if spectrum == nil {
		return
	}

	ctx.JSON(200, spectrum)
}

func (a *Application) getRenderedSpectrum(ctx *gin.Context) {
	scaled := ctx.Query("scaled")
	spectrum := a.getSpectrumFromParameters(ctx)

	if spectrum == nil {
		return
	}

	image, err := spectrum.Render(scaled == "true")
	if err != nil {
		sendError(ctx, 500, "Spectra Rendering has Failed: %s", err.Error())
		return
	}
	ctx.Data(200, "image/png", image)
}

func (a *Application) getParsedConnectionInformation(ctx *gin.Context) {
	spectrum := a.getSpectrumFromParameters(ctx)

	if spectrum == nil {
		ctx.Data(404, "application/json", []byte("{}"))
		return
	}

	conInfo, err := spectrum.GetConnectionInformation()
	if err != nil {
		ctx.Data(400, "application/json", []byte("{}"))
	} else {
		ctx.JSON(200, conInfo)
	}
}

func (a *Application) getConnectionInformation(ctx *gin.Context) {
	spectrum := a.getSpectrumFromParameters(ctx)

	if spectrum == nil {
		return
	}

	ctx.Data(200, "text/html", []byte(spectrum.ConnectionInformation))
}

func (a *Application) getSpectrumFromParameters(ctx *gin.Context) *fritz.Spectrum {
	timestampString := ctx.Param("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	key := repository.GetFromTimestamp(timestamp)
	if !key.IsValid() || err != nil {
		sendError(ctx, 404, "A Invalid Key was requested: %s - Timestamp: %s", key.String(), timestampString)
		return nil
	}
	spectrum, err := a.repo.GetSpectrum(timestamp)
	if err != nil {
		sendError(ctx, 404, "Spectra Retrieval Failed: %s", err.Error())
		return nil
	}
	return spectrum
}
