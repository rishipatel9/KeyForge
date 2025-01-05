package web_test

import (
	"KeyForge/config"
	"KeyForge/db"
	"KeyForge/web"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func createShardDB(t *testing.T, idx int) *db.Database {
	t.Helper()
	tempFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("db%d", idx))

	if err != nil {
		t.Fatalf("error occured while creating temp file: %v", err)
	}
	tempFile.Close()

	name := tempFile.Name()

	t.Cleanup(func() { os.Remove(name) })

	db, closefunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("error occured while creating database : %v", db)
	}
	t.Cleanup(func() { closefunc() })

	return db
}
func createShardServer(t *testing.T, idx int, addrs map[int]string) (*db.Database, *web.Server) {
	t.Helper()

	db := createShardDB(t, idx)

	cfg := &config.Shards{
		Addrs:   addrs,
		Count:   len(addrs),
		CurrIdx: idx,
	}

	s := web.NewServer(db, cfg)
	return db, s
}

func TestServer(t *testing.T) {
	var ts1SetHandler, ts1GetHandler func(w http.ResponseWriter, r *http.Request)
	var ts2SetHandler, ts2GetHandler func(w http.ResponseWriter, r *http.Request)
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			ts1GetHandler(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			ts1SetHandler(w, r)
		}
	}))
	defer ts1.Close()
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			ts2GetHandler(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			ts2SetHandler(w, r)
		}
	}))

	defer ts2.Close()

	addrs := map[int]string{
		0: strings.TrimPrefix(ts1.URL, "http://"),
		1: strings.TrimPrefix(ts2.URL, "http://"),
	}

	db1, web1 := createShardServer(t, 0, addrs)
	db2, web2 := createShardServer(t, 1, addrs)

	// Calculated manually and depends on the sharding function.
	keys := map[string]int{
		"Soviet": 1,
		"USA":    0,
	}

	ts1GetHandler = web1.GetHandler
	ts2GetHandler = web2.GetHandler
	ts2SetHandler = web2.SetHandler
	ts1SetHandler = web1.SetHandler

	for key := range keys {
		_, err := http.Get(fmt.Sprintf("%s/set?key=%s&value=value-%s", ts1.URL, key, key))
		if err != nil {
			t.Fatalf("could not set key: %v", err)
		}
	}

	for key := range keys {
		resp, err := http.Get(fmt.Sprintf("%s/get?key=%s", ts1.URL, key))
		if err != nil {
			t.Fatalf("could not get key: %v", err)
		}
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}
		want := []byte("value-" + key)
		if !bytes.Contains(contents, want) {
			t.Errorf("Unexpected contents of the key %q: got %q, want %q", key, contents, want)
		}
	}

	value1, err := db1.GetKey("USA")
	if err != nil {
		t.Fatalf("USA key error: %v", err)
	}

	want1 := "value-USA"
	if !bytes.Equal(value1, []byte(want1)) {
		t.Errorf("Unexpected value of USA key: got %q, want %q", value1, want1)
	}

	value2, err := db2.GetKey("Soviet")
	if err != nil {
		t.Fatalf("Soviet key error: %v", err)
	}

	want2 := "value-Soviet"
	if !bytes.Equal(value2, []byte(want2)) {
		t.Errorf("Unexpected value of Soviet key: got %q, want %q", value2, want2)
	}
}
