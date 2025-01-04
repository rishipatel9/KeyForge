package config

import (
	"fmt"
	"hash/fnv"

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

func ParseShards(shards []Shard, currShardName string) (*Shards, error) {
	shardCount := len(shards)
	shardIdx := -1
	addrs := make(map[int]string)

	for _, s := range shards {
		if _, ok := addrs[s.Idx]; ok {
			return nil, fmt.Errorf("duplicate shard Index : %d", s.Idx)
		}

		addrs[s.Idx] = s.Address
		if s.Name == currShardName {
			shardIdx = s.Idx
		}
	}

	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("shard with index %d not found", i)
		}
	}

	if shardIdx < 0 {
		return nil, fmt.Errorf("shard with name %q was not found ", currShardName)
	}

	return &Shards{
		Addrs:   addrs,
		Count:   shardCount,
		CurrIdx: shardIdx,
	}, nil

}
func (s *Shards) Index(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}
