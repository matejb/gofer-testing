package wittr_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matejb/gofer-testing/weather/with-complete-tests/wittr"
)

// što testirati:
// - da radi dobar upit na wittr servis
// - da podnosi grešku mreže
// - da podnosi grešku servisa
// - da prihvača dobar JSON odgovor
// - da podnosi loš JSON odgovor

func TestIsForShorts(t *testing.T) {
	cases := map[string]struct {
		simulateNetworkError bool
		city                 string
		responseStatusCode   int
		responseBody         string
		expectedRequestCity  string
		expectedTemperature  int
		expectedNetworkCall  bool
		expectedError        error
	}{
		"returns temperature for Zagreb provided by service": {
			city:                "Zagreb",
			responseStatusCode:  http.StatusOK,
			responseBody:        `{"current_condition": [{"FeelsLikeC": "40"}]}`,
			expectedTemperature: 40,
			expectedRequestCity: "Zagreb",
			expectedNetworkCall: true,
		},

		"returns temperature for Berlin provided by service": {
			city:                "Berlin",
			responseStatusCode:  http.StatusOK,
			responseBody:        `{"current_condition": [{"FeelsLikeC": "40"}]}`,
			expectedTemperature: 40,
			expectedRequestCity: "Berlin",
			expectedNetworkCall: true,
		},

		"returns error if service returned bad data for current temperature": {
			responseStatusCode:  http.StatusOK,
			responseBody:        `{"current_condition": [{"FeelsLikeC": "WAT"}]}`,
			expectedNetworkCall: true,
			expectedError:       errors.New("bad value for current temperature received"),
		},

		"returns error if service did not returned the current temperature": {
			responseStatusCode:  http.StatusOK,
			responseBody:        `{"current_condition": [{}]}`,
			expectedNetworkCall: true,
			expectedError:       errors.New("current temperature not received"),
		},

		"returns error when there is a network error": {
			simulateNetworkError: true,
			expectedError:        errors.New("connection refused"),
		},

		"returns error when server returns a non 2XX status code": {
			responseStatusCode:  http.StatusInternalServerError,
			expectedNetworkCall: true,
			expectedError:       errors.New("500 Internal Server Error"),
		},

		"returns error when server returns non JSON response": {
			responseStatusCode:  http.StatusOK,
			responseBody:        "NOT JSON",
			expectedNetworkCall: true,
			expectedError:       errors.New("decoding weather data failed"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			networkCallMade := false
			defer func() {
				if networkCallMade != tc.expectedNetworkCall {
					t.Errorf("expected network call to be %v got %v", tc.expectedNetworkCall, networkCallMade)
				}
			}()
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				networkCallMade = true

				if req.Method != http.MethodGet {
					t.Errorf("expected %q got %q", http.MethodGet, req.Method)
				}

				if req.URL.Path != "/"+tc.expectedRequestCity {
					t.Errorf("expected request path %q got %q", "/"+tc.expectedRequestCity, req.URL.Path)
				}

				if req.URL.Query().Get("format") != "j1" {
					t.Errorf("expected request query name %q to have value %q got %q", "format", "j1", req.URL.Query().Get("format"))
				}

				w.WriteHeader(tc.responseStatusCode)

				_, err := io.Copy(w, bytes.NewBuffer([]byte(tc.responseBody)))
				if err != nil {
					t.Fatalf("unexpected error in mock server body copy: %s", err)
				}
			}))
			defer ts.Close()

			wtr := wittr.New(tc.city)
			if wtr.BaseURL != "https://wttr.in" {
				t.Errorf("expected default BaseURL to be %q got %q", "https://wttr.in", wtr.BaseURL)
			}

			wtr.BaseURL = ts.URL

			if tc.simulateNetworkError {
				wtr.BaseURL = "http://localhost:0"
			}

			temp, err := wtr.CurrentTemperature()
			switch {
			case err != nil && tc.expectedError == nil:
				t.Fatalf("unexpected error: %s", err)
			case err == nil && tc.expectedError != nil:
				t.Fatal("expected error")
			case err != nil && tc.expectedError != nil && !strings.Contains(err.Error(), tc.expectedError.Error()):
				t.Fatalf("expected error %q but got %q", tc.expectedError, err)
			}

			if temp != tc.expectedTemperature {
				t.Errorf("expected %d got %d", tc.expectedTemperature, temp)
			}

		})
	}
}
