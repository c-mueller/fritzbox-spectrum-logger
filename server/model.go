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

type ComparisonRequest struct {
	Timestamps []int64 `json:"timestamps"`
}

type StatusResponse struct {
	State  string `json:"state"`
	Uptime int64  `json:"uptime"`
}

type NeighboursResponse struct {
	PreviousTimestamp int64 `json:"previous_timestamp"`
	NextTimestamp     int64 `json:"next_timestamp"`
	RequestTimestamp  int64 `json:"request_timestamp"`
}

type StatResponse struct {
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
