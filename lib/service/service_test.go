package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	testPkgPath := "../../samples/init_collaboration"
	runner := "Media"
	destinationOwner := "Research"
	tRef := "/Research/transformation/private_total_customers"
	destRef := "/Research/destination/customer_overlap_count"

	_, err := NewService(testPkgPath, runner, destinationOwner, tRef, destRef)
	//service, err := NewService(testPkgPath, runner, destinationOwner, tRef, destRef)
	require.NoError(t, err)

	// err = service.RunCollaborationEvent()
	// require.NoError(t, err)
}
