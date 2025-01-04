package config_test

import (
	"KeyForge/config"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func createConfig(t *testing.T, contents string) (config.Config, error) {
	t.Helper()

	f, err := ioutil.TempFile(os.TempDir(), "config.toml")

	if err != nil {
		log.Fatalf("Couldn't create temp file : %v", err)
		return config.Config{}, err
	}

	defer f.Close()
	defer os.Remove(f.Name())

	_, err = f.WriteString(contents)

	if err != nil {
		log.Fatalf("Could not write into the file : %v", err)
		return config.Config{}, err
	}

	c, err := config.ParseFile(f.Name())
	if err != nil {
		log.Fatalf("Could not parse the file : %v", err)
		return config.Config{}, err
	}
	return c, nil

}
func TestConfigParse(t *testing.T) {
	contents := `[[shard]]
		name="shard1"
		idx=0
		address="localhost:5000"
	`

	c, err := createConfig(t, contents)

	if err != nil {
		log.Fatalf("Error occured while creating config : %v", err)
	}

	want := config.Config{
		Shard: []config.Shard{
			{
				Name:    "shard1",
				Idx:     0,
				Address: "localhost:5000",
			},
		},
	}

	if !reflect.DeepEqual(c, want) {
		t.Errorf("The config does not match : %v and %v", c, want)
	}

}

func TestParseShard(t *testing.T) {
	c, err := createConfig(t, `
	[[shard]]
		name="shard1"
		idx=0
		address="localhost:5000"
	[[shard]]
		name="shard2"
		idx=1
		address="localhost:5001"
	`)

	if err != nil {
		t.Fatalf("Error occurred while creating config: %v", err)
	}

	// Update currShardName to "shard2"
	s, err := config.ParseShards(c.Shard, "shard2")
	if err != nil {
		t.Fatalf("Could not parse shards: %v", err)
	}

	want := &config.Shards{
		Count:   2,
		CurrIdx: 1,
		Addrs: map[int]string{
			0: "localhost:5000",
			1: "localhost:5001",
		},
	}

	if !reflect.DeepEqual(s, want) {
		t.Errorf("Mismatch in shard config. Got: %+v, Want: %+v", s, want)
	}
}
