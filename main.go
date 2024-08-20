package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/mdns"
)

func main() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)

	entries := []*mdns.ServiceEntry{}

	go func() {
		for entry := range entriesCh {
			entries = append(entries, entry)
		}
	}()

	// QueryParams
	params := mdns.QueryParam{
		Service: "_services._dns-sd._udp",
		Timeout: time.Second * 5,
		Entries: entriesCh,
	}

	// Start the lookup
	mdns.Query(&params)
	close(entriesCh)

	debug(entries)
}

func debug(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(b))
}
