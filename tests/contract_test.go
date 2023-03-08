package tests

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/require"
	 
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