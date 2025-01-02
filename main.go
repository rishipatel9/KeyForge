package main

import (
	"KeyForge/config"
	"KeyForge/db"
	"KeyForge/web"
	"flag"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
)

var (
	dbLocation  = flag.String("db-location", "", "Path to boltDB database")
	httpAddress = flag.String("http-address", "127.0.0.1:5000", "Http port and host")
	configFile  = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard       = flag.String("shard", "", "Name of shard for data")
)

func ParseFlag() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatal("Must Provide DB Location")
	}
	if *shard == "" {
		log.Fatal("Must Provide Shard")
	}
}

func main() {
	ParseFlag()

	var c config.Config

	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("toml.DecodeFile(%q) : %v ", *configFile, err)
	}

	var shardCount int
	var shardIdx int = -1
	var shardAddrs = make(map[int]string)
	for _, s := range c.Shard {
		shardAddrs[s.Idx] = s.Address
		if s.Idx+1 > shardCount {
			shardCount = s.Idx + 1
		}
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatal("Shard not found")
	}

	log.Printf("Shard Count is %d, current Shard %d", shardCount, shardIdx)
	DB, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(DB, shardIdx, shardCount, shardAddrs)

	http.HandleFunc("/set", srv.SetHandler)

	http.HandleFunc("/get", srv.GetHandler)

	log.Fatal(http.ListenAndServe(*httpAddress, nil))

}
