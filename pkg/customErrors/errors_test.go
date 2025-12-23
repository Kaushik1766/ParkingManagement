package customerrors

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestBadRequestError(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		err WebError
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   WebError
	}{
		{
			name: "bad request with validation error",
			args: args{
				w: httptest.NewRecorder(),
				err: WebError{
					Message: "validation failed for user kaushik",
					Code:    400,
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: WebError{
				Message: "validation failed for user kaushik",
				Code:    400,
			},
		},
		{
			name: "bad request with missing email",
			args: args{
				w: httptest.NewRecorder(),
				err: WebError{
					Message: "email kaushik@a.com is required",
					Code:    400,
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: WebError{
				Message: "email kaushik@a.com is required",
				Code:    400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			recorder := tt.args.w.(*httptest.ResponseRecorder)
			BadRequestError(tt.args.w, tt.args.err)

			// Check status code
			if recorder.Code != tt.expectedStatus {
				t.Errorf("BadRequestError() status = %v, want %v", recorder.Code, tt.expectedStatus)
			}

			// Check content type
			expectedContentType := "application/json"
			if contentType := recorder.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("BadRequestError() content-type = %v, want %v", contentType, expectedContentType)
			}

			// Check response body
			var actualBody WebError
			if err := json.Unmarshal(recorder.Body.Bytes(), &actualBody); err != nil {
				t.Errorf("BadRequestError() failed to unmarshal response: %v", err)
			}
			if actualBody != tt.expectedBody {
				t.Errorf("BadRequestError() body = %v, want %v", actualBody, tt.expectedBody)
			}

			// Check log output
			logOutput := buf.String()
			if !strings.Contains(logOutput, "Bad request:") {
				t.Errorf("BadRequestError() log output should contain 'Bad request:', got: %v", logOutput)
			}
		})
	}
}

func TestDisplayError(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "display user not found error",
			args: args{
				msg: "user kaushik not found",
			},
		},
		{
			name: "display authentication error",
			args: args{
				msg: "invalid credentials for kaushik@a.com",
			},
		},
		{
			name: "display admin access error",
			args: args{
				msg: "admin@a.com unauthorized access",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			// Redirect stdin to provide input for fmt.Scanln()
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write a newline to simulate user pressing Enter
			go func() {
				defer w.Close()
				if _, err := w.Write([]byte("\n")); err != nil {
					t.Logf("Error writing to pipe: %v", err)
				}
			}()
			defer func() { os.Stdin = oldStdin }()

			DisplayError(tt.args.msg)

			// Check log output
			logOutput := buf.String()
			expectedLog := "Error: " + tt.args.msg
			if !strings.Contains(logOutput, expectedLog) {
				t.Errorf("DisplayError() log output should contain '%v', got: %v", expectedLog, logOutput)
			}
		})
	}
}

func TestInternalServerError(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		err WebError
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   WebError
	}{
		{
			name: "internal server error database connection",
			args: args{
				w: httptest.NewRecorder(),
				err: WebError{
					Message: "database connection failed for user kaushik",
					Code:    500,
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: WebError{
				Message: "database connection failed for user kaushik",
				Code:    500,
			},
		},
		{
			name: "internal server error service unavailable",
			args: args{
				w: httptest.NewRecorder(),
				err: WebError{
					Message: "service temporarily unavailable",
					Code:    500,
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: WebError{
				Message: "service temporarily unavailable",
				Code:    500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			recorder := tt.args.w.(*httptest.ResponseRecorder)
			InternalServerError(tt.args.w, tt.args.err)

			// Check status code
			if recorder.Code != tt.expectedStatus {
				t.Errorf("InternalServerError() status = %v, want %v", recorder.Code, tt.expectedStatus)
			}

			// Check content type
			expectedContentType := "application/json"
			if contentType := recorder.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("InternalServerError() content-type = %v, want %v", contentType, expectedContentType)
			}

			// Check response body
			var actualBody WebError
			if err := json.Unmarshal(recorder.Body.Bytes(), &actualBody); err != nil {
				t.Errorf("InternalServerError() failed to unmarshal response: %v", err)
			}
			if actualBody != tt.expectedBody {
				t.Errorf("InternalServerError() body = %v, want %v", actualBody, tt.expectedBody)
			}

			// Check log output
			logOutput := buf.String()
			if !strings.Contains(logOutput, "Internal server error:") {
				t.Errorf("InternalServerError() log output should contain 'Internal server error:', got: %v", logOutput)
			}
		})
	}
}

func TestUnathorized_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "unauthorized error message",
			want: "user unauthorized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Unauthorized{}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnauthorizedError(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		err WebError
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   WebError
	}{
		{
			name: "unauthorized access for kaushik",
			args: args{
				w: httptest.NewRecorder(),
				err: WebError{
					Message: "user kaushik@a.com not authorized",
					Code:    401,
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: WebError{
				Message: "user kaushik@a.com not authorized",
				Code:    401,
			},
		},
		{
			name: "unauthorized admin access",
			args: args{
				w: httptest.NewRecorder(),
				err: WebError{
					Message: "admin@a.com requires elevated privileges",
					Code:    401,
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: WebError{
				Message: "admin@a.com requires elevated privileges",
				Code:    401,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			recorder := tt.args.w.(*httptest.ResponseRecorder)
			UnauthorizedError(tt.args.w, tt.args.err)

			// Check status code
			if recorder.Code != tt.expectedStatus {
				t.Errorf("UnauthorizedError() status = %v, want %v", recorder.Code, tt.expectedStatus)
			}

			// Check content type
			expectedContentType := "application/json"
			if contentType := recorder.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("UnauthorizedError() content-type = %v, want %v", contentType, expectedContentType)
			}

			// Check response body
			var actualBody WebError
			if err := json.Unmarshal(recorder.Body.Bytes(), &actualBody); err != nil {
				t.Errorf("UnauthorizedError() failed to unmarshal response: %v", err)
			}
			if actualBody != tt.expectedBody {
				t.Errorf("UnauthorizedError() body = %v, want %v", actualBody, tt.expectedBody)
			}

			// Check log output
			logOutput := buf.String()
			if !strings.Contains(logOutput, "Unauthorized access attempt:") {
				t.Errorf("UnauthorizedError() log output should contain 'Unauthorized access attempt:', got: %v", logOutput)
			}
		})
	}
}

func TestUserNotFound_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "user not found error message",
			want: "user not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := UserNotFound{}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
