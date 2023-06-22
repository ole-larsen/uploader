package util

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/ole-larsen/uploader/service"
	"github.com/ole-larsen/uploader/service/log"
)

func SetupPrometheusHandler(handler http.Handler) http.Handler {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := service.NewPrometheusResponseWriter(w)
		if lrw != nil && r != nil {
			handler.ServeHTTP(lrw, r)
			statusCode := lrw.StatusCode
			duration := time.Since(start)
			logger := log.NewLogger()
			logger.Infoln(r.URL.String(), r.Method, fmt.Sprintf("%d", statusCode), duration.Seconds())
		}
	})
	return h
}

func SetupCorsHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("origin"))
		w.Header().Add("Content-Type", "application/json, multipart/form-data, application/x-www-form-urlencoded")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Token")
		w.Header().Add("Access-Control-Expose-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func SetupCsrfHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrf, err := generateNonce()
		if err != nil {
			panic(err)
		}
		expiration := time.Now().Add(time.Hour)
		cookie := http.Cookie{Name: "_csrf", Value: csrf, Expires: expiration, Path: "/"}
		http.SetCookie(w, &cookie)
		handler.ServeHTTP(w, r)
	})
}

func generateNonce() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	return fmt.Sprintf("%x", b), err
}
