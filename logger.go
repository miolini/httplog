package httplog

import (
	"github.com/dustin/go-humanize"
	"github.com/miolini/datacounter"
	"log"
	"net/http"
	"os"
	"time"
)

// Logger simple func for wrapping http.Handler with standart stderr logger
func Logger(h http.Handler) http.Handler {
	return LoggerWithImpl(h, log.New(os.Stderr, "", 0))
}

// LoggerWithImpl func for wrapping http.Handler and custom logger implementation
func LoggerWithImpl(h http.Handler, logimpl *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wcounter := datacounter.NewResponseWriterCounter(w)
		defer func(t time.Time) {
			log.Printf("[%.2fms %s] %s %s", time.Since(t).Seconds()*1000, humanize.Bytes(wcounter.Count()), r.Method, r.URL.Path)
		}(time.Now())
		h.ServeHTTP(wcounter, r)
	})
}
