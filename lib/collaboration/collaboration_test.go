package collaboration

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qascade/dcr/lib/collaboration/config"
)

func TestGraph(t *testing.T) {
	fmt.Println("Running TestCollaboration")
	// testInitCollaborationPath, err := filepath.Abs("../../samples/init_collaboration")
	// require.NoError(t, err)

	// testAddressGraph(t, testInitCollaborationPath)

	// testOneDepthPath, err := filepath.Abs("../../samples/test_graph/test_onedepth")
	// require.NoError(t, err)

	// testAddressGraph(t, testOneDepthPath)

	testTwoDepthPath, err := filepath.Abs("../../samples/test_graph/test_twodepth")
	require.NoError(t, err)

	testAddressGraph(t, testTwoDepthPath)
}

func TestParse(t *testing.T) {
	testOneDepthPath, err := filepath.Abs("../../samples/init_collaboration")
	require.NoError(t, err)
	testConfig(t, testOneDepthPath)

	// testTwoDepthPath, err := filepath.Abs("../../samples/test_graph/test_twodepth")
	// require.NoError(t, err)
	// testConfig(t, testTwoDepthPath)
}

func testAddressGraph(t *testing.T, pkgPath string) *Collaboration {
	fmt.Printf("Running TestAddressGraph, for pkgPath = %s\n", pkgPath)
	collaboration, err := NewCollaboration(pkgPath)
	require.NoError(t, err)
	graph := *collaboration.AddressGraph

	fmt.Println("AddressGraph =============")
	fmt.Println("Count of AddressNodes =============")
	fmt.Println(graph.Count)
	fmt.Println("AdjacencyList =============")
	fmt.Println(graph.AdjacencyList)

	refs, err := graph.GetOrderedRunnableRefs()
	require.NoError(t, err)
	fmt.Println("OrderedRunnableRefs =============")
	fmt.Println(refs)
	fmt.Println("AuthorityStatus =============")
	fmt.Println(graph.AuthorityStatus)
	require.NotNil(t, collaboration.AddressGraph)
	return collaboration
}

func testConfig(t *testing.T, path string) *config.CollaborationConfig {
	parser := config.Parser(config.NewCollaborationConfig())
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
