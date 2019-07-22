package restfulapi

import (
	"net/http"

	"go-live/av"
	"go-live/protocol/rtmp"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(stream av.Handler) *httprouter.Router {
	router := httprouter.New()

	// APP Restful API
	router.POST("/app/:appname/create", CreateAppHandler)
	router.GET("/app", ListAppsHandler)
	router.GET("/app/:appname/get", GetAppByNameHandler)
	router.DELETE("/app/:appname/del", DeleteAppByIdHandler)

	// Live Restful API
	router.POST("/live/:appname/:livename/create", CreateLiveHandler)
	router.GET("/live", ListLivesHandler)
	router.GET("/live/:appname", ListLivesByAppnameHandler)
	router.GET("/live/:appname/:livename/get", GetLiveByIdHandler)
	router.GET("/live/:appname/:livename/status", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		livestream := stream.(*rtmp.RtmpStream)
		GetLiveStatusHandler(w, r, ps, livestream)
	})
	router.PUT("/live/:appname/:livename/refershtoken", RefershLiveTokenByIdHandler)
	router.DELETE("/live/:appname/:livename/del", DeleteLiveByIdHandler)

	return router
}
