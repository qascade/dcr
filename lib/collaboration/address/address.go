package address

import (
	"os"

	"github.com/qascade/dcr/lib/collaboration/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Graph struct {
	AdjacencyList map[DcrAddress][]DcrAddress
}

var adjList map[DcrAddress][]DcrAddress

// All Data Access Grants are to be the part of DcrAddress not source, transformation interface.

func NewGraph(collabConfig config.CollaborationConfig) *Graph {
	log.Info("a new Graph function yet to be implemented")
	adjList = make(map[DcrAddress][]DcrAddress)
	return &Graph{
		AdjacencyList: adjList,
	}
}

// All AddressNodeTypes must implement this interface
type DcrAddress interface {
	Authorize(AddressRef) (bool, error) // Is a Collaborator Allowed to Move further into the graph.
	//Deref  // Function that returns the real transformation
	Type() AddressType // Returns the type of Address.
}

type SourceAddress struct {
	Ref AddressRef
	//Source              *source.Source
	Owner               AddressRef //CollaboratorName
	ConsumersAllowed    []AddressRef
	DestinationsAllowed []AddressRef
}

func NewSourceAddress(ref string, owner string, consumersAllowed []AddressRef, destAllowed []AddressRef) DcrAddress {
	return &SourceAddress{
		Ref:                 AddressRef(ref),
		Owner:               NewCollaboratorRef(owner),
		ConsumersAllowed:    consumersAllowed,
		DestinationsAllowed: destAllowed,
	}
}
func (sa *SourceAddress) Authorize(collabName AddressRef) (bool, error) {
	log.Info("Authorize for SourceAddress still needs to be implemented")
	return true, nil
}

func (sa *SourceAddress) Type() AddressType {
	return ADDRESS_TYPE_SOURCE
}

type TransformationAddress struct {
	Ref                 AddressRef
	Owner               AddressRef
	Runner              AddressRef
	ConsumersAllowed    []AddressRef
	DestinationsAllowed []AddressRef
	//Transformation      *transformation.Transformation
}

func (ta *TransformationAddress) Authorize(collabName AddressRef) (bool, error) {
	log.Info("Authorize for Transformation Address still needs to be implemented")
	return true, nil
}

func (ta *TransformationAddress) Type() AddressType {
	return ADDRESS_TYPE_TRANSFORMATION
}

type DestinationAddress struct {
	Ref       AddressRef
	Requester AddressRef
	//Destination *destination.Destination
}

func (da *DestinationAddress) Authorize(collabName AddressRef) (bool, error) {
	log.Info("Authorize for destination still needs to be implemented")
	return true, nil
}

func (da *DestinationAddress) Type() AddressType {
	return ADDRESS_TYPE_DESTINATION
}

func getAddressRefSlice(s []string) []AddressRef {
	addRefS := make([]AddressRef, 0)
	for c, e := range s {
		addRefS[c] = AddressRef(e)
	}
	return addRefS
}

func getTransformationRefSlice(destAllowed []config.SourceDestinationAllowedSpec) []AddressRef {
	addS := make([]AddressRef, 0)
	for _, dest := range destAllowed {
		addS = append(addS, AddressRef(dest.Ref))
	}
	return addS
}
