package application

import (
    "github.com/c-mueller/fritzbox-spectrum-logger/config"
    "github.com/c-mueller/fritzbox-spectrum-logger/repository"
    "github.com/c-mueller/fritzbox-spectrum-logger/fritz"
)

type Application struct {
    config  config.Configuration
    bindAdr string
    repo    *repository.Repository
    session *fritz.Session
}
