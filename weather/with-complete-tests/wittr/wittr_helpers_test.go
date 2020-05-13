package wittr_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/matejb/gofer-testing/weather/with-complete-tests/wittr"
)

func setupFakeConnection(ws *wittr.WeatherSource) (cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, bytes.NewBuffer([]byte(
			`{"current_condition": [{"FeelsLikeC": "20"}]}`,
		)))
	}))

	ws.BaseURL = ts.URL

	return func() { ts.Close() }
}
