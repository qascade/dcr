package tests

import (
	"path/filepath"
	"testing"

	"github.com/qascade/dcr/collaboration"
)

func TestExecuteWarehouse(t *testing.T) {

	sqlFilePath := filepath.Join("/home/shisuimadara/dcr/collaboration/test.sql")

	c := collaboration.CollaborationPackage{}
	err := c.ExecuteWarehouse(sqlFilePath)
	if err != nil {
		t.Errorf("Unexpected error executing warehouse: %v", err)
	}
}
