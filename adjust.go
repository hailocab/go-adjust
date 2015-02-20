package goadjust

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05Z-0700"
)

var (
	// ErrDeviceNotFound is returned when adjust did not recognise the device
	ErrDeviceNotFound = errors.New("Device not found")

	// APIURL is the URL for the S2S API
	APIURL = "https://s2s.adjust.com"
)

// Environment represents where the data should be stored in Adjust (defaults
// to "production")
type Environment string

const (
	// Production tells Adjust to store the event as a production event
	Production Environment = "production"
	// Sandbox tells Adjust to store the event in the sandbox
	Sandbox = "sandbox"
)

// DeviceIDType represents which environment the data should be stored for.
type DeviceIDType string

const (
	// IDFA is the iOS ID for Advertisers
	IDFA DeviceIDType = "idfa"
	// IDFV is the iOS ID for Vendors
	IDFV = "idfv"
	// Mac is the MAC address of device (without “:”) (Android only)
	Mac = "mac"
	// MacMD5 is the MD5 of MAC (upper case without “:”) (Android only)
	MacMD5 = "mac_md5"
	// MacSHA1 is the SHA1 of MAC (upper case with “:”) (Android only)
	MacSHA1 = "mac_sha1"
	// AndroidID is an Android ID
	AndroidID = "android_id"
	// GPSAdID is the Google Play Advertiser ID
	GPSAdID = "gps_adid"
)

// A Client is used to send track requests to Adjust.
type Client struct {
	AppToken    string
	Environment Environment

	HTTPClient *http.Client
}

// New creates a new Adjust Client.
func New(appToken string, env Environment) *Client {
	return &Client{
		AppToken:    appToken,
		Environment: env,
	}
}

// Response represents the response from an track request.
type Response struct {
	Status       string `json:"status"`
	TrackerToken string `json:"tracker_token"`
	TrackerName  string `json:"tracker_name"`
	Network      string `json:"network"`
	Country      string `json:"country"`
}

// TrackEvent tracks a non-revenue event
func (c *Client) TrackEvent(deviceIDType DeviceIDType, deviceID string, eventToken string, t time.Time) (resp *Response, err error) {
	return c.send("/event", url.Values{
		"event_token":        {eventToken},
		"created_at":         {t.Format(timeLayout)},
		string(deviceIDType): {deviceID},
	}, map[string]string{})
}

// TrackEventWithParams tracks a non-revenue event with custom parameters. These
// parameters can be used in callbacks.
func (c *Client) TrackEventWithParams(deviceIDType DeviceIDType, deviceID string, eventToken string, t time.Time, params map[string]string) (resp *Response, err error) {
	return c.send("/event", url.Values{
		"event_token":        {eventToken},
		"created_at":         {t.Format(timeLayout)},
		string(deviceIDType): {deviceID},
	}, params)
}

// TrackRevenue tracks a revenue event
func (c *Client) TrackRevenue(deviceIDType DeviceIDType, deviceID string, eventToken string, amount int, t time.Time) (resp *Response, err error) {
	return c.send("/revenue", url.Values{
		"event_token":        {eventToken},
		"amount":             {strconv.FormatInt(int64(amount), 10)},
		"created_at":         {t.Format(timeLayout)},
		string(deviceIDType): {deviceID},
	}, map[string]string{})
}

// TrackRevenueWithParams tracks a revenue event with custom parameters. These
// parameters can be used in callbacks.
func (c *Client) TrackRevenueWithParams(deviceIDType DeviceIDType, deviceID string, eventToken string, amount int64, t time.Time, params map[string]string) (resp *Response, err error) {
	return c.send("/revenue", url.Values{
		"event_token":        {eventToken},
		"amount":             {strconv.FormatInt(int64(amount), 10)},
		"created_at":         {t.Format(timeLayout)},
		string(deviceIDType): {deviceID},
	}, params)
}

func (c *Client) send(path string, req url.Values, params map[string]string) (resp *Response, err error) {
	// Add client fields to request
	req.Add("s2s", "1")
	req.Add("app_token", c.AppToken)
	req.Add("params", base64EncodeMap(params))
	if c.Environment != "" {
		req.Add("environment", string(c.Environment))
	}

	// Send request
	var httpClient = c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	httpResp, err := httpClient.PostForm(APIURL+path, req)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, httpResp.Body)
	httpResp.Body.Close()

	switch {
	case strings.HasPrefix(buf.String(), "Event failed (Device not found, contact support@adjust.com)"):
		return nil, ErrDeviceNotFound
	case strings.HasPrefix(buf.String(), "Event failed"):
		msg := strings.TrimSpace(buf.String())
		return nil, errors.New(msg)
	}

	// Unmarshal response
	if err := json.NewDecoder(buf).Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func base64EncodeMap(m map[string]string) string {
	b, _ := json.Marshal(m)

	return base64.StdEncoding.EncodeToString(b)
}
