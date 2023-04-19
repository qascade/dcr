package address

import (
	"os"

	"github.com/qascade/dcr/lib/collaboration/destination"
	"github.com/qascade/dcr/lib/collaboration/source"
	"github.com/qascade/dcr/lib/collaboration/transformation"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Graph map[DcrAddress][]DcrAddress

func NewGraph() *Graph {
	return &Graph{}
}

type CollaboratorName string

// All AddressNodeTypes must implement this interface
type DcrAddress interface {
	Authorize() (bool, error) // Is a Collaborator Allowed to Move further into the graph.
	//Deref  // Function that returns the real transformation
	Type() AddressType // Returns the type of Address.
}

type SourceAddress struct {
	Owner               CollaboratorName
	ConsumersAllowed    []CollaboratorName
	DestinationsAllowed []CollaboratorName
	Source              *source.Source
}

func (sa *SourceAddress) Authorize(collabName CollaboratorName) (bool, error) {
	return true, nil
}

func (sa *SourceAddress) Type() AddressType {
	return ADDRESS_TYPE_SOURCE
}

type TransformationAddress struct {
	Owner               CollaboratorName
	Runner              CollaboratorName
	ConsumersAllowed    []CollaboratorName
	DestinationsAllowed []CollaboratorName
	Transformation      *transformation.Transformation
}

func (ta *TransformationAddress) Authorize(collabName CollaboratorName) (bool, error) {
	return true, nil
}

func (ta *TransformationAddress) Type() AddressType {
	return ADDRESS_TYPE_TRANSFORMATION
}

type DestinationAddress struct {
	Requester   CollaboratorName
	Destination *destination.Destination
}

func (da *DestinationAddress) Authorize(collabName CollaboratorName) (bool, error) {
	return true, nil
}

func (da *DestinationAddress) Type() AddressType {
	return ADDRESS_TYPE_DESTINATION
}
