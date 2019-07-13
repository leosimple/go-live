package models

import (
	"go-live/orm"
	"log"

	"github.com/jinzhu/gorm"
)

type App struct {
	gorm.Model
	Appname string `gorm:"not null;unique"`
	Liveon  string `gorm:"not null"`
}

func init() {
	orm.Gorm.AutoMigrate(new(App))
}

func CreateApp(app *App) error {
	err := orm.Gorm.Create(app).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAppById(id int) (*App, error) {
	var app App

	err := orm.Gorm.First(&app, id).Error

	return &app, err
}

func GetAppByName(appname string) (*App, error) {
	var app App

	err := orm.Gorm.Where("appname = ?", appname).First(&app).Error

	if err != nil {
		return nil, err
	}

	return &app, nil
}

func GetAllApps() ([]App, error) {
	var apps []App
	err := orm.Gorm.Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func DeleteApp(app *App) error {
	err := orm.Gorm.Delete(app).Error
	if err != nil {
		return err
	}

	return nil
}

func GetAppsByNameorLiveon(appname string) ([]App, error) {
	var apps []App

	err := orm.Gorm.Where("appname = ?", appname).Where("liveon = ?", "on").Find(&apps).Error

	if err != nil {
		return nil, err
	}

	return apps, nil
}

func CheckAppByName(name string) bool {
	count := 0
	err := orm.Gorm.Where("appname = ?", name).Count(&count).Error

	if err != nil {
		log.Println(err)
		return false
	}

	if count == 1 {
		return true
	}

	return false
}
