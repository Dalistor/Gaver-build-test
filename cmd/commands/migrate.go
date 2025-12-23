package commands

import (
	"log"
	"os"
	"strings"

	"test/internal/database"
	"test/internal/migrations"
)

func Migrate() {
	module, err := migrations.ReadGaverModule()
	if err != nil {
		log.Panic("Erro ao ler gaverModule.json: ", err)
	}

	if module.ProjectDatabaseType != "sqlite" {
		log.Panicf("Tipo de banco de dados incorreto. Esperado: sqlite, obtido: %s", module.ProjectDatabaseType)
	}

	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Println("Diretório de migrações não encontrado. Nada para migrar.")
		return
	}

	migrationFiles, err := migrations.ListMigrationFiles(migrationsDir, module.MigrationTag)
	if err != nil {
		log.Panic("Erro ao listar arquivos de migração: ", err)
	}

	if len(migrationFiles) == 0 {
		log.Println("Nenhuma migração pendente.")
		return
	}

	log.Printf("Encontradas %d migração(ões) pendente(s).", len(migrationFiles))

	for _, migrationFile := range migrationFiles {
		log.Printf("Executando migração: %s", migrationFile.FullName)

		sql, err := migrations.ReadMigrationFile(migrationFile.Path)
		if err != nil {
			log.Panicf("Erro ao ler arquivo de migração %s: %v", migrationFile.FullName, err)
		}

		statements := migrations.SplitSQLStatements(sql)
		
		for _, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" || statement == ";" {
				continue
			}

			if err := database.DB.Exec(statement).Error; err != nil {
				log.Panicf("Erro ao executar migração %s: %v\nSQL: %s", migrationFile.FullName, err, statement)
			}
		}

		if err := migrations.UpdateMigrationTag(migrationFile.Number); err != nil {
			log.Panicf("Erro ao atualizar migrationTag após migração %s: %v", migrationFile.FullName, err)
		}

		log.Printf("Migração %s executada com sucesso. MigrationTag atualizado para: %d", migrationFile.FullName, migrationFile.Number)
	}

	log.Println("Todas as migrações foram executadas com sucesso!")
}

