package restfulapi

import "github.com/julienschmidt/httprouter"

func NewRouter() *httprouter.Router {
	router := httprouter.New()

	// APP Restful API
	router.POST("/app/:appname/create", CreateAppHandler)
	router.GET("/app", ListAppsHandler)
	router.GET("/app/:appid/get", GetAppByIdHandler)
	router.DELETE("/app/:appid/del", DeleteAppByIdHandler)

	// Live Restful API
	router.POST("/live/:appname/:livename/create", CreateLiveHandler)
	router.GET("/live", ListLivesHandler)
	router.GET("/live/:appname", ListLivesByAppnameHandler)
	router.GET("/live/:appname/:liveid/get", GetLiveByIdHandler)
	router.PUT("/live/:appname/:liveid/:token/refershtoken", RefershLiveTokenByIdHandler)
	router.DELETE("/live/:appname/:liveid/:token/del", DeleteLiveByIdHandler)

	return router
}
