package db_test

import (
	"KeyForge/db"
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetSet(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "testdb")
	if err != nil {
		t.Fatalf("could not create temp file : %v", err)
	}
	name := f.Name()
	f.Close()

	defer os.Remove(name)

	db, closeFunc, err := db.NewDatabase(name)

	if err != nil {
		t.Fatalf("could not create database :%v", err)
	}
	defer closeFunc()
	if err := db.SetKey("part", []byte("great")); err != nil {
		t.Fatalf("count not write key to database : %v", err)
	}
	value, err := db.GetKey("party")

	if err != nil {
		t.Fatalf("could not get the key :%v", err)
	}
	if bytes.Equal(value, []byte("great")) {
		t.Fatalf("unexpected value got key party : got  %q and got %q ", value, "great")
	}

}
