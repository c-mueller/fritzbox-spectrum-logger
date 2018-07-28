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
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
)

func (a *Application) getSpectraComparison(ctx *gin.Context) {
	requestBodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(400, "")
		return
	}

	var spectraList *ComparisonRequest
	err = json.Unmarshal(requestBodyBytes, &spectraList)
	if err != nil {
		ctx.String(400, "")
		return
	}

	spectraCollection := make(fritz.ComparisonSet, 0)
	for _, timestamp := range spectraList.Timestamps {
		spectrum, err := a.repo.GetSpectrumForTimestamp(timestamp)
		if err != nil {
			ctx.String(404, "")
			return
		}
		spectraCollection = append(spectraCollection, *spectrum)
	}

	imageBytes, err := spectraCollection.RenderComparison(false)
	if err != nil {
		ctx.String(500,"")
		return
	}
	ctx.Data(200, "image/png", imageBytes)
}
