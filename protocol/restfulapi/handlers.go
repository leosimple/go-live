package restfulapi

import (
	"errors"
	"fmt"
	"go-live/functions"
	"go-live/models"
	"go-live/orm"
	"go-live/protocol/rtmp"
	"net/http"
	"strings"

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

func GetAppByNameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appname := ps.ByName("appname")
	if appname == "" {
		SendErrorResponse(w, http.StatusBadRequest, "Appname is not be null.")
		return
	}

	app, err := models.GetAppByName(appname)

	if err != nil {
		SendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	SendResponse(w, http.StatusOK, &AppResponse{
		Code:    http.StatusOK,
		Data:    app,
		Message: "Successfully obtained the corresponding application.",
	})
}

func GetLiveStatusHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, livestream *rtmp.RtmpStream) {
	msgs := new(Streams)
	for item := range livestream.GetStreams().IterBuffered() {
		if s, ok := item.Val.(*rtmp.Stream); ok {
			if s.GetReader() != nil {
				switch s.GetReader().(type) {
				case *rtmp.VirReader:
					v := s.GetReader().(*rtmp.VirReader)
					msg := Stream{item.Key, v.Info().URL, v.ReadBWInfo.StreamId, v.ReadBWInfo.VideoDatainBytes, v.ReadBWInfo.VideoSpeedInBytesperMS,
						v.ReadBWInfo.AudioDatainBytes, v.ReadBWInfo.AudioSpeedInBytesperMS}
					msgs.Publishers = append(msgs.Publishers, msg)
				}
			}
		}
	}

	for item := range livestream.GetStreams().IterBuffered() {
		ws := item.Val.(*rtmp.Stream).GetWs()
		for s := range ws.IterBuffered() {
			if pw, ok := s.Val.(*rtmp.PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					switch pw.GetWriter().(type) {
					case *rtmp.VirWriter:
						v := pw.GetWriter().(*rtmp.VirWriter)
						msg := Stream{item.Key, v.Info().URL, v.WriteBWInfo.StreamId, v.WriteBWInfo.VideoDatainBytes, v.WriteBWInfo.VideoSpeedInBytesperMS,
							v.WriteBWInfo.AudioDatainBytes, v.WriteBWInfo.AudioSpeedInBytesperMS}
						msgs.Players = append(msgs.Players, msg)
					}
				}
			}
		}
	}

	publisheronline := false

	for _, v := range msgs.Publishers {
		data := fmt.Sprintf("%s/%s", ps.ByName("appname"), ps.ByName("livename"))

		if strings.Split(v.Key, "_")[0] == data {
			publisheronline = true
		}
	}

	playercount := len(msgs.Players)

	SendResponse(w, http.StatusOK, &StatusResponse{
		Code:        http.StatusOK,
		Message:     "Player Status Get OK.",
		PlayerCount: playercount,
		IsPublisher: publisheronline,
	})
}

func DeleteAppByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appname := ps.ByName("appname")
	if appname == "" {
		SendErrorResponse(w, http.StatusBadRequest, "Appname is not be null.")
		return
	}

	if !models.CheckAppByName(appname) {
		SendErrorResponse(w, http.StatusNotFound, "This app is not in the database.")
		return
	}

	app, err := models.GetAppByName(appname)

	if err != nil {
		SendErrorResponse(w, http.StatusNotFound, err.Error())
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
	var live models.Live
	appname := ps.ByName("appname")
	livename := ps.ByName("livename")

	err := orm.Gorm.Where("app = ?", appname).Where("livename = ?", livename).First(&live).Error
	if err != nil {
		SendErrorResponse(w, http.StatusNotFound, errors.New("lives cannot find.").Error())
		return
	}

	SendResponse(w, http.StatusOK, &LiveResponse{
		Code:    http.StatusOK,
		Message: "Successfully obtained the corresponding live.",
		Data:    live,
	})
}

func RefershLiveTokenByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var live models.Live
	appname := ps.ByName("appname")
	livename := ps.ByName("livename")
	token := functions.RandomString(6)

	err := orm.Gorm.Where("app = ?", appname).Where("livename = ?", livename).First(&live).Error

	if err != nil {
		SendErrorResponse(w, http.StatusNotFound, errors.New("lives cannot find.").Error())
		return
	}

	err = orm.Gorm.Model(&live).Update("Token", token).Error

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
	livename := ps.ByName("livename")

	if !models.CheckLive(livename) {
		SendErrorResponse(w, http.StatusNotFound, "This live not in database.")
		return
	}

	live, err := models.GetLiveByApporId(appname, livename)

	if err != nil {
		SendErrorResponse(w, http.StatusNotFound, "404 live not found")
		return
	}

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
