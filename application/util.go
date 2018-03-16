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
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
)

func sendError(ctx *gin.Context, code int, format string, data ...interface{}) {
	message := fmt.Sprintf(format, data...)
	log.Error(message)
	ctx.String(code, fmt.Sprintf("%d: %s", code, message))
	return
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
