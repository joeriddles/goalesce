package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Basic(t *testing.T) {
	err := run("../../examples/basic")
	require.NoError(t, err)
}

func Test_Cars(t *testing.T) {
	err := run("../../examples/cars")
	require.NoError(t, err)
}
