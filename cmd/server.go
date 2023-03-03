package main

import (
	"aquiduct/config"
	"aquiduct/network"
	"flag"
)

var master *string = flag.String("host", "localhost", "Host address of cluster to join.")
var local *bool = flag.Bool("local", false, "Used to run the service locally without clustering (useful for testing).")

func main() {

	flag.Parse()

	config.Load()
	if local == nil || !*local {
		network.ClusterSetup(false)
		if master != nil {
			err := network.Join(*master)
			if err != nil {
				panic(err)
			}
		}
	} else {
		network.ClusterSetup(true)
	}

	network.Serve()

}
