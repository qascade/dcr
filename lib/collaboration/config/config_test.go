package config

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	parser := ConfigParser(ConfigFolder{})
	testPath, err := filepath.Abs("../../../samples/test_collaboration")
	require.NoError(t, err)
	pkgConfig, err := parser.Parse(testPath)
	require.NoError(t, err)
	for _, pkg := range pkgConfig.PackagesInfo {
		fmt.Println("SourceSpec =============")
		fmt.Println(*pkg.SourceSpec)
		fmt.Println("DestinationSpec =============")
		fmt.Println(*pkg.DestinationGroupSpec)
		fmt.Println("TransformationSpec =============")
		fmt.Println(*pkg.TransformationGroupSpec)
	}
}
