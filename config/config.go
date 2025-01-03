package config

import (
	"github.com/BurntSushi/toml"
)

type Shard struct {
	Name    string
	Idx     int
	Address string
}
type Config struct {
	Shard []Shard
}

func ParseFile(filename string) (Config, error) {
	var c Config
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return Config{}, err
	}
	return c, nil

}

type Shards struct {
	Count   int
	CurrIdx int
	Addrs   map[int]string
}
