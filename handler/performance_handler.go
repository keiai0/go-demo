package handler

import (
	"net/http"

	"github.com/keiai0/go-demo/demo"
)

func PerformanceHandler(w http.ResponseWriter, r *http.Request) {
	demo.PerformanceDemo()
	_, _ = w.Write([]byte("performance comparison!"))
}

func PerformanceProfHandler(w http.ResponseWriter, r *http.Request) {
	demo.PerformanceProfDemo()
	_, _ = w.Write([]byte("performance comparison!"))
}

func PerformanceTraceHandler(w http.ResponseWriter, r *http.Request) {
	demo.PerformanceTraceDemo()
	_, _ = w.Write([]byte("performance comparison!"))
}
