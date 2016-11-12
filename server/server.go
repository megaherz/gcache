package server

import (
	"net/http"
	"log"
	"gcache"
)

type Server struct {
	cache *gcache.Cache
}

func (s *Server) Run(addr string) {
	http.HandleFunc("/keys", &keysHandler{cache: s.cache}.Handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func NewServer() *Server {
	return &Server{
		cache: gcache.NewCache(),
	}
}