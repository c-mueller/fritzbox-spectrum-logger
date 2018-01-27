package application

import (
    "github.com/c-mueller/fritzbox-spectrum-logger/config"
    "github.com/c-mueller/fritzbox-spectrum-logger/repository"
    "github.com/c-mueller/fritzbox-spectrum-logger/fritz"
    "time"
)

type Application struct {
    config    config.Configuration
    bindAdr   string
    repo      *repository.Repository
    session   *fritz.Session
    state     ApplicationState
    startTime time.Time
}

type StatusResponse struct {
    State  string `json:"state"`
    Uptime int64  `json:"uptime"`
}
