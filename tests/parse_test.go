package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCollabPkg(t *testing.T) {
	fmt.Println("Executing TestParseCollabPkg")
	var testStructs = []testPackage{
		{
			path: "../samples/small_collab",
		},	
	}
	for _, testStruct := range testStructs {
		testParsing(t, testStruct)
	}

}

// TODO - Add custom stubs to compare. 
func testParsing(t *testing.T, testStruct testPackage) {
	collabPkg := setupCollabPkg(t, testStruct)
	cSpec, tSpecs, err := collabPkg.Parse(testStruct.path)
	require.NoError(t, err)
	fmt.Println(cSpec)
	fmt.Println()
	for _, tSpec := range *tSpecs {
		fmt.Println(tSpec)
	}
	require.NotNil(t, cSpec)
	require.NotNil(t, tSpecs)
}



