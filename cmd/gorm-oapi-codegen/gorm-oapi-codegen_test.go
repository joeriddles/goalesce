package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Basic(t *testing.T) {
	err := run(
		"../../examples/basic",
		"./generated/basic",
		"github.com/joeriddles/gorm-oapi-codegen",
		"github.com/joeriddles/gorm-oapi-codegen/examples/basic",
	)
	require.NoError(t, err)
}

func Test_Cars(t *testing.T) {
	err := run(
		"../../examples/cars",
		"./generated/cars",
		"github.com/joeriddles/gorm-oapi-codegen",
		"github.com/joeriddles/gorm-oapi-codegen/examples/cars",
	)
	require.NoError(t, err)
}
