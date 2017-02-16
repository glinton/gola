// build and run this in order to serve up static js api client page that will connect to it

package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/glinton/gola/api"
)

var (
	port = "8080"
)

// if a port was passed, use it
func init() {
	flag.StringVar(&port, "port", port, "Port to lisen on")
	flag.Parse()

	if portnum, err := strconv.Atoi(port); portnum < 0 || err != nil {
		port = "8080"
	}
}

// serve the worlds
func main() {
	err := api.Start(fmt.Sprintf("0.0.0.0:" + port))
	if err != nil {
		fmt.Printf("Failed to start server - %s\n", err)
	}
}
