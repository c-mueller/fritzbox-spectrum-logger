package server

import "github.com/prometheus/client_golang/prometheus"

var requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "request_counter",
	Help:      "Number of requests processed in since application launch",
}, []string{"path"})

var processingTimeHistogram = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "request_processing_time_summary",
	Help:      "Processing time summary",
}, []string{"path"})

func init() {
	prometheus.MustRegister(requestCounter, processingTimeHistogram)
}
