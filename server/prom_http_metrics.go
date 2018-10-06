package server

import "github.com/prometheus/client_golang/prometheus"

var requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "request_counter",
	Help:      "Number of requests processed in since application launch",
}, []string{"path"})

var processingTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "",
	Subsystem: "fsl",
	Name:      "request_processing_time",
	Help:      "Processing time Histogram",
	Buckets:   []float64{1, 2, 5, 10, 20, 50, 100, 250, 500, 1000, 2500,},
}, []string{"path"})

func init() {
	prometheus.MustRegister(requestCounter, processingTimeHistogram)
}
