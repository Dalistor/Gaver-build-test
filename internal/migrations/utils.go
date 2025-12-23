package migrations

import (
	"encoding/json"
	"fmt"
	"os"
)

type GaverModule struct {
	Type                string   `json:"type"`
	ProjectName         string   `json:"projectName"`
	ProjectVersion      string   `json:"projectVersion"`
	ProjectModules      []string `json:"projectModules"`
	ProjectDatabaseType string   `json:"projectDatabaseType"`
	MigrationTag        int      `json:"migrationTag"`
}

func ReadGaverModule() (*GaverModule, error) {
	jsonFile, err := os.Open("gaverModule.json")
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir gaverModule.json: %w", err)
	}
	defer jsonFile.Close()

	var module GaverModule
	if err := json.NewDecoder(jsonFile).Decode(&module); err != nil {
		return nil, fmt.Errorf("erro ao decodificar gaverModule.json: %w", err)
	}

	return &module, nil
}

func WriteGaverModule(module *GaverModule) error {
	jsonFile, err := os.OpenFile("gaverModule.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir gaverModule.json para escrita: %w", err)
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(module); err != nil {
		return fmt.Errorf("erro ao codificar gaverModule.json: %w", err)
	}

	return nil
}

func UpdateMigrationTag(tag int) error {
	module, err := ReadGaverModule()
	if err != nil {
		return err
	}

	module.MigrationTag = tag
	return WriteGaverModule(module)
}



