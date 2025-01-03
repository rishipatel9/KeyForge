package config_test

import (
	"KeyForge/config"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestConfigParse(t *testing.T) {
	contents := `[[shard]]
	name="shard1"
	idx=0
	address="localhost:5000"
	`

	f, err := ioutil.TempFile(os.TempDir(), "config.toml")

	if err != nil {
		log.Fatalf("Couldn't create temp file : %v", err)
	}

	defer f.Close()
	defer os.Remove(f.Name())

	_, err = f.WriteString(contents)

	if err != nil {
		log.Fatalf("Could not write into the file : %v", err)
	}

	c, err := config.ParseFile(f.Name())
	if err != nil {
		log.Fatalf("Could not parse the file : %v", err)
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

	if reflect.DeepEqual(c, want) {
		t.Errorf("The config does not match : %v and %v", c, want)
	}

}
