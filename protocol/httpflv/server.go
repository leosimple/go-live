package httpflv

import (
	"encoding/json"
	"go-live/av"
	"go-live/limiters"
	"go-live/models"
	"go-live/protocol/rtmp"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	handler av.Handler
	mux     *http.ServeMux
	limiter *limiters.ConnectionLimiter
}

type stream struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type streams struct {
	Publishers []stream `json:"publishers"`
	Players    []stream `json:"players"`
}

func NewServer(h av.Handler) *Server {
	maxconn, err := strconv.Atoi(os.Getenv("MAX_CONNECTION"))

	if err != nil {
		return nil
	}

	return &Server{
		handler: h,
		mux:     http.NewServeMux(),
		limiter: limiters.NewConnectionLimiter(maxconn),
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !server.limiter.GetConnection() {
		w.WriteHeader(http.StatusTooManyRequests)
		io.WriteString(w, "too many requests.")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	server.mux.ServeHTTP(w, r)

	defer server.limiter.FreeConnection()
}

func (server *Server) Serve(l net.Listener) error {
	server.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.handleConn(w, r)
	})
	server.mux.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
		server.getStream(w, r)
	})

	http.Serve(l, server)
	return nil
}

// 获取发布和播放器的信息
func (server *Server) getStreams(w http.ResponseWriter, r *http.Request) *streams {
	rtmpStream := server.handler.(*rtmp.RtmpStream)
	if rtmpStream == nil {
		return nil
	}
	msgs := new(streams)
	for item := range rtmpStream.GetStreams().IterBuffered() {
		if s, ok := item.Val.(*rtmp.Stream); ok {
			if s.GetReader() != nil {
				msg := stream{item.Key, s.GetReader().Info().UID}
				msgs.Publishers = append(msgs.Publishers, msg)
			}
		}
	}

	for item := range rtmpStream.GetStreams().IterBuffered() {
		ws := item.Val.(*rtmp.Stream).GetWs()
		for s := range ws.IterBuffered() {
			if pw, ok := s.Val.(*rtmp.PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					msg := stream{item.Key, pw.GetWriter().Info().UID}
					msgs.Players = append(msgs.Players, msg)
				}
			}
		}
	}

	return msgs
}

func (server *Server) getStream(w http.ResponseWriter, r *http.Request) {
	msgs := server.getStreams(w, r)
	if msgs == nil {
		return
	}
	resp, _ := json.Marshal(msgs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (server *Server) handleConn(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("http flv handleConn panic: ", r)
		}
	}()

	url := r.URL.String()
	u := r.URL.Path
	if pos := strings.LastIndex(u, "."); pos < 0 || u[pos:] != ".flv" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	path := strings.TrimSuffix(strings.TrimLeft(u, "/"), ".flv")
	paths := strings.SplitN(path, "/", 2)
	log.Println("url:", u, "path:", path, "paths:", paths)

	if len(paths) != 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// 判断视屏流是否发布,如果没有发布,直接返回404
	msgs := server.getStreams(w, r)
	if msgs == nil || len(msgs.Publishers) == 0 {
		http.Error(w, "invalid path", http.StatusNotFound)
		return
	} else {
		include := false
		for _, item := range msgs.Publishers {
			if item.Key == path {
				include = true
				break
			}
		}
		if include == false {
			http.Error(w, "invalid path", http.StatusNotFound)
			return
		}
	}

	appname := paths[0]
	splited := strings.Split(paths[1], "_")
	if len(splited) < 2 {
		http.Error(w, "param is too short.", http.StatusBadRequest)
		return
	}
	livename := splited[0]
	if !models.CheckLive(livename) {
		http.Error(w, "player livename is error", http.StatusBadRequest)
		return
	}

	token := splited[1]

	if !models.CheckToken(appname, livename, token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	writer := NewFLVWriter(paths[0], paths[1], url, w)

	server.handler.HandleWriter(writer)
	writer.Wait()
}
