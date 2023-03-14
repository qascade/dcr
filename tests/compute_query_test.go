package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/qascade/dcr/collaboration"
	"github.com/stretchr/testify/require"
)

func TestExecuteWarehouse(t *testing.T) {

	homDir, err := os.UserHomeDir()
	require.NoError(t, err, "Unexpected error getting user home directory")

	sqlFilePath := filepath.Join(homDir, "dcr", "collaboration", "test.sql")

	c := collaboration.CollaborationPackage{}
	err = c.ExecuteSql(sqlFilePath)
	require.NoError(t, err, "Unexpected error executing warehouse")
}
