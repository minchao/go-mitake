package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/minchao/go-mitake"
)

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, `Usage: send [options]`)
	flag.PrintDefaults()
}

func main() {
	var (
		username string
		password string
		to       string
		message  string
	)

	flag.StringVar(&username, "u", os.Getenv("MITAKE_USERNAME"), "Username")
	flag.StringVar(&password, "p", os.Getenv("MITAKE_PASSWORD"), "Password")
	flag.StringVar(&to, "t", "", "Destination phone number, for example: 0987654321")
	flag.StringVar(&message, "m", "", "Message content")

	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 3 {
		flag.Usage()
		os.Exit(1)
	}

	client := mitake.NewClient(username, password, nil)

	resp, err := client.Send(mitake.Message{
		Dstaddr: to,
		Smbody:  message,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, "result: %+v\n", resp.INI)
}
