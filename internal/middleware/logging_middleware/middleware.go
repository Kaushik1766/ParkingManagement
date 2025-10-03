package loggingmiddleware

import (
	"net/http"
)

func LoggingMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// <-time.After(2 * time.Second)
		// start := time.Now()
		// log.Printf("\n\nStarted %s %s\n", r.Method, r.URL.Path)
		// log.Printf("Request Headers: %v\n", r.Header)
		next(w, r)
		// log.Printf("Completed %s %s in %v\n\n", r.Method, r.URL.Path, time.Since(start))
	}
}
