package main

import (
	"log"
)

func main() {
	// load the config
	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}
	registry, err := initServer()
	if err != nil {
		log.Fatal(err)
	}
	// run the server
	err = registry.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
