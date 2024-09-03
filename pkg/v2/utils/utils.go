package utils

import (
	"errors"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/joeriddles/goalesce/pkg/v2/entity"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var camelCaseRegex = regexp.MustCompile("([a-z0-9])([A-Z])")

// Convert the string to camelCase.
func ToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	s = camelCaseRegex.ReplaceAllString(s, "${1} ${2}")

	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	for i := 1; i < len(words); i++ {
		caser := cases.Title(language.AmericanEnglish)
		words[i] = caser.String(strings.ToLower(words[i]))
	}

	words[0] = strings.ToLower(words[0])
	return strings.Join(words, "")

}

// Convert the string to snake_case.
func ToSnakeCase(s string) string {
	camelCase := ToCamelCase(s)
	snake := camelCaseRegex.ReplaceAllString(camelCase, "${1}_${2}")
	return strings.ToLower(snake)
}

// Convert the string to html-case.
func ToHtmlCase(s string) string {
	snakeCase := ToSnakeCase(s)
	return strings.ReplaceAll(snakeCase, "_", "-")
}

// Convert the string to PascalCase.
func ToPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	camelCase := ToCamelCase(s)
	return strings.ToUpper(camelCase[0:1]) + camelCase[1:]
}

func WrapID(model *entity.GormModelMetadata) string {
	result := "id"

	idField, err := First(model.AllFields(), func(f *entity.GormModelField) bool {
		return f.Name == "ID"
	})
	if err == nil {
		if basicType, ok := idField.GetType().(*types.Basic); ok {
			wrapper := basicType.Name()
			result = fmt.Sprintf("%v(id)", wrapper)
		}
	}

	return result
}

// FindGoMod searches upwards from the given path for a go.mod file
func FindGoMod(startPath string, module string) (string, error) {
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(absPath, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			goModBytes, err := os.ReadFile(goModPath)
			if err != nil {
				return "", err
			}

			goModFile, err := modfile.Parse(goModPath, goModBytes, nil)
			if err == nil {
				if goModFile.Module.Mod.Path == module {
					return goModPath, nil
				}
				// If it's not the right module, keep looking
			}
		}

		// Move up one directory level
		parentPath := filepath.Dir(absPath)
		if parentPath == absPath {
			// We have reached the root directory
			break
		}
		absPath = parentPath
	}

	return "", fmt.Errorf("go.mod file not found for mod: %v", module)
}

func First[S ~[]E, E any](s S, f func(E) bool) (E, error) {
	index := slices.IndexFunc(s, f)
	if index != -1 {
		return s[index], nil
	}
	var empty E
	return empty, errors.New("no matching object found in slice")
}

func Map[S ~[]E, E any, R any](s S, f func(E) R) []R {
	result := []R{}
	for _, e := range s {
		r := f(e)
		result = append(result, r)
	}
	return result
}

func MapPointers[S ~[]*R, R any](s S) []R {
	return Map(s, func(val *R) R { return *val })
}

func StripModulePackage(s, moduleName string) string {
	// moduleName = strings.ReplaceAll(moduleName, "/", `\/`)
	pattern := fmt.Sprintf(`%v([/A-z-]+)?\.`, moduleName)
	re := regexp.MustCompile(pattern)
	s = re.ReplaceAllString(s, "")
	return s
}

// Is the type a non-simple type?
func IsComplexType(typ string) bool {
	if typ == "" {
		return false
	}
	return strings.HasPrefix(typ, "*") || !strings.HasPrefix(typ, "[]") || typ[0:1] != strings.ToUpper(typ[0:1])
}

// Is the type a simple, built-in type?
func IsSimpleType(t string) bool {
	return !IsComplexType(t)
}

var goalesceTagPattern *regexp.Regexp = regexp.MustCompile(`goalesce:"(.*?)"`)

// Parse goalesce settings from a field's tag
func ParseGoalesceTagSettings(tag string) (map[string]string, error) {
	settings := map[string]string{}

	if !goalesceTagPattern.MatchString(tag) {
		return settings, nil
	}

	matches := goalesceTagPattern.FindStringSubmatch(tag)
	gormSettings := strings.Split(matches[1], ";")
	for _, kvp := range gormSettings {
		kvp = strings.TrimSpace(kvp)
		if kvp == "" {
			continue
		}

		keyAndValue := strings.SplitN(kvp, ":", 2)
		if len(keyAndValue) != 2 {
			return nil, fmt.Errorf("cannot parse goalesce settings key-value pair: %v", kvp)
		}
		settings[keyAndValue[0]] = keyAndValue[1]
	}

	return settings, nil
}
