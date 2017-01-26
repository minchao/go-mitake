# go-mitake

go-mitake is a Go client library for accessing the Mitake SMS API.

## Usage

```go
import "github.com/minchao/go-mitake"
```

Construct a new Mitake SMS client, then use to access the Mitake API. For example:

```go
client := mitake.NewClient("USERNAME", "PASSWORD", nil)

// Retrieving your account balance
balance, err := client.QueryAccountPoint()

// Send an SMS
message := mitake.Message{
    Dstaddr: "0987654321",
    Smbody:  "Test SMS",
}

response, err := client.Send(message)

```

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).