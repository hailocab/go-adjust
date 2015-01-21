package adjust

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestNew(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(adjust, `{"status":"OK"}`, func() {
		assert.Equal(t, "token", adjust.AppToken)
		assert.Equal(t, Sandbox, adjust.Environment)
	})
}

func TestTrackEvent(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(adjust, `{"status":"OK"}`, func() {
		resp, err := adjust.TrackEvent(IDFA, "1234", "5678", time.Unix(12345678, 0))
		require.Nil(t, err)
		assert.Equal(t, "OK", resp.Status)
	})
}

func TestTrackEventInvalidToken(t *testing.T) {
	adjust := New("invalidtoken", Sandbox)
	mockHTTPClient(adjust, `Event failed (Event Token for wrong app: '5qvlia') (app_token: za9taw6qel8j, adid: a594e6bdfef2a0112c9b53ae76c1f20d)`, func() {
		_, err := adjust.TrackEvent(IDFA, "1234", "5678", time.Unix(12345678, 0))
		require.NotNil(t, err)
	})
}

func TestTrackEventUnknownDevice(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(adjust, `Event failed (Device not found, contact support@adjust.com) (app_token: foo, adid: bar)`, func() {
		_, err := adjust.TrackEvent(IDFA, "1234", "5678", time.Unix(12345678, 0))
		require.NotNil(t, err)
		assert.Equal(t, ErrDeviceNotFound, err)
	})
}

func TestTrackRevenue(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(adjust, `{"status":"OK"}`, func() {
		resp, err := adjust.TrackRevenue(IDFA, "1234", "5678", 200, time.Unix(12345678, 0))
		require.Nil(t, err)
		assert.Equal(t, "OK", resp.Status)
	})
}

func TestTrackRevenueInvalidAmount(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(adjust, `Event failed (Negative revenue: -1) (app_token: za9taw6qel8d, adid: a594e6bdfef2a0112c9b53ae76c1f20d)`, func() {
		_, err := adjust.TrackRevenue(IDFA, "1234", "5678", -1, time.Unix(12345678, 0))
		require.NotNil(t, err)
	})
}

func mockHTTPClient(client *Client, response string, f func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
	client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(ts.URL)
			},
		},
	}
	APIURL = ts.URL

	defer func() {
		ts.Close()
		client.HTTPClient = nil
		APIURL = "https://s2s.adjust.com"
	}()

	f()
}
