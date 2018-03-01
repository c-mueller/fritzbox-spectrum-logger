// Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger)
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
	"encoding/json"
	"github.com/c-mueller/fritzbox-spectrum-logger/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

func (a *Application) getStatus(ctx *gin.Context) {
	ctx.JSON(200, StatusResponse{
		State:  a.state.String(),
		Uptime: int64(time.Since(a.startTime).Seconds()),
	})
}

func (a *Application) getStats(ctx *gin.Context) {
	stats, err := a.repo.GetStatistics()
	if err != nil {
		ctx.String(500, "")
		return
	}
	if a.latest == nil {
		ctx.JSON(200, StatResponse{
			SpectrumCount: a.sessionLogCounter,
			Stats:         stats,
		})
	} else {
		ctx.JSON(200, StatResponse{
			SpectrumCount: a.sessionLogCounter,
			Latest:        a.latest,
			Stats:         stats,
		})
	}
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
