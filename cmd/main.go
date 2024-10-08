package main

import (
	"github/pbreedt/testcontainers/nats/pub"
	"github/pbreedt/testcontainers/nats/sub"
	"os"
	"strconv"

	"github.com/nats-io/nats.go"
)

/*
Usage: (requires a running NATS server)

	args:
		1. pub or sub
		2. number of seconds to run

example:

	$ go run cmd/main.go sub 10 &
	$ go run cmd/main.go pub 5
*/
func main() {
	svr := os.Args[1]
	secs, _ := strconv.Atoi(os.Args[2])

	if svr == "pub" {
		pub.Publish(nats.DefaultURL, secs)
	} else {
		sub.Subscribe(nats.DefaultURL, secs)
	}
}
