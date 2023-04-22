package collaboration

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/collaboration/config"
)

func TestCollaboration(t *testing.T) {
	collabConfig := testConfig(t)
	graph := testAddressGraph(t, collabConfig)
	testAuthorization(t, graph)
}

func testAuthorization(t *testing.T, graph *address.Graph) {
}

func testAddressGraph(t *testing.T, collabConfig *config.CollaborationConfig) *address.Graph {
	collaboration, err := NewCollaboration(*collabConfig)
	require.NoError(t, err)
	fmt.Println("AddressGraph =============")
	fmt.Println(*collaboration.AddressGraph)
	require.NotNil(t, collaboration.AddressGraph)
	return collaboration.AddressGraph
}

func testConfig(t *testing.T) *config.CollaborationConfig {
	fmt.Println("Running TestConfig")
	parser := config.ConfigParser(config.ConfigFolder{})
	testPath1, err := filepath.Abs("../../samples/test_collaboration")
	require.NoError(t, err)

	testPath2, err := filepath.Abs("../../samples/init_collaboration")
	require.NoError(t, err)

	_, err = parser.Parse(testPath1)
	require.NoError(t, err)
	// for _, pkg := range pkgConfig.PackagesInfo {
	// 	fmt.Println("SourceSpec =============")
	// 	fmt.Println(*pkg.SourceSpec)
	// 	fmt.Println("DestinationSpec =============")
	// 	fmt.Println(*pkg.DestinationGroupSpec)
	// 	fmt.Println("TransformationSpec =============")
	// 	fmt.Println(*pkg.TransformationGroupSpec)
	// }

	pkgConfig, err := parser.Parse(testPath2)
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
