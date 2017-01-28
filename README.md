# go-mitake

[![GoDoc](https://godoc.org/github.com/minchao/go-mitake?status.svg)](https://godoc.org/github.com/minchao/go-mitake)
[![Build Status](https://travis-ci.org/minchao/go-mitake.svg?branch=master)](https://travis-ci.org/minchao/go-mitake)
[![Go Report Card](https://goreportcard.com/badge/github.com/minchao/go-mitake)](https://goreportcard.com/report/github.com/minchao/go-mitake)

go-mitake is a Go client library for accessing the Mitake SMS API (Taiwan mobile phone number only).

## Installation

```bash
go get github.com/minchao/go-mitake
```

## Usage

```go
import "github.com/minchao/go-mitake"
```

Construct a new Mitake SMS client, then use to access the Mitake API. For example:

```go
client := mitake.NewClient("USERNAME", "PASSWORD", nil)

// Retrieving your account balance
balance, err := client.QueryAccountPoint()
```

Send an SMS:

```go
message := mitake.Message{
    Dstaddr: "0987654321",
    Smbody:  "Test SMS",
}

response, err := client.Send(message)
```

Send multiple SMS:

```go
messages := []mitake.Message{
    {
        Dstaddr: "0987654321",
        Smbody:  "Test SMS",
    },
    // ...
}

response, err := client.SendBatch(messages)
```

Query the status of messages:

```go
response, err := err := client.QueryMessageStatus([]string{"MESSAGE_ID1", "MESSAGE_ID2"})
```

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).