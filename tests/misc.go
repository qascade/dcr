package tests 

import ( 
	"testing"
	"github.com/stretchr/testify/require"
	
	collab "github.com/qascade/dcr/collaboration"
)


type testPackage struct {
	path string
}

func setupCollabPkg(t *testing.T, testStruct testPackage) *collab.CollaborationPackage {
	collabPkg, err := collab.NewCollaborationPkg(testStruct.path)
	require.NoError(t, err)
	return collabPkg
}