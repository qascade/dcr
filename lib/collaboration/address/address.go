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
	Count                 int
	AdjacencyList         map[AddressRef][]AddressRef
	AuthorityStatus       map[AddressRef]bool
	CachedSources         map[AddressRef]DcrAddress
	CachedTransformations map[AddressRef]DcrAddress
	CachedDestinations    map[AddressRef]DcrAddress
}

func NewGraph(cSources map[AddressRef]DcrAddress, cTransformations map[AddressRef]DcrAddress, cDestinations map[AddressRef]DcrAddress) (*Graph, error) {
	log.Info("Graph is being populated...")
	count := len(cSources) + len(cTransformations) + len(cDestinations)
	adjList := make(map[AddressRef][]AddressRef)
	authorityStatus := make(map[AddressRef]bool)
	for tRef, tAddressI := range cTransformations {
		authorityStatus[tRef] = false
		tAddress, ok := tAddressI.(*TransformationAddress)
		if !ok {
			log.Error("The address is not of type TransformationAddress")
			return nil, fmt.Errorf("the address is not of type TransformationAddress for addressRef: %s", tRef)
		}

		sourcesInfo := tAddress.Transformation.GetSourcesInfo()
		for _, sourceMetadata := range sourcesInfo {
			sAddress := cSources[AddressRef(sourceMetadata.AddressRef)]
			if sAddress != nil {
				adjList[tRef] = append(adjList[tRef], AddressRef(sourceMetadata.AddressRef))
				continue
			}
			// A transformation may consume another transformation output as a potential source.
			if AddressRef(sourceMetadata.AddressRef).IsTransformation() {
				fmt.Printf("Info::transformation %s is consuming transformation output %s as a source\n", tRef, sourceMetadata.AddressRef)
				adjList[tRef] = append(adjList[tRef], AddressRef(sourceMetadata.AddressRef))
			} else {
				log.Error("Source Address not found")
				return nil, fmt.Errorf("source Address not found for addressRef: %s", sourceMetadata.AddressRef)
			}
		}
	}
	for dRef, dAddressI := range cDestinations {
		authorityStatus[dRef] = false
		dAddress, ok := dAddressI.(*DestinationAddress)
		if !ok {
			log.Error("The address is not of type DestinationAddress")
			return nil, fmt.Errorf("the address is not of type DestinationAddress for addressRef: %s", dRef)
		}
		adjList[dRef] = append(adjList[dRef], AddressRef(dAddress.Destination.GetTransformationRef()))
	}

	for sRef := range cSources {
		authorityStatus[sRef] = true // Sources are always authorized.
	}

	graph := &Graph{
		Count:                 count,
		AdjacencyList:         adjList,
		AuthorityStatus:       authorityStatus,
		CachedSources:         cSources,
		CachedTransformations: cTransformations,
		CachedDestinations:    cDestinations,
	}
	return graph, nil
}

// All AddressNodeTypes must implement this interface
type DcrAddress interface {
	Authorize([]AddressRef, AddressRef) (bool, error) // Are we allowed to move further down the graph?
	//Deref  // Function that returns the real transformation
	Type() AddressType // Returns the type of Address.
}

type SourceAddress struct {
	Ref                         AddressRef
	Source                      source.Source
	Owner                       AddressRef //CollaboratorName
	TransformationOwnersAllowed []AddressRef
	DestinationsAllowed         []AddressRef
	SourceNoises                map[AddressRef]map[string]string
}

func NewSourceAddress(ref AddressRef, owner string, transformationOwnersAllowed []AddressRef, destAllowed []AddressRef, source source.Source, sourceNoises map[AddressRef]map[string]string) DcrAddress {
	// Owner is always allowed to consume its own source.
	transformationOwnersAllowed = append(transformationOwnersAllowed, NewCollaboratorRef(owner))
	return &SourceAddress{
		Ref:                         ref,
		Owner:                       NewCollaboratorRef(owner),
		TransformationOwnersAllowed: transformationOwnersAllowed,
		DestinationsAllowed:         destAllowed,
		Source:                      source,
		SourceNoises:                sourceNoises,
	}
}

// Transformaton Owner is checked against Source TransformationOwnersAllowed.
// Destination owner is checked against Source DestinationsAllowed.
func (sa *SourceAddress) Authorize(parents []AddressRef, root AddressRef) (bool, error) {
	for _, ref := range parents {
		if ref.IsDestination() {
			isAuthorized, err := sa.AuthorizeDestination(ref)
			if !isAuthorized {
				return false, err
			}
		}
		if ref.IsTransformation() {
			isAuthorized, err := sa.AuthorizeTransformation(ref)
			if !isAuthorized {
				return false, err
			}
		}
	}
	if root.IsDestination() {
		return sa.AuthorizeDestination(root)
	}
	if root.IsTransformation() {
		return sa.AuthorizeTransformation(root)
	}
	// If root is a source, it is always authorized.
	return true, nil
}

