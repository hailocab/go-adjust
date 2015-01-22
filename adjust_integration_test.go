// +build integration

package goadjust

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	AppToken     = os.Getenv("ADJUST_APP_TOKEN")
	DeviceID     = os.Getenv("ADJUST_DEVICE_ID")
	EventToken   = os.Getenv("ADJUST_EVENT_TOKEN")
	RevenueToken = os.Getenv("ADJUST_REVENUE_TOKEN")
)

func TestIntegrationTrackEvent(t *testing.T) {
	adjust := New(AppToken, Sandbox)

	_, err := adjust.TrackEvent(IDFA, DeviceID, EventToken, time.Now())
	require.Nil(t, err)
}

func TestIntegrationTrackRevenue(t *testing.T) {
	adjust := New(AppToken, Sandbox)

	_, err := adjust.TrackRevenue(IDFA, DeviceID, EventToken, 1000, time.Now())
	require.Nil(t, err)

}

func TestIntegrationCustomAttributes(t *testing.T) {
	adjust := New(AppToken, Sandbox)

	_, err := adjust.TrackEventWithParams(IDFA, DeviceID, EventToken, time.Now(), map[string]string{
		"testdata": "1234",
	})
	require.Nil(t, err)

}
