package main

import (
	_ "go-live/conf"
	"go-live/protocol/httpflv"
	"go-live/protocol/httpopera"
	"go-live/protocol/restfulapi"
	"go-live/protocol/rtmp"
	"log"
	"net"
	"os"
	"time"
)

var (
	rtmpAddr    = os.Getenv("RTMP_ADDR")
	httpFlvAddr = os.Getenv("HTTPFLV_ADDR")
	operaAddr   = os.Getenv("MANAGE_ADDR")
	apiAddr     = os.Getenv("API_ADDR")
)

func startRtmp(stream *rtmp.RtmpStream) {
	rtmpListen, err := net.Listen("tcp", rtmpAddr)
	if err != nil {
		log.Fatal(err)
	}

	var rtmpServer *rtmp.Server

	rtmpServer = rtmp.NewRtmpServer(stream, nil)

	defer func() {
		if r := recover(); r != nil {
			log.Println("RTMP server panic: ", r)
		}
	}()
	log.Println("RTMP Listen On", rtmpAddr)
	rtmpServer.Serve(rtmpListen)
}

func startAPI() {
	apiListen, err := net.Listen("tcp", apiAddr)
	if err != nil {
		log.Fatal(err)
	}

	var apiServer *restfulapi.Server

	apiServer = restfulapi.NewServer()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("API server panic: ", r)
			}
		}()

		log.Println("API Listen On", apiAddr)
		apiServer.Serve(apiListen)
	}()
}

func startHTTPFlv(stream *rtmp.RtmpStream) {
	flvListen, err := net.Listen("tcp", httpFlvAddr)
	if err != nil {
		log.Fatal(err)
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("HTTP-FLV server panic: ", r)
			}
		}()
		log.Println("HTTP-FLV listen On", httpFlvAddr)
		hdlServer.Serve(flvListen)
	}()
}

func startHTTPOpera(stream *rtmp.RtmpStream) {
	if operaAddr != "" {
		opListen, err := net.Listen("tcp", operaAddr)
		if err != nil {
			log.Fatal(err)
		}
		opServer := httpopera.NewServer(stream, rtmpAddr)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("HTTP-Operation server panic: ", r)
				}
			}()
			log.Println("HTTP-Operation listen On", operaAddr)
			opServer.Serve(opListen)
		}()
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("go-live panic: ", r)
			time.Sleep(1 * time.Second)
		}
	}()

	stream := rtmp.NewRtmpStream()
	startHTTPFlv(stream)
	startHTTPOpera(stream)
	startAPI()
	startRtmp(stream)
}
