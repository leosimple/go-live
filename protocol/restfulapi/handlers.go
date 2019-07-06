package restfulapi

import (
	"go-live/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// App Restful API
func CreateAppHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appname := ps.ByName("appname")
	liveon := r.FormValue("liveon")
	if liveon == "" {
		liveon = "on"
	}

	err := models.CreateApp(&models.App{
		Appname: appname,
		Liveon:  liveon,
	})

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, &Response{
		Code:    http.StatusOK,
		Message: "Successfully created this app.",
	})
}

func ListAppsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	apps, err := models.GetAllApps()

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, &AppsResponse{
		Code:    http.StatusOK,
		Data:    apps,
		Message: "Successfully acquired all applications.",
	})
}

func GetAppByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func UpdateAppByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func DeleteAppByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

// Live Restful API
func CreateLiveHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func ListLivesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func GetLiveByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func UpdateLiveByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func UpdateLiveTokenByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func DeleteLiveByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