func (sa *SourceAddress) AuthorizeTransformation(root AddressRef) (bool, error) {
	tOwner := root.Collaborator()
	tAllowed := false
	for _, c := range sa.TransformationOwnersAllowed {
		if c == tOwner {
			tAllowed = true
			break
		}
	}
	if !tAllowed {
		err := fmt.Errorf("transformation: %s not allowed to consume source: %s", root, sa.Ref)
		log.Error(err)
		return false, err
	}
	log.Infof("transformation: %s allowed to consume source: %s", root, sa.Ref)
	return true, nil
}

func (sa *SourceAddress) AuthorizeDestination(root AddressRef) (bool, error) {
	dAllowed := false
	for _, ref := range sa.DestinationsAllowed {
		if root == ref {
			dAllowed = true
			break
		}
	}
	if !dAllowed {
		err := fmt.Errorf("destination: %s not allowed to consume source: %s", root, sa.Ref)
		return false, err
	}
	log.Infof("destination: %s allowed to consume source: %s", root, sa.Ref)
	return true, nil
}

func (sa *SourceAddress) Type() AddressType {
	return ADDRESS_TYPE_SOURCE
}

type TransformationAddress struct {
	Ref                      AddressRef
	Owner                    AddressRef
	DestinationOwnersAllowed []AddressRef
	DestinationsAllowed      []AddressRef
	Transformation           transformation.Transformation
	NoiseParams              []string
	NoiseType                string
}

func NewTransformationAddress(ref AddressRef, owner string, destinationOwnersAllowed []AddressRef, destAllowed []AddressRef, t transformation.Transformation, noiseParams []string) DcrAddress {
	// Owner is always allowed to consume its own transformation.
	destinationOwnersAllowed = append(destinationOwnersAllowed, NewCollaboratorRef(owner))
	destAllowed = append(destAllowed, NewCollaboratorRef(owner))
	return &TransformationAddress{
		Ref:                      ref,
		Owner:                    NewCollaboratorRef(owner),
		DestinationOwnersAllowed: destinationOwnersAllowed,
		DestinationsAllowed:      destAllowed,
		Transformation:           t,
		NoiseParams:              noiseParams,
	}
}

// Destination Owners are to checked against transformation DestinationOwnersAllowed.
func (ta *TransformationAddress) Authorize(parents []AddressRef, root AddressRef) (bool, error) {
	log.Infof("Root %s is trying to consume transformation %s, Performing authorization", root, ta.Ref)
	for _, ref := range parents {
		if ref.IsDestination() {
			isAuthorized, err := ta.AuthorizeDestination(ref)
			if !isAuthorized {
				return false, err
			}
		}
	}
	if root.IsDestination() {
		return ta.AuthorizeDestination(root)
	}
	// If root is not a destination, it must be transformation itself.
	return true, nil
}

func (ta *TransformationAddress) AuthorizeDestination(root AddressRef) (bool, error) {
	dAllowed := false
	dOwner := root.Collaborator()
	for _, ref := range ta.DestinationOwnersAllowed {
		if dOwner == ref {
			dAllowed = true
			break
		}
	}
	if !dAllowed {
		err := fmt.Errorf("destination: %s not allowed to consume transformation: %s", root, ta.Ref)
		log.Error(err)
		return false, err
	}
	log.Infof("destination: %s allowed to consume transformation: %s", root, ta.Ref)
	return true, nil
}

func (ta *TransformationAddress) Type() AddressType {
	return ADDRESS_TYPE_TRANSFORMATION
}

type DestinationAddress struct {
	Ref         AddressRef
	Owner       AddressRef
	Destination destination.Destination
}

func NewDestinationAddress(ref AddressRef, owner AddressRef, dest destination.Destination) DcrAddress {
	return &DestinationAddress{
		Ref:         ref,
		Owner:       AddressRef(owner),
		Destination: dest,
	}
}

func (da *DestinationAddress) Authorize(_ []AddressRef, _ AddressRef) (bool, error) {
	// Movement from destination is always authorized.
	log.Infof("Performing authorization for Destination %s.", da.Ref)
	return true, nil
}

// Helper functions

func (da *DestinationAddress) Type() AddressType {
	return ADDRESS_TYPE_DESTINATION
}

func getAddressRefSlice(s []string) []AddressRef {
	addRefS := make([]AddressRef, 0)
	for _, str := range s {
		addRefS = append(addRefS, NewCollaboratorRef(str))
	}
	return addRefS
}

// Returns a slice of AddressRef from a slice of SourceDestinationAllowedSpec
func getTransformationRefSlice(destAllowed []config.SourceDestinationAllowedSpec) []AddressRef {
	addS := make([]AddressRef, 0)
	for _, dest := range destAllowed {
		addS = append(addS, AddressRef(dest.Ref))
	}
	return addS
}
