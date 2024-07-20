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
