package restfulapi

import (
	"errors"
	"fmt"
	"go-live/functions"
	"go-live/models"
	"go-live/orm"
	"net/http"
	"strconv"

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

	SendResponse(w, http.StatusOK, &Response{
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

	SendResponse(w, http.StatusOK, &AppsResponse{
		Code:    http.StatusOK,
		Data:    apps,
		Message: "Successfully acquired all applications.",
	})
}

func GetAppByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appid := ps.ByName("appid")
	if appid == "" {
		SendErrorResponse(w, http.StatusBadRequest, "Appid is not be null.")
		return
	}

	id, err := strconv.Atoi(appid)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	app, err := models.GetAppById(id)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, http.StatusOK, &AppResponse{
		Code:    http.StatusOK,
		Data:    app,
		Message: "Successfully obtained the corresponding application.",
	})
}

func DeleteAppByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appid := ps.ByName("appid")
	if appid == "" {
		SendErrorResponse(w, http.StatusBadRequest, "Appid is not be null.")
		return
	}

	id, err := strconv.Atoi(appid)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !models.CheckAppById(id) {
		SendErrorResponse(w, http.StatusInternalServerError, "This app is not in the database.")
		return
	}

	app, err := models.GetAppById(id)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = models.DeleteApp(app)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = models.DeleteLive(&models.Live{App: app.Appname})

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, http.StatusOK, &Response{
		Code:    http.StatusOK,
		Message: "Successfully deleted this app.",
	})
}

// Live Restful API
func CreateLiveHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appname := ps.ByName("appname")
	livename := ps.ByName("livename")

	token := functions.RandomString(6)

	err := models.CreateLive(&models.Live{
		App:      appname,
		Livename: livename,
		Token:    token,
	})

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, http.StatusOK, &LiveTokenResponse{
		Code:    http.StatusOK,
		Message: "Successfully created this live.",
		Token:   token,
	})
}

func ListLivesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lives, err := models.GetAllLives()
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	SendResponse(w, http.StatusOK, &LivesResponse{
		Code:    http.StatusOK,
		Message: "Successfully acquired all lives.",
		Data:    lives,
	})
}

func ListLivesByAppnameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appname := ps.ByName("appname")

	lives, err := models.GetAllLivesByappname(appname)
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	SendResponse(w, http.StatusOK, &LivesResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Successfully acquired all lives : %s.", appname),
		Data:    lives,
	})
}

func GetLiveByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var lives []models.Live
	appname := ps.ByName("appname")
	liveid := ps.ByName("liveid")

	err := orm.Gorm.Where("app = ?", appname).Where("id = ?", liveid).Find(&lives).Error
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(lives) == 0 {
		SendErrorResponse(w, http.StatusBadRequest, errors.New("lives cannot find.").Error())
		return
	}

	SendResponse(w, http.StatusOK, &LiveResponse{
		Code:    http.StatusOK,
		Message: "Successfully obtained the corresponding live.",
		Data:    lives[0],
	})
}

func RefershLiveTokenByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var lives []models.Live
	appname := ps.ByName("appname")
	liveid := ps.ByName("liveid")
	token := functions.RandomString(6)

	err := orm.Gorm.Where("app = ?", appname).Where("id = ?", liveid).Find(&lives).Error

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(lives) == 0 {
		SendErrorResponse(w, http.StatusBadRequest, errors.New("lives cannot find.").Error())
		return
	}

	err = orm.Gorm.Model(&lives[0]).Update("Token", token).Error

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, http.StatusOK, &LiveTokenResponse{
		Code:    http.StatusOK,
		Message: "Successfully refreshed Token.",
		Token:   token,
	})
}

func DeleteLiveByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appname := ps.ByName("appname")
	liveid := ps.ByName("liveid")

	id, err := strconv.Atoi(liveid)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !models.CheckLiveById(id) {
		SendErrorResponse(w, http.StatusBadRequest, "This live not in database.")
		return
	}

	live, err := models.GetLiveByApporId(appname, id)

	err = models.DeleteLive(live)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, http.StatusOK, &Response{
		Code:    http.StatusOK,
		Message: "Successfully deleted this live.",
	})
}
