package address

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/qascade/dcr/lib/collaboration/address/destination"
	"github.com/qascade/dcr/lib/collaboration/address/source"
	"github.com/qascade/dcr/lib/collaboration/address/transformation"
	"github.com/qascade/dcr/lib/collaboration/config"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Graph struct {
	AdjacencyList map[AddressRef][]AddressRef
}

func NewGraph(cSources map[AddressRef]DcrAddress, cTransformations map[AddressRef]DcrAddress, cDestinations map[AddressRef]DcrAddress) (*Graph, error) {
	log.Info("Graph is being populated...")
	adjList := make(map[AddressRef][]AddressRef)
	for tRef, tAddressI := range cTransformations {
		tAddress, ok := tAddressI.(*TransformationAddress)
		if !ok {
			log.Error("The address is not of type TransformationAddress")
			return nil, fmt.Errorf("the address is not of type TransformationAddress for addressRef: %s", tRef)
		}
		sourcesInfo := tAddress.Transformation.GetSourcesInfo()
		for _, sourceMetadata := range sourcesInfo {
			sAddress := cSources[AddressRef(sourceMetadata.AddressRef)]
			if sAddress == nil {
				log.Error("Source Address not found")
				return nil, fmt.Errorf("source Address not found for addressRef: %s", sourceMetadata.AddressRef)
			}
			adjList[tRef] = append(adjList[tRef], AddressRef(sourceMetadata.AddressRef))
		}
	}
	for dRef, dAddressI := range cDestinations {
		dAddress, ok := dAddressI.(*DestinationAddress)
		if !ok {
			log.Error("The address is not of type DestinationAddress")
			return nil, fmt.Errorf("the address is not of type DestinationAddress for addressRef: %s", dRef)
		}
		adjList[dRef] = append(adjList[dRef], AddressRef(dAddress.Destination.GetTransformationRef()))
	}
	graph := &Graph{
		AdjacencyList: adjList,
	}
	return graph, nil
}

// All AddressNodeTypes must implement this interface
type DcrAddress interface {
	Authorize(AddressRef) (bool, error) // Is a Collaborator Allowed to Move further into the graph.
	//Deref  // Function that returns the real transformation
	Type() AddressType // Returns the type of Address.
}

type SourceAddress struct {
	Ref                 AddressRef
	Source              source.Source
	Owner               AddressRef //CollaboratorName
	ConsumersAllowed    []AddressRef
	DestinationsAllowed []AddressRef
}

func NewSourceAddress(ref AddressRef, owner string, consumersAllowed []AddressRef, destAllowed []AddressRef, source source.Source) DcrAddress {
	// Owner is always allowed to consume its own source.
	consumersAllowed = append(consumersAllowed, NewCollaboratorRef(owner))
	return &SourceAddress{
		Ref:                 ref,
		Owner:               NewCollaboratorRef(owner),
		ConsumersAllowed:    consumersAllowed,
		DestinationsAllowed: destAllowed,
		Source:              source,
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
	ConsumersAllowed    []AddressRef
	DestinationsAllowed []AddressRef
	Transformation      transformation.Transformation
}

func NewTransformationAddress(ref AddressRef, owner string, consumersAllowed []AddressRef, destAllowed []AddressRef, t transformation.Transformation) DcrAddress {
	// Owner is always allowed to consume its own transformation.
	consumersAllowed = append(consumersAllowed, NewCollaboratorRef(owner))
	destAllowed = append(destAllowed, NewCollaboratorRef(owner))
	return &TransformationAddress{
		Ref:                 ref,
		Owner:               NewCollaboratorRef(owner),
		ConsumersAllowed:    consumersAllowed,
		DestinationsAllowed: destAllowed,
		Transformation:      t,
	}
}
func (ta *TransformationAddress) Authorize(collabName AddressRef) (bool, error) {
	log.Info("Authorize for Transformation Address still needs to be implemented")
	return true, nil
}

func (ta *TransformationAddress) Type() AddressType {
	return ADDRESS_TYPE_TRANSFORMATION
}

type DestinationAddress struct {
	Ref         AddressRef
	Requester   AddressRef
	Destination destination.Destination
}

func NewDestinationAddress(ref AddressRef, requester AddressRef, dest destination.Destination) DcrAddress {
	return &DestinationAddress{
		Ref:         ref,
		Requester:   AddressRef(requester),
		Destination: dest,
	}
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
	for _, str := range s {
		addRefS = append(addRefS, AddressRef(str))
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
