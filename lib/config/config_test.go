package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	parser := ConfigParser(ConfigFolder{})
	cSpec, pkgConfig, err := parser.Parse("../../samples/test_collaboration")
	require.NoError(t, err)
	fmt.Println(cSpec)
	fmt.Println(pkgConfig)
}
