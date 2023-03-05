package tests

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/qascade/dcr/core/contract"
)

type testPackage struct {
	path string
}

func TestContract(t *testing.T) {
	var testStructs []testPackage = []testPackage{
		{
			path: "../samples/small_collab",
		},
	}
	for _, testStruct := range testStructs {
		tempContract, err := contract.ParseContract(testStruct.path)
		require.NoError(t, err)
		err = tempContract.Validate()
		require.NoError(t, err)
	}

}
