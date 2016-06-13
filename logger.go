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

// FormatterFunc user defined formatter func
type FormatterFunc func(w io.Writer, r *http.Request, sentBytes uint64, elapsedTime time.Duration)

// DefaultFormatter default formatter func
var DefaultFormatter FormatterFunc = func(w io.Writer, r *http.Request, sentBytes uint64, elapsedTime time.Duration) {
	fmt.Fprintf(w, "%s [%.2fms %s] %s %s\n", time.Now().Format("2005-Mo-2 15:04:05"),
		elapsedTime.Seconds()*1000, humanize.Bytes(sentBytes), r.Method, r.URL.Path)
}

// Logger simple func for wrapping http.Handler and log to stderr
func Logger(h http.Handler) http.Handler {
	return LoggerWithWriterAndFormatter(h, os.Stderr, DefaultFormatter)
}

// LoggerWithFormatter wrapping func with user defined formatter
func LoggerWithFormatter(h http.Handler, formatter FormatterFunc) http.Handler {
	return LoggerWithWriterAndFormatter(h, os.Stderr, formatter)
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
