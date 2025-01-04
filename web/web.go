package web

import (
	"KeyForge/db"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
)

type Server struct {
	db         *db.Database
	shardIdx   int
	shardCount int
	shardAddrs map[int]string
}

func NewServer(db *db.Database, shardIdx int, shardCount int, addrs map[int]string) *Server {
	return &Server{
		db:         db,
		shardIdx:   shardIdx,
		shardCount: shardCount,
		shardAddrs: addrs,
	}
}

func (s *Server) GetShard(key string) int {
	h := fnv.New64a()
	h.Write([]byte(key))
	fmt.Println(((h.Sum64()) % uint64(s.shardCount)))
	return int((h.Sum64()) % uint64(s.shardCount))
}

func (s *Server) Redirect(w http.ResponseWriter, r *http.Request, shard int) {
	url := "http://" + s.shardAddrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d to shard %d (%q)\n", s.shardIdx, shard, url)
	res, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error redirecting the request : %v", err)
		return
	}
	defer res.Body.Close()
	io.Copy(w, res.Body)
}

// Get Handler
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shard := s.GetShard(key)

	if shard != s.shardIdx {
		s.Redirect(w, r, shard)
	}

	ans, err := s.db.GetKey(key)
	fmt.Fprintf(w, "Value : %q, error : %v, shard : %d , addrs : %q, current shardIdx : %d", ans, err, shard, s.shardAddrs[shard], s.shardIdx)
}

// Set Handler
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shard := s.GetShard(key)

	if shard != s.shardIdx {
		s.Redirect(w, r, shard)
	}
	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "error: %v, shard : %d, shardIdx : %d", err, shard, s.shardIdx)
}
