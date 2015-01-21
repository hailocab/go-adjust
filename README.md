Go Adjust API Client
====================

This is a Go library for Adjust that currently supports server side event tracking.

For more information see https://docs.adjust.com/en/event-tracking/

Documentation
=============

The full documentation is available on [Godoc](http://godoc.org/github.com/hashicorp/consul/api)

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
