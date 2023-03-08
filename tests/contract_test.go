package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	collab "github.com/qascade/dcr/collaboration"
)

type testPackage struct {
	path string
}

func TestContract(t *testing.T) {
	var testStructs = []testPackage{
		{
			path: "../samples/small_collab",
		},
	}

	for _, testStruct := range testStructs {
		testParsing(t, testStruct)
		testVerification(t, testStruct)
	}

}

func testParsing(t *testing.T, testStruct testPackage) {
	var collabPkg collab.CollaborationParser = &collab.CollaborationPackage{}
	cSpec, tSpecs, err := collabPkg.Parse(testStruct.path)
	require.NoError(t, err)
	fmt.Println(cSpec)
	fmt.Println()
	for _, tSpec := range *tSpecs {
		fmt.Println(tSpec)
		fmt.Println()
	}
	require.NotNil(t, cSpec)
	require.NotNil(t, tSpecs)
}

func testVerification(t *testing.T, testStruct testPackage) {
	// Contents of GitRepo Contract and Local Contract should match. 
	err := collab.Verify(testStruct.path)
	require.NoError(t, err)
}