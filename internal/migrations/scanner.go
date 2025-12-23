package migrations

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type ModelInfo struct {
	Name      string
	TableName string
	Fields    []FieldInfo
}

type FieldInfo struct {
	Name         string
	Type         string
	SQLType      string
	IsPrimaryKey bool
	IsNotNull    bool
	IsUnique     bool
	IsIndex      bool
	DefaultValue string
	Tags         map[string]string
}

func ScanModels(directory string) ([]ModelInfo, error) {
	var models []ModelInfo

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fileModels, err := parseGoFile(path)
		if err != nil {
			return fmt.Errorf("erro ao parsear arquivo %s: %w", path, err)
		}

		models = append(models, fileModels...)
		return nil
	})

	return models, err
}

func parseGoFile(filePath string) ([]ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var models []ModelInfo

	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		modelInfo := ModelInfo{
			Name:      ts.Name.Name,
			TableName: getTableName(ts.Name.Name),
			Fields:    []FieldInfo{},
		}

		for _, field := range st.Fields.List {
			if len(field.Names) == 0 {
				continue
			}

			fieldInfo := parseField(field)
			if fieldInfo != nil {
				modelInfo.Fields = append(modelInfo.Fields, *fieldInfo)
			}
		}

		if len(modelInfo.Fields) > 0 {
			models = append(models, modelInfo)
		}

		return true
	})

	return models, nil
}

func parseField(field *ast.Field) *FieldInfo {
	if len(field.Names) == 0 {
		return nil
	}

	fieldInfo := &FieldInfo{
		Name:  field.Names[0].Name,
		Tags:  make(map[string]string),
		Type:  getTypeName(field.Type),
		IsNotNull: true,
	}

	if field.Tag != nil {
		tagValue := strings.Trim(field.Tag.Value, "`")
		tags := parseStructTag(tagValue)
		fieldInfo.Tags = tags

		if gormTag, ok := tags["gorm"]; ok {
			parseGormTag(gormTag, fieldInfo)
		}
	}

	return fieldInfo
}

func parseStructTag(tag string) map[string]string {
	tags := make(map[string]string)
	parts := strings.Split(tag, " ")

	for _, part := range parts {
		if idx := strings.Index(part, ":"); idx > 0 {
			key := part[:idx]
			value := strings.Trim(part[idx+1:], "\"")
			tags[key] = value
		}
	}

	return tags
}

func parseGormTag(tag string, fieldInfo *FieldInfo) {
	parts := strings.Split(tag, ";")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		if part == "primary_key" || part == "primaryKey" {
			fieldInfo.IsPrimaryKey = true
		}
		if part == "not null" || part == "notnull" {
			fieldInfo.IsNotNull = true
		}
		if strings.HasPrefix(part, "unique") {
			fieldInfo.IsUnique = true
		}
		if part == "index" {
			fieldInfo.IsIndex = true
		}
		if strings.HasPrefix(part, "type:") {
			fieldInfo.SQLType = strings.TrimPrefix(part, "type:")
		}
		if strings.HasPrefix(part, "default:") {
			fieldInfo.DefaultValue = strings.TrimPrefix(part, "default:")
		}
	}
}

func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", getTypeName(t.X), t.Sel.Name)
	case *ast.ArrayType:
		return "[]" + getTypeName(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", getTypeName(t.Key), getTypeName(t.Value))
	default:
		return "unknown"
	}
}

func getTableName(structName string) string {
	return strings.ToLower(structName) + "s"
}

// FindModelDirectories percorre a pasta modules e encontra todos os diretórios models
// dentro de cada módulo. A pasta modules será criada pelo framework.
func FindModelDirectories() ([]string, error) {
	var dirs []string
	modulesDir := "modules"

	// Verificar se modules existe (será criada pelo framework)
	// Se não existir, retorna lista vazia sem erro
	if _, err := os.Stat(modulesDir); os.IsNotExist(err) {
		return dirs, nil // Retorna vazio, não é erro - pasta será criada pelo framework
	}

	// Percorrer cada subdiretório em modules
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Verificar se existe models/ dentro deste módulo
		modelPath := filepath.Join(modulesDir, entry.Name(), "models")
		if info, err := os.Stat(modelPath); err == nil && info.IsDir() {
			dirs = append(dirs, modelPath)
		}
	}

	return dirs, nil
}

// ScanModelsFromModules escaneia todos os diretórios models encontrados em modules
// e retorna uma lista consolidada de todos os models encontrados
func ScanModelsFromModules() ([]ModelInfo, error) {
	var allModels []ModelInfo

	modelDirs, err := FindModelDirectories()
	if err != nil {
		return nil, err
	}

	if len(modelDirs) == 0 {
		return allModels, nil
	}

	for _, dir := range modelDirs {
		models, err := ScanModels(dir)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear %s: %w", dir, err)
		}
		allModels = append(allModels, models...)
	}

	return allModels, nil
}

