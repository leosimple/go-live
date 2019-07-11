package orm

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Gorm *gorm.DB

func init() {
	var err error

	Gorm, err = gorm.Open("mysql", os.Getenv("MYSQL_DSN"))

	if err != nil {
		log.Fatal(err)
	}
}
