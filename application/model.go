package application

import "github.com/c-mueller/fritzbox-spectrum-logger/config"

type Application struct {
    config  config.Configuration
    bindAdr string
}
