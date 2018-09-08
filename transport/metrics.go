package transport

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
)

func metricsMiddleware(key string, metricsService *metrics.Metrics) func(next http.Handler) http.Handler {
	metricsKeys := make(map[int][]string)
	var m sync.RWMutex
	getMKey := func(code int) []string {
		m.RLock()
		k, ok := metricsKeys[code]
		m.RUnlock()
		if ok {
			return k
		}

		m.Lock()
		k = []string{"http_req", key, strconv.Itoa(code)}
		metricsKeys[code] = k
		m.Unlock()

		return k
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			rw := &metricedResponseWriter{ResponseWriter: w}
			timeStart := time.Now()

			next.ServeHTTP(rw, req)

			statusCode := rw.statusCode
			if !rw.codeWritten {
				statusCode = 200
			}
			metricsService.MeasureSince(getMKey(statusCode), timeStart)
		})
	}
}

type metricedResponseWriter struct {
	http.ResponseWriter

	codeWritten bool
	statusCode  int
}

func (w *metricedResponseWriter) WriteHeader(code int) {
	if !w.codeWritten {
		w.codeWritten = true
		w.statusCode = code
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *metricedResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w *metricedResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
