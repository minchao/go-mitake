package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/minchao/go-mitake"
)

func usage() {
	fmt.Println(`Usage: send [options]
Options are:
    -u  Mitake Username
    -p  Mitake Password
    -t  Destination phone number, for example: 0987654321
    -m  Message content`)
	os.Exit(0)
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
	flag.StringVar(&to, "t", "", "Destination phone number")
	flag.StringVar(&message, "m", "", "Message content")

	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
	}

	client := mitake.NewClient(username, password, nil)

	resp, err := client.Send(mitake.Message{
		Dstaddr: to,
		Smbody:  message,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "result: %+v\n", resp.INI)
}
