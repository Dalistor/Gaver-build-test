package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env = envConfig{}
var GaverSettings = gaverSettings{}

type gaverSettings struct {
	Type                string   `json:"type"`
	ProjectName         string   `json:"projectName"`
	ProjectVersion      string   `json:"projectVersion"`
	ProjectModules      []string `json:"projectModules"`
	ProjectDatabaseType string   `json:"projectDatabaseType"`
}

type envConfig struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string

	GinMode string
	GinJWT  string
	GinPort string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Erro ao carregar o arquivo .env: ", err)
	}

	Env = envConfig{
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		GinMode:    os.Getenv("GIN_MODE"),
		GinJWT:     os.Getenv("GIN_JWT"),
		GinPort:    os.Getenv("GIN_PORT"),
	}

	if err := loadGaverSettings(&GaverSettings); err != nil {
		log.Panic("Erro ao carregar as configurações do Gaver: ", err)
	}
}

func loadGaverSettings(settings *gaverSettings) error {
	jsonFile, err := os.Open("gaverModule.json")
	if err != nil {
		return fmt.Errorf("Erro ao carregar o arquivo gaverModule.json: %w", err)
	}
	defer jsonFile.Close()

	if err := json.NewDecoder(jsonFile).Decode(&settings); err != nil {
		return fmt.Errorf("Erro ao decodificar o arquivo gaverModule.json: %w", err)
	}

	return nil
}
