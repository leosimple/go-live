package models

import (
	"errors"
	"go-live/orm"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Live struct {
	gorm.Model
	App      string `gorm:"not null"`
	Livename string `gorm:"not null"`
	Token    string `gorm:"not null"`
}

func init() {
	orm.Gorm.AutoMigrate(new(Live))
}

func CreateLive(live *Live) error {
	lives, err := GetAllLives()
	if err != nil {
		return err
	}

	for _, l := range lives {
		if l.App == live.App {
			if l.Livename == live.Livename {
				return errors.New("live name is exist.")
			}
		}
	}

	err = orm.Gorm.Create(live).Error
	if err != nil {
		return err
	}

	return nil
}

func GetAllLives() ([]Live, error) {
	var lives []Live
	err := orm.Gorm.Find(&lives).Error
	if err != nil {
		return nil, err
	}
	return lives, nil
}

func GetAllLivesByappname(appname string) ([]Live, error) {
	var lives []Live
	err := orm.Gorm.Where("app = ?", appname).Find(&lives).Error
	if err != nil {
		return nil, err
	}
	return lives, nil
}

func GetLiveByApporId(appname string, livename string) (*Live, error) {
	var live Live
	err := orm.Gorm.Where("app = ?", appname).Where("livename = ?", livename).First(&live).Error

	if err != nil {
		return nil, err
	}

	return &live, nil
}

func DeleteLive(live *Live) error {
	err := orm.Gorm.Delete(live).Error
	if err != nil {
		return err
	}

	return nil
}

func CheckToken(appname string, livename string, token string) bool {
	var lives []Live
	err := orm.Gorm.Where("app = ?", appname).Where("livename = ?", livename).Where("token = ?", token).Find(&lives).Error

	if err != nil {
		log.Println(err)
		return false
	}

	if len(lives) == 1 {
		return true
	}

	return false
}

func CheckLive(livename string) bool {
	var lives []Live
	err := orm.Gorm.Where("livename = ?", livename).Find(&lives).Error

	if err != nil {
		log.Println(err)
		return false
	}

	if len(lives) == 1 {
		return true
	}

	return false
}

func CheckLiveById(id int) bool {
	var lives []Live

	err := orm.Gorm.Where("id = ?", strconv.Itoa(id)).Find(&lives).Error

	if err != nil {
		log.Println(err)
		return false
	}

	if len(lives) == 1 {
		return true
	}

	return false
}
