package restfulapi

import (
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router *httprouter.Router
}

func NewServer() *Server {
	return &Server{
		router: NewRouter(),
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	server.router.ServeHTTP(w, r)
}

func (server *Server) Serve(l net.Listener) error {
	http.Serve(l, server)
	return nil
}
