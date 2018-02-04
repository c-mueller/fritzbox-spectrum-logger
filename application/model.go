package application

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/config"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"time"
)

type Application struct {
	config            config.Configuration
	bindAdr           string
	repo              *repository.Repository
	session           *fritz.Session
	state             ApplicationState
	startTime         time.Time
	updateTicker      *time.Ticker
	sessionLogCounter int64
}

type StatusResponse struct {
	State         string `json:"state"`
	Uptime        int64  `json:"uptime"`
	SpectrumCount int64  `json:"spectrum_count"`
}

type InfoResponse struct {
	State   string `json:"state"`
	Message string `json:"message"`
}
