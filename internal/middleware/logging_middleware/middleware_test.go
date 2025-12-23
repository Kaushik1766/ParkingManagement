package loggingmiddleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	type args struct {
		next func(w http.ResponseWriter, r *http.Request)
	}
	tests := []struct {
		name       string
		args       args
		request    *http.Request
		wantLogs   []string
		nextCalled bool
	}{
		{
			name: "successful request logging",
			args: args{
				next: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("OK"))
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/test", nil),
			wantLogs: []string{
				"Started GET /test",
				"Request Headers:",
				"Completed GET /test in",
			},
			nextCalled: true,
		},
		{
			name: "POST request with headers",
			args: args{
				next: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(`{"name":"test"}`))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer token123")
				return req
			}(),
			wantLogs: []string{
				"Started POST /api/users",
				"Request Headers:",
				"Content-Type",
				"Authorization",
				"Completed POST /api/users in",
			},
			nextCalled: true,
		},
		{
			name: "request with query parameters",
			args: args{
				next: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/search?q=test&page=1", nil),
			wantLogs: []string{
				"Started GET /search",
				"Request Headers:",
				"Completed GET /search in",
			},
			nextCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			w := httptest.NewRecorder()

			middleware := LoggingMiddleware(tt.args.next)
			middleware(w, tt.request)

			logOutput := buf.String()

			if tt.nextCalled && w.Code == 0 {
				t.Error("Expected next handler to be called, but response writer was not used")
			}

			for _, wantLog := range tt.wantLogs {
				if !strings.Contains(logOutput, wantLog) {
					t.Errorf("Expected log to contain %q, but got: %s", wantLog, logOutput)
				}
			}

			if !strings.Contains(logOutput, "in ") {
				t.Error("Expected log to contain timing information")
			}
		})
	}
}

func TestLoggingMiddleware_NextHandlerPanic(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	middleware := LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic to bubble up through middleware")
		}

		logOutput := buf.String()
		if !strings.Contains(logOutput, "Started GET /panic") {
			t.Error("Expected start log even with panic")
		}
	}()

	middleware(w, req)
}
