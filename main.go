package main

import (
	"KeyForge/config"
	"KeyForge/db"
	"KeyForge/web"
	"flag"
	"log"
	"net/http"
)

var (
	dbLocation  = flag.String("db-location", "", "Path to boltDB database")
	httpAddress = flag.String("http-address", "127.0.0.1:5000", "Http port and host")
	configFile  = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard       = flag.String("shard", "", "Name of shard for data")
)

func main() {
	flag.Parse()

	c, err := config.ParseFile(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config %q: %v", *configFile, err)
	}

	shards, err := config.ParseShards(c.Shard, *shard)
	if err != nil {
		log.Fatalf("Error parsing shards config: %v", err)
	}

	log.Printf("Shard count is %d, current shard: %d", shards.Count, shards.CurrIdx)

	DB, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Error creating %q: %v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(DB, shards)

	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/purge", srv.DeleteExtraKeysHandler)

	log.Fatal(http.ListenAndServe(*httpAddress, nil))

}
