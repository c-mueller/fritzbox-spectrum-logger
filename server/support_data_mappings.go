package server

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func (a *Application) getSupportTimestamps(ctx *gin.Context) {
	ctx.JSON(200, a.repo.ListSupportDataEntries())
}

func (a *Application) getSupportData(ctx *gin.Context) {
	timestampStr := ctx.Param("timestamp")
	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		ctx.Data(400, "text/plain", []byte("Invalid Timestamp"))
		return
	}

	data, err := a.repo.GetSupportData(int(ts))
	if err != nil {
		ctx.Data(404, "text/plain", []byte("Timestamp not found!"))
		return
	}

	ctx.Data(200, "text/plain", data)
}
