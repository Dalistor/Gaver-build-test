package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"test/internal/migrations"
)

func MakeMigrations() {
	module, err := migrations.ReadGaverModule()
	if err != nil {
		log.Panic("Erro ao ler gaverModule.json: ", err)
	}

	if module.ProjectDatabaseType != "sqlite" {
		log.Panicf("Tipo de banco de dados incorreto. Esperado: sqlite, obtido: %s", module.ProjectDatabaseType)
	}

	models, err := migrations.ScanModelsFromModules()
	if err != nil {
		log.Panic("Erro ao escanear models: ", err)
	}

	if len(models) == 0 {
		log.Println("Nenhum model encontrado em modules/*/models. Nada para migrar.")
		return
	}

	nextMigrationTag := module.MigrationTag + 1
	migrationFileName := fmt.Sprintf("%04d_%s.sql", nextMigrationTag, generateMigrationName(models))

	migrationsDir := "migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		log.Panic("Erro ao criar diretório de migrações: ", err)
	}

	migrationPath := filepath.Join(migrationsDir, migrationFileName)

	sql := migrations.GenerateSQLForSQLite(models)

	if err := os.WriteFile(migrationPath, []byte(sql), 0644); err != nil {
		log.Panic("Erro ao escrever arquivo de migração: ", err)
	}

	module.MigrationTag = nextMigrationTag
	if err := migrations.WriteGaverModule(module); err != nil {
		log.Panic("Erro ao atualizar gaverModule.json: ", err)
	}

	log.Printf("Migração criada: %s", migrationFileName)
	log.Printf("MigrationTag atualizado para: %d", nextMigrationTag)
}

func generateMigrationName(models []migrations.ModelInfo) string {
	if len(models) == 0 {
		return "empty"
	}

	if len(models) == 1 {
		return strings.ToLower(models[0].TableName)
	}

	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("multiple_%s", timestamp)
}

