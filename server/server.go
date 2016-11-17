package server

import (
	"net/http"
	"log"
	"gcache"
	"gcache/server/handlers"
	"fmt"
)

const headerAuthorization  = "Authorization"

type Server struct {
	cache *gcache.Cache
	urlLoggingEnabled bool
	pws string
}

func (s *Server) Run(addr string) {

	keysHandler := new(handlers.KeysHandler).Init(s.cache)
	listsHandler := new(handlers.ListsHandler).Init(s.cache)
	hashesHandler := new(handlers.HashesHandler).Init(s.cache)

	s.middleware("/keys", keysHandler)
	s.middleware("/lists", listsHandler)
	s.middleware("/hashes", hashesHandler)


	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s *Server) middleware(route string, handler handlers.Handler){
	http.HandleFunc(route,
		s.urlLoggingHandler(s.authHandler(handler)).ServeHTTP)
}

func NewServer() *Server {
	return &Server{
		cache: gcache.NewCache(),
	}
}

func NewServerWithAuth(pws string) *Server {
	return &Server{
		cache: gcache.NewCache(),
		pws: pws,
	}
}

func (s *Server) SetUrlLogging(enabled bool)  {
	s.urlLoggingEnabled = true
}

// Authentication handler
func (s *Server) authHandler(h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if (s.pws != "") {

			psw := r.Header.Get(headerAuthorization)

			if (psw != s.pws) {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

// Url logging handler. Logs urls if url logging is enabled
func (s *Server) urlLoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if s.urlLoggingEnabled {
			fmt.Println(*r.URL)
		}

		h.ServeHTTP(w, r)
	})
}