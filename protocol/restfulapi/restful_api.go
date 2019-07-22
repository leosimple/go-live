package restfulapi

import (
	"go-live/av"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router *httprouter.Router
}

func NewServer(stream av.Handler) *Server {
	return &Server{
		router: NewRouter(stream),
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
