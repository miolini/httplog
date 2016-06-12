package httplog

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/miolini/datacounter"
	"io"
	"net/http"
	"os"
	"time"
)

type FormatterFunc func(w io.Writer, r *http.Request, sentBytes uint64, elapsedTime time.Duration)

var DefaultFormatter FormatterFunc = func(w io.Writer, r *http.Request, sentBytes uint64, elapsedTime time.Duration) {
	fmt.Fprintf(w, "[%.2fms %s] %s %s", elapsedTime.Seconds()*1000, humanize.Bytes(sentBytes), r.Method, r.URL.Path)
}

// Logger simple func for wrapping http.Handler and log to stderr
func Logger(h http.Handler) http.Handler {
	return LoggerWithWriterAndFormatter(h, os.Stderr, DefaultFormatter)
}

// LoggerWithWriter func for wrapping http.Handler and user defined output io.Writer
func LoggerWithWriter(h http.Handler, w io.Writer) http.Handler {
	return LoggerWithWriterAndFormatter(h, w, DefaultFormatter)
}

// LoggerWithWriterAndFormatter func for wrapping http.Handler and user defined output io.Writer and FormatterFunc
func LoggerWithWriterAndFormatter(h http.Handler, lw io.Writer, formatter FormatterFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wcounter := datacounter.NewResponseWriterCounter(w)
		defer func(t time.Time) {
			formatter(lw, r, wcounter.Count(), time.Since(t))
		}(time.Now())
		h.ServeHTTP(wcounter, r)
	})
}
