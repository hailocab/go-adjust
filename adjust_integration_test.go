// +build integration

package adjust

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

}

func TestIntegrationCustomAttributes(t *testing.T) {

}
