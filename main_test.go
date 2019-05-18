package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "good",
			in:             httptest.NewRequest("GET", "/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   message,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := newHandlers(nil)
			h.rootHandler(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Logf("Expected: %d\nGot: %d\n", test.expectedStatus, test.out.Code)
				t.Fail()
			}

			respBody := test.out.Body.String()
			if respBody != test.expectedBody {
				t.Logf("Expected: %s\nGot: %s\n", respBody, test.expectedBody)
				t.Fail()
			}
		})
	}
}
