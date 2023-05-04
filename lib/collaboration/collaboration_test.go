package collaboration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/qascade/dcr/lib/collaboration/config"
)

func init() {
	log.SetLevel(log.ErrorLevel)
	log.SetOutput(os.Stdout)
}

func TestGraph(t *testing.T) {
	fmt.Println("Running TestCollaboration")
	testOneDepthPath, err := filepath.Abs("../../samples/test_graph/test_onedepth")
	require.NoError(t, err)

	testAddressGraph(t, testOneDepthPath)

	testTwoDepthPath, err := filepath.Abs("../../samples/test_graph/test_twodepth")
	require.NoError(t, err)

	testAddressGraph(t, testTwoDepthPath)
}

func TestParse(t *testing.T) {
	testOneDepthPath, err := filepath.Abs("../../samples/init_collaboration")
	require.NoError(t, err)
	log.Printf("Running TestParse for pkgPath = %s\n", testOneDepthPath)
	testConfig(t, testOneDepthPath)

	testTwoDepthPath, err := filepath.Abs("../../samples/test_graph/test_twodepth")
	require.NoError(t, err)
	log.Printf("Running TestParse for pkgPath = %s\n", testTwoDepthPath)
	testConfig(t, testTwoDepthPath)
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
	fmt.Println("TopoOrder =============")
	fmt.Println(graph.TopoOrder)

	require.NotNil(t, collaboration.AddressGraph)
	return collaboration
}

func testConfig(t *testing.T, path string) *config.CollaborationConfig {
	log.Printf("Running TestConfig, path = %s\n", path)
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
