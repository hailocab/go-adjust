Go Adjust API Client
====================

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg "GoDoc")](http://godoc.org/github.com/hailocab/go-adjust) 
[![Build Status](https://img.shields.io/travis/hailocab/g-adjust/master.svg "Build Status")](https://travis-ci.org/hailocab/go-adjust) 
[![Civerage](http://gocover.io/_badge/github.com/hailocab/go-adjust "Coverage")](http://gocover.io/github.com/hailocab/go-adjust)

This is a Go library for Adjust that currently supports server side event tracking.

For more information see https://docs.adjust.com/en/event-tracking/

Documentation
=============

The full documentation is available on [Godoc](http://godoc.org/github.com/hailocab/go-adjust)

## Usage

```go
client := adjust.New(goadjust.Config{
    AppToken: "4w565xzmb54d",
    Environment: goadjust.Sandbox,
})

_, err := client.TrackEvent("event_token", goadjust.IDFA, "D2CADB5F-410F-4963-AC0C-2A78534BDF1E", time.Now())
if err != nil {
    log.Fatal(err)
}

_, err = client.TrackRevenue("event_token", 990, goadjust.IDFA, "D2CADB5F-410F-4963-AC0C-2A78534BDF1E", time.Now())
if err != nil {
    log.Fatal(err)
}
```

## Tests

Go-adjust contains both unit tests and integration tests. To run the tests run the commands below:

```
go test ./...
```

To run the integration tests you must pass the tokens and a valid device ID (IDFA) via environment variables along with the `integration` build tag. When running the integration tests events are published to the specified Adjust account in Sandbox mode.

```
ADJUST_APP_TOKEN=... ADJUST_EVENT_TOKEN=... ADJUST_DEVICE_ID=... go test -tags=integration
```

*Note* Please note that the tests currently do not work in Go 1.3 as the mocked `net/http.Transport` is being ignored (possibly due to a bug in the stdlib?), this should hopefully be fixed soon.
