package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCamelCase(t *testing.T) {
	assert.Equal(t, "helloWorld", ToCamelCase("hello world"))
	assert.Equal(t, "helloWorld", ToCamelCase("Hello world"))
	assert.Equal(t, "helloWorld", ToCamelCase("HELLO WORLD"))
	assert.Equal(t, "helloWorld", ToCamelCase("helloWorld"))
	assert.Equal(t, "helloWorld", ToCamelCase("HelloWorld"))
}

func TestToSnakeCase(t *testing.T) {
	assert.Equal(t, "hello_world", ToSnakeCase("hello world"))
	assert.Equal(t, "hello_world", ToSnakeCase("Hello world"))
	assert.Equal(t, "hello_world", ToSnakeCase("HELLO WORLD"))
	assert.Equal(t, "hello_world", ToSnakeCase("helloWorld"))
	assert.Equal(t, "hello_world", ToSnakeCase("HelloWorld"))
}

func TestToPascalCase(t *testing.T) {
	assert.Equal(t, "HelloWorld", ToPascalCase("hello world"))
	assert.Equal(t, "HelloWorld", ToPascalCase("Hello world"))
	assert.Equal(t, "HelloWorld", ToPascalCase("HELLO WORLD"))
	assert.Equal(t, "HelloWorld", ToPascalCase("helloWorld"))
	assert.Equal(t, "HelloWorld", ToPascalCase("HelloWorld"))
}

func TestStripModulePacakge(t *testing.T) {
	assert.Equal(t, "User", StripModulePackage(`github.com/joeriddles/goalesce/pkg/model.User`, `github.com/joeriddles/goalesce`))
	assert.Equal(t, "User", StripModulePackage(`command-line-arguments.User`, `command-line-arguments`))

	assert.Equal(t, "gorm.DeletedAt", StripModulePackage("gorm.DeletedAt", `github.com/joeriddles/goalesce`))
	assert.Equal(t, "time.Time", StripModulePackage("time.Time", `github.com/joeriddles/goalesce`))
	assert.Equal(t, `github.come/some-user/some-package/pkg/model/User`, StripModulePackage(`github.come/some-user/some-package/pkg/model/User`, `github.com/joeriddles/goalesce`))
}
