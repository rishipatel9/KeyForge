package db_test

import (
	"KeyForge/db"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestDeleteExtraKeys(t *testing.T) {
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

	setKey(t, db, "party", "great")
	setKey(t, db, "usa", "capitalistPigs")
	// value := getKey(t, db, "party")

	// if !bytes.Equal(value, []byte("great")) {
	// 	t.Fatalf("unexpected value got key party : got  %q and want %q ", value, "great")
	// }

	if err := db.DeleteExtraKeys(func(name string) bool {
		return name == "usa"
	}); err != nil {
		t.Fatalf("could not delete extra keys : %v", err)
	}

	if value := getKey(t, db, "party"); !bytes.Equal(value, []byte("great")) {
		t.Fatalf("unexpected value got key party : got  %q and want %q ", value, "great")
	}
	if value := getKey(t, db, "us"); !bytes.Equal(value, []byte("")) {
		t.Fatalf("unexpected value got key usa : got  %q and want %q ", value, "")
	}

}

func setKey(t *testing.T, d *db.Database, key, value string) {
	t.Helper()

	if err := d.SetKey(key, []byte(value)); err != nil {
		t.Fatalf("Set key failed : %v", err)
	}
	fmt.Printf("set the key %q to value %q successfully\n", key, value)
}

func getKey(t *testing.T, d *db.Database, key string) []byte {
	t.Helper()
	value, err := d.GetKey(key)
	if err != nil {
		t.Fatalf("Get key failed : %v", err)
		return []byte("")
	}
	return value

}
