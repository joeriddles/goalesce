package utils

import (
	"regexp"
	"strings"

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

// Convert the string to PascalCase.
func ToPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	camelCase := ToCamelCase(s)
	return strings.ToUpper(camelCase[0:1]) + camelCase[1:]
}
