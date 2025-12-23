package database

import (
	"test/internal/config"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error

	DB, err = gorm.Open(sqlite.Open("internal/database/" + config.Env.DBName + ".db"))
	if err != nil {
		log.Panic("Erro ao conectar ao banco de dados: ", err)
	}
}
