package web

import (
	"KeyForge/db"
	"fmt"
	"net/http"
)

type Server struct {
	db *db.Database
}

func NewServer(db *db.Database) *Server {
	return &Server{
		db: db,
	}
}

// Get Handler
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	ans, err := s.db.GetKey(key)

	fmt.Fprintf(w, "Value : %q, error : %v", ans, err)
	fmt.Print("Get called!!!")
}

// Set Handler
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err := s.db.SetKey(key, []byte(value))

	fmt.Fprintf(w, " error : %v", err)
	fmt.Print("Set called!!")
}
