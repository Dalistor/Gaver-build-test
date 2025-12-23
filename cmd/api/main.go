package main

import (
	"fmt"
	"log"
	"os"

	"test/cmd/commands"
	"test/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	initialization()

	// Rodar comandos
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			commands.Migrate()
		case "makemigrations":
			commands.MakeMigrations()
		default:
			log.Panicf("Comando inválido: %s", os.Args[1])
		}
	}
}

func initialization() {
	showLogoCLI()

	ginEngine := gin.Default()

	switch config.Env.GinMode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		log.Panicf("Modo de execução inválido: %s", config.Env.GinMode)
	}

	log.Println("Modo de execução: ", config.Env.GinMode)

	ginEngine.Use(config.Cors())

	ginEngine.Run(":" + config.Env.GinPort)
}

func showLogoCLI() {
	const orange = "\033[38;5;208m"
	const reset = "\033[0m"

	fmt.Printf("%s================================================%s\n", orange, reset)
	fmt.Printf("%s#   Gaver: %s - %s                    %s\n", orange, config.GaverSettings.ProjectName, config.GaverSettings.ProjectVersion, reset)
	fmt.Printf("%s================================================%s\n", orange, reset)
}
