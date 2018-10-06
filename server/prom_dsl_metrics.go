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

import "github.com/prometheus/client_golang/prometheus"

var maxDataRateVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_dslam_max_data_rate",
	Help:      "Maximum datarate of the DSLAM",
}, []string{"direction"})

var minDataRateVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_dslam_min_data_rate",
	Help:      "Minimum datarate of the DSLAM",
}, []string{"direction"})

var capacityVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_capacity",
	Help:      "Capacity of the DSL Connection",
}, []string{"direction"})

var currentDataRateVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_current_data_rate",
	Help:      "Data rate of the connection",
}, []string{"direction"})

var lineLatencyVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_latency",
	Help:      "Latency in MS of the connection",
}, []string{"direction"})

var inpValueVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_inp",
	Help:      "INP value of the connection",
}, []string{"direction"})

var snrVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_snr",
	Help:      "SNR Value in dB of the connection",
}, []string{"direction"})

var attenuationVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_attenuation",
	Help:      "Attenuation in dB of the connection",
}, []string{"direction"})

var errorVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_error_seconds",
	Help:      "Errors ES",
}, []string{"direction"})

var manyErrorVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_many_error_seconds",
	Help:      "Errors SES",
}, []string{"direction"})

var errorsPerMinVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_errors_per_minute",
	Help:      "Errors per minute",
}, []string{"direction"})

var errorsLast15MinVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_line_errors_last_15_min",
	Help:      "Errors in the last 15 minutes",
}, []string{"direction"})

var failCountVec = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "fritzbox_collection_failures",
	Help:      "The count of failed attempts to download a spectrum",
}, []string{})

func init() {
	prometheus.MustRegister(maxDataRateVec, minDataRateVec)
	prometheus.MustRegister(capacityVec, currentDataRateVec)
	prometheus.MustRegister(lineLatencyVec, inpValueVec)
	prometheus.MustRegister(snrVec, attenuationVec)
	prometheus.MustRegister(failCountVec)
	prometheus.MustRegister(errorVec, manyErrorVec, errorsPerMinVec, errorsLast15MinVec)
}
