package collaboration

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qascade/dcr/lib/collaboration/config"
)

func TestCollaboration(t *testing.T) {
	fmt.Println("Running TestCollaboration")
	testPath1, err := filepath.Abs("../../samples/test_collaboration")
	require.NoError(t, err)

	testPath2, err := filepath.Abs("../../samples/init_collaboration")
	require.NoError(t, err)
	testCompile(t, testPath1)
	testCompile(t, testPath2)
}

func testCompile(t *testing.T, path string) {
	fmt.Println("Running TestCompile")
	testConfig(t, path)
	//testAddressGraph(t, path)
}

func testAddressGraph(t *testing.T, pkgPath string) *Collaboration {
	fmt.Println("Running TestAddressGraph")
	collaboration, err := NewCollaboration(pkgPath)
	require.NoError(t, err)
	fmt.Println("AddressGraph =============")
	fmt.Println(*collaboration.AddressGraph)
	require.NotNil(t, collaboration.AddressGraph)
	return collaboration
}

func testConfig(t *testing.T, path string) *config.CollaborationConfig {
	fmt.Println("Running TestConfig")
	parser := config.ConfigParser(config.ConfigFolder{})
	_, err := parser.Parse(path)
	require.NoError(t, err)
	// for _, pkg := range pkgConfig.PackagesInfo {
	// 	fmt.Println("SourceSpec =============")
	// 	fmt.Println(*pkg.SourceSpec)
	// 	fmt.Println("DestinationSpec =============")
	// 	fmt.Println(*pkg.DestinationGroupSpec)
	// 	fmt.Println("TransformationSpec =============")
	// 	fmt.Println(*pkg.TransformationGroupSpec)
	// }

	pkgConfig, err := parser.Parse(path)
	require.NoError(t, err)
	for _, pkg := range pkgConfig.PackagesInfo {
		fmt.Println("SourceSpec =============")
		fmt.Println(*pkg.SourceSpec)
		fmt.Println("DestinationSpec =============")
		fmt.Println(*pkg.DestinationGroupSpec)
		fmt.Println("TransformationSpec =============")
		fmt.Println(*pkg.TransformationGroupSpec)
	}
	return pkgConfig

}
