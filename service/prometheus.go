package service

import (
	"net/http"
)

//type DBStats = sql.DBStats
//
//type Metrics interface {
//	SetDBStats(dbstats DBStats)
//	HTTPHandler() http.Handler
//}

type prometheusResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewPrometheusResponseWriter(w http.ResponseWriter) *prometheusResponseWriter {
	return &prometheusResponseWriter{w, http.StatusOK}
}

func (w *prometheusResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}
