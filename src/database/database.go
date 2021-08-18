package database

import (
	"log"

	"github.com/datti-to/purrmannplus-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() {

	type Open func(string) gorm.Dialector
	var o Open

	switch config.DATABASE_TYPE {
	case "POSTGRES":
		o = postgres.Open
	case "MYSQL":
		o = mysql.Open
	case "SQLITE":
		o = sqlite.Open
	default:
		log.Fatalf("DATABASE_TYPE env has to one of ('POSTGRES', 'MYSQL', 'SQLITE')")
	}

	var err error
	DB, err = gorm.Open(o(config.DATABASE_URI), &gorm.Config{})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

}

func CloseDB() {
	dialect, err := DB.DB()
	if err != nil {
		log.Fatalf("models.CloseDB err: %v", err)
	}
	defer dialect.Close()
}
