package tests

import (
	"fmt"
	"testing"

	"github.com/qascade/dcr/collaboration/contract"
	"github.com/stretchr/testify/require"
)

func TestContracter(t *testing.T) {
	fmt.Println("Executing TestContracter")
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
	contracter := contract.Contracter(&contract.Contract{})
	cSpec, tSpecs, err := contracter.Parse(testStruct.path)
	require.NoError(t, err)
	fmt.Println(cSpec)
	fmt.Println()
	for _, tSpec := range *tSpecs {
		fmt.Println(tSpec)
	}
	require.NotNil(t, cSpec)
	require.NotNil(t, tSpecs)
}
