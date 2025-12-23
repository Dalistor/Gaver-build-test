package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MigrationFile struct {
	Path     string
	Number   int
	Name     string
	FullName string
}

func ListMigrationFiles(migrationsDir string, currentTag int) ([]MigrationFile, error) {
	var files []MigrationFile

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, fmt.Errorf("erro ao ler diretório de migrações: %w", err)
	}

	pattern := regexp.MustCompile(`^(\d{4})_(.+)\.sql$`)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		matches := pattern.FindStringSubmatch(entry.Name())
		if len(matches) != 3 {
			continue
		}

		migrationNumber, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		if migrationNumber <= currentTag {
			continue
		}

		files = append(files, MigrationFile{
			Path:     filepath.Join(migrationsDir, entry.Name()),
			Number:   migrationNumber,
			Name:     matches[2],
			FullName: entry.Name(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Number < files[j].Number
	})

	return files, nil
}

func ReadMigrationFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo de migração: %w", err)
	}

	sql := string(content)
	sql = strings.TrimSpace(sql)

	return sql, nil
}

func SplitSQLStatements(sql string) []string {
	statements := []string{}
	lines := strings.Split(sql, "\n")
	var currentStatement strings.Builder
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		
		currentStatement.WriteString(line)
		currentStatement.WriteString(" ")
		
		if strings.HasSuffix(line, ";") {
			stmt := strings.TrimSpace(currentStatement.String())
			if stmt != "" && stmt != ";" {
				statements = append(statements, stmt)
			}
			currentStatement.Reset()
		}
	}
	
	if currentStatement.Len() > 0 {
		stmt := strings.TrimSpace(currentStatement.String())
		if stmt != "" {
			if !strings.HasSuffix(stmt, ";") {
				stmt += ";"
			}
			statements = append(statements, stmt)
		}
	}

	return statements
}

