package main

import (
	"KeyForge/db"
	"KeyForge/web"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	dbLocation  = flag.String("db-location", "", "Path to boltDB database")
	httpAddress = flag.String("http-address", "127.0.0.1:5000", "Http port and host")
	configFile  = flag.String("config-file", "sharding.toml", "Config file for static sharding")
)

func ParseFlag() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatal("Must Provide DB Location")
	}
}

func main() {
	ParseFlag()

	configConents, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Read File (%q) : %w", *configFile, err)
	}

	DB, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(DB)

	http.HandleFunc("/set", srv.SetHandler)

	http.HandleFunc("/get", srv.GetHandler)

	log.Fatal(http.ListenAndServe(*httpAddress, nil))

}
