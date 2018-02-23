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
	state             APIState
	startTime         time.Time
	updateTicker      *time.Ticker
	sessionLogCounter int64
	latest            *LatestSpectrumResponse
}

type StatusResponse struct {
	State         string                   `json:"state"`
	Uptime        int64                    `json:"uptime"`
	SpectrumCount int64                    `json:"spectrum_count"`
	Latest        *LatestSpectrumResponse  `json:"latest"`
	Stats         *repository.SpectraStats `json:"stats"`
}

type InfoResponse struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

type KeysResponse struct {
	Keys             repository.SpectraKeys `json:"keys"`
	RequestTimestamp int64                  `json:"timestamp"`
}

type TimestampResponse struct {
	Timestamps       []int64                `json:"timestamps"`
	Key              repository.SpectrumKey `json:"requested_day"`
	RequestTimestamp int64                  `json:"timestamp"`
}

type LatestSpectrumResponse struct {
	Key       repository.SpectrumKey `json:"date"`
	Timestamp int64                  `json:"timestamp"`
}
