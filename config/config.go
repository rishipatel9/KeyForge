package config

type Shard struct {
	Name    string
	Idx     int
	Address string
}
type Config struct {
	Shard []Shard
}
