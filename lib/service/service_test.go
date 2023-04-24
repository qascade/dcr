package service

// import (
// 	"fmt"
// 	"path/filepath"
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	"github.com/qascade/dcr/lib/collaboration"
// 	"github.com/qascade/dcr/lib/collaboration/address"
// 	"github.com/qascade/dcr/lib/collaboration/config"

// )
// func TestService(t *testing.T) {
// 	fmt.Println("Running TestCollaboration")
// 	collabConfig := testConfig(t)
// 	collaboration := testAddressGraph(t, collabConfig)
// 	authorized := testAuthorization(t, collaboration)

// 	require.Equal(t, true, authorized)
// 	tAddress, err := collaboration.DeRefTransformation()
// 	require.NoError(t, err)
// 	runTransformationEvent := NewRunTransformationEvent(collabConfig, tAddress)
// 	runTransformationEvent.Execute()
// }

// func testAuthorization(t *testing.T, collaboration *collaboration.Collaboration) bool {
// 	/* test 1 */
// 	destRequester := address.AddressRef("/Research")
// 	destination := address.AddressRef("/Research/destination/customer_overlap_count")

// 	allowed, err := collaboration.AuthorizeCollaborationEvent(destRequester, destination)
// 	require.NoError(t, err)
// 	require.Equal(t, true, allowed)

// 	/*test 2 */
// 	transformationRunner := address.AddressRef("/Media")
// 	transformation := address.AddressRef("/Research/transformation/total_customers")

// 	allowed, err = collaboration.AuthorizeCollaborationEvent(transformationRunner, transformation)
// 	require.NoError(t, err)
// 	require.Equal(t, true, allowed)
// 	return allowed
// }

// func testAddressGraph(t *testing.T, collabConfig *config.CollaborationConfig) *collaboration.Collaboration {
// 	fmt.Println("Running TestAddressGraph")
// 	collaboration, err := collaboration.NewCollaboration(*collabConfig)
// 	require.NoError(t, err)
// 	fmt.Println("AddressGraph =============")
// 	fmt.Println(*collaboration.AddressGraph)
// 	require.NotNil(t, collaboration.AddressGraph)
// 	return collaboration
// }

// func testConfig(t *testing.T) *config.CollaborationConfig {
// 	fmt.Println("Running TestConfig")
// 	parser := config.ConfigParser(config.ConfigFolder{})
// 	testPath1, err := filepath.Abs("../../samples/test_collaboration")
// 	require.NoError(t, err)

// 	testPath2, err := filepath.Abs("../../samples/init_collaboration")
// 	require.NoError(t, err)

// 	_, err = parser.Parse(testPath1)
// 	require.NoError(t, err)
// 	// for _, pkg := range pkgConfig.PackagesInfo {
// 	// 	fmt.Println("SourceSpec =============")
// 	// 	fmt.Println(*pkg.SourceSpec)
// 	// 	fmt.Println("DestinationSpec =============")
// 	// 	fmt.Println(*pkg.DestinationGroupSpec)
// 	// 	fmt.Println("TransformationSpec =============")
// 	// 	fmt.Println(*pkg.TransformationGroupSpec)
// 	// }

// 	pkgConfig, err := parser.Parse(testPath2)
// 	require.NoError(t, err)
// 	for _, pkg := range pkgConfig.PackagesInfo {
// 		fmt.Println("SourceSpec =============")
// 		fmt.Println(*pkg.SourceSpec)
// 		fmt.Println("DestinationSpec =============")
// 		fmt.Println(*pkg.DestinationGroupSpec)
// 		fmt.Println("TransformationSpec =============")
// 		fmt.Println(*pkg.TransformationGroupSpec)
// 	}
// 	return pkgConfig

// }
