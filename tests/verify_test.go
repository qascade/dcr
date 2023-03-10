package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyCollabPkg(t *testing.T) {

	fmt.Println("Executing TestVerifyCollabPkg")
	var testStructs = []testPackage{
		{
			path: "../samples/small_collab",
		},
	}

	for _, testStruct := range testStructs {
		testVerification(t, testStruct)
	}

}

func testVerification(t *testing.T, testStruct testPackage) {
	// Contents of GitRepo Contract and Local Contract should match.
	collabPkg := setupCollabPkg(t, testStruct)

	err := collabPkg.UploadToRepo(testStruct.path)
	require.NoError(t, err)

	err = collabPkg.Verify(testStruct.path)
	require.NoError(t, err)

	err = collabPkg.Terminate()
	require.NoError(t, err)
}
