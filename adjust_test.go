package goadjust

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	adjust := New("token", Sandbox)
	assert.Equal(t, "token", adjust.AppToken)
	assert.Equal(t, Sandbox, adjust.Environment)
}

func TestTrackEvent(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(t, adjust, url.Values{
		"app_token":   {"token"},
		"event_token": {"5678"},
		"idfa":        {"1234"},
		"created_at":  {"1970-05-23T22:21:18Z+0100"},
		"environment": {"sandbox"},
		"params":      {"e30="},
		"s2s":         {"1"},
	}, 200, `{"status":"OK"}`, func() {
		resp, err := adjust.TrackEvent(IDFA, "1234", "5678", time.Unix(12345678, 0))
		require.Nil(t, err)
		assert.Equal(t, "OK", resp.Status)
	})
}

func TestTrackEventWithParams(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(t, adjust, url.Values{
		"app_token":   {"token"},
		"event_token": {"5678"},
		"idfa":        {"1234"},
		"created_at":  {"1970-05-23T22:21:18Z+0100"},
		"environment": {"sandbox"},
		"params":      {"eyJ0ZXN0ZGF0YSI6IjQzMjEifQ=="},
		"s2s":         {"1"},
	}, 200, `{"status":"OK"}`, func() {
		resp, err := adjust.TrackEventWithParams(IDFA, "1234", "5678", time.Unix(12345678, 0), map[string]string{
			"testdata": "4321",
		})
		require.Nil(t, err)
		assert.Equal(t, "OK", resp.Status)
	})
}

func TestTrackEventInvalidToken(t *testing.T) {
	adjust := New("invalidtoken", Sandbox)
	mockHTTPClient(t, adjust, url.Values{
		"app_token":   {"invalidtoken"},
		"event_token": {"5678"},
		"idfa":        {"1234"},
		"created_at":  {"1970-05-23T22:21:18Z+0100"},
		"environment": {"sandbox"},
		"params":      {"e30="},
		"s2s":         {"1"},
	}, 500, `Event failed (Event Token for wrong app: '5qvlia') (app_token: za9taw6qel8j, idfa: 1234)`, func() {
		_, err := adjust.TrackEvent(IDFA, "1234", "5678", time.Unix(12345678, 0))
		require.NotNil(t, err)
	})
}

func TestTrackEventUnknownDevice(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(t, adjust, url.Values{
		"app_token":   {"token"},
		"event_token": {"5678"},
		"idfa":        {"unknowndevice"},
		"created_at":  {"1970-05-23T22:21:18Z+0100"},
		"environment": {"sandbox"},
		"params":      {"e30="},
		"s2s":         {"1"},
	}, 500, `Event failed (Device not found, contact support@adjust.com) (app_token: token, idfa: 1234)`, func() {
		_, err := adjust.TrackEvent(IDFA, "unknowndevice", "5678", time.Unix(12345678, 0))
		require.NotNil(t, err)
		assert.Equal(t, ErrDeviceNotFound, err)
	})
}

func TestTrackRevenue(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(t, adjust, url.Values{
		"app_token":   {"token"},
		"event_token": {"5678"},
		"idfa":        {"1234"},
		"amount":      {"200"},
		"created_at":  {"1970-05-23T22:21:18Z+0100"},
		"environment": {"sandbox"},
		"params":      {"e30="},
		"s2s":         {"1"},
	}, 200, `{"status":"OK"}`, func() {
		resp, err := adjust.TrackRevenue(IDFA, "1234", "5678", 200, time.Unix(12345678, 0))
		require.Nil(t, err)
		assert.Equal(t, "OK", resp.Status)
	})
}

func TestTrackRevenueWithParams(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(t, adjust, url.Values{
		"app_token":   {"token"},
		"event_token": {"5678"},
		"idfa":        {"1234"},
		"amount":      {"200"},
		"created_at":  {"1970-05-23T22:21:18Z+0100"},
		"environment": {"sandbox"},
		"params":      {"eyJ0ZXN0ZGF0YSI6IjQzMjEifQ=="},
		"s2s":         {"1"},
	}, 200, `{"status":"OK"}`, func() {
		resp, err := adjust.TrackRevenueWithParams(IDFA, "1234", "5678", 200, time.Unix(12345678, 0), map[string]string{
			"testdata": "4321",
		})
		require.Nil(t, err)
		assert.Equal(t, "OK", resp.Status)
	})
}

func TestTrackRevenueInvalidAmount(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(t, adjust, url.Values{
		"app_token":   {"token"},
		"event_token": {"5678"},
		"idfa":        {"1234"},
		"amount":      {"-1"},
		"created_at":  {"1970-05-23T22:21:18Z+0100"},
		"environment": {"sandbox"},
		"params":      {"e30="},
		"s2s":         {"1"},
	}, 500, `Event failed (Negative revenue: -1) (app_token: za9taw6qel8d, idfa: 1234)`, func() {
		_, err := adjust.TrackRevenue(IDFA, "1234", "5678", -1, time.Unix(12345678, 0))
		require.NotNil(t, err)
	})
}

func TestInvalidJSONResponse(t *testing.T) {
	adjust := New("token", Sandbox)
	mockHTTPClient(t, adjust, url.Values{
		"app_token":   {"token"},
		"event_token": {"5678"},
		"idfa":        {"1234"},
		"amount":      {"1000"},
		"created_at":  {"1970-05-23T22:21:18Z+0100"},
		"environment": {"sandbox"},
		"params":      {"e30="},
		"s2s":         {"1"},
	}, 500, `?!?!?!`, func() {
		_, err := adjust.TrackRevenue(IDFA, "1234", "5678", 1000, time.Unix(12345678, 0))
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})
}

func TestHTTPClientError(t *testing.T) {
	adjust := New("token", Sandbox)
	adjust.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return nil, fmt.Errorf("client error")
			},
		},
	}

	_, err := adjust.TrackRevenue(IDFA, "1234", "5678", -1, time.Unix(12345678, 0))
	require.NotNil(t, err)
	assert.Equal(t, err.Error(), "Post https://s2s.adjust.com/event: client error")
}

func mockHTTPClient(t *testing.T, client *Client, expectedReq url.Values, code int, response string, f func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, expectedReq, r.PostForm)

		w.WriteHeader(code)
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
