package address

import (
	"fmt"
	"os"

	"github.com/edwingeng/deque"
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
	Count         int
	AdjacencyList map[AddressRef][]AddressRef
	IndegreeList  map[AddressRef]int
	TopoOrder     []AddressRef
}

func NewGraph(cSources map[AddressRef]DcrAddress, cTransformations map[AddressRef]DcrAddress, cDestinations map[AddressRef]DcrAddress) (*Graph, error) {
	log.Info("Graph is being populated...")
	count := len(cSources) + len(cTransformations) + len(cDestinations)
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
			if sAddress != nil {
				adjList[tRef] = append(adjList[tRef], AddressRef(sourceMetadata.AddressRef))
				continue
			}
			// A transformation may consume destination as a potential source.
			if AddressRef(sourceMetadata.AddressRef).IsDestination() {
				fmt.Printf("Info::transformation %s has a destination %s as a source\n", tRef, sourceMetadata.AddressRef)
				adjList[tRef] = append(adjList[tRef], AddressRef(sourceMetadata.AddressRef))
			} else {
				log.Error("Source Address not found")
				return nil, fmt.Errorf("source Address not found for addressRef: %s", sourceMetadata.AddressRef)
			}
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

	indegreeList := make(map[AddressRef]int)
	// Create indegreeList
	for v, neighbourList := range adjList {
		if _, ok := indegreeList[v]; !ok {
			indegreeList[v] = 0
		}
		for _, neighbour := range neighbourList {
			if _, ok := indegreeList[neighbour]; !ok {
				indegreeList[neighbour] = 0
			}
			indegreeList[neighbour]++
		}
	}

	// Do topological sort
	topoOrder, err := topoSort(count, adjList, indegreeList)
	if err != nil {
		err = fmt.Errorf("not able to do topological sort, %v", err)
		log.Error(err)
		return nil, err
	}
	graph := &Graph{
		Count:         count,
		AdjacencyList: adjList,
		IndegreeList:  indegreeList,
		TopoOrder:     topoOrder,
	}
	return graph, nil
}

func topoSort(count int, adjList map[AddressRef][]AddressRef, indegreeList map[AddressRef]int) ([]AddressRef, error) {
	// Create a queue and enqueue all vertices with indegree 0
	var queue deque.Deque = deque.NewDeque()
	for k, v := range indegreeList {
		if v == 0 {
			queue.PushBack(k)
		}
	}
	var topoOrder []AddressRef

	for queue.Len() != 0 {
		// Dequeue a vertex from queue and add it to topoOrder
		v, ok := queue.PopFront().(AddressRef)
		if !ok {
			err := fmt.Errorf("could not convert to AddressRef, %v", v)
			log.Error(err)
			return nil, err
		}
		topoOrder = append(topoOrder, v)
		// Iterate through all its neighbouring nodes of dequeued node u and decrease their in-degree by 1
		for _, neighbour := range adjList[v] {
			indegreeList[neighbour]--
			// If in-degree becomes zero, add it to queue
			if indegreeList[neighbour] == 0 {
				queue.PushBack(neighbour)
			}
		}
	}
	// Check if there was a cycle
	if len(topoOrder) != count {
		err := fmt.Errorf("there exists a cycle in the graph")
		log.Error(err)
		return nil, err
	}
	return topoOrder, nil
}

// All AddressNodeTypes must implement this interface
type DcrAddress interface {
	Authorize(AddressRef, AddressRef) (bool, error) // Is a Collaborator Allowed to Move further into the graph.
	//Deref  // Function that returns the real transformation
	Type() AddressType // Returns the type of Address.
}

type SourceAddress struct {
	Ref                 AddressRef
	Source              source.Source
	Owner               AddressRef //CollaboratorName
	ConsumersAllowed    []AddressRef
	DestinationsAllowed []AddressRef
	SourceNoises        map[AddressRef]map[string]string
}

func NewSourceAddress(ref AddressRef, owner string, consumersAllowed []AddressRef, destAllowed []AddressRef, source source.Source, sourceNoises map[AddressRef]map[string]string) DcrAddress {
	// Owner is always allowed to consume its own source.
	consumersAllowed = append(consumersAllowed, NewCollaboratorRef(owner))
	return &SourceAddress{
		Ref:                 ref,
		Owner:               NewCollaboratorRef(owner),
		ConsumersAllowed:    consumersAllowed,
		DestinationsAllowed: destAllowed,
		Source:              source,
		SourceNoises:        sourceNoises,
	}
}

func (sa *SourceAddress) Authorize(collabName AddressRef, tName AddressRef) (bool, error) {
	log.Infof("Collaborator %s is trying to consume source %s, Performing authorization", collabName, sa.Ref)
	collabAllowed := false
	destAllowed := false
	for _, allowedCollab := range sa.ConsumersAllowed {
		if collabName == allowedCollab {
			collabAllowed = true
		}
	}
	for _, allowedDest := range sa.DestinationsAllowed {
		if tName == allowedDest {
			destAllowed = true
		}
	}
	if collabAllowed && destAllowed {
		log.Infof("Collaborator %s is allowed to consume source %s. Authorization Successful", collabName, sa.Ref)
		return true, nil
	}
	err := fmt.Errorf("collaborator %s is not allowed to consume source %s", collabName, sa.Ref)
	log.Error(err)
	return false, err
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
	NoiseParams         []string
	NoiseType           string
}

func NewTransformationAddress(ref AddressRef, owner string, consumersAllowed []AddressRef, destAllowed []AddressRef, t transformation.Transformation, noiseParams []string) DcrAddress {
	// Owner is always allowed to consume its own transformation.
	consumersAllowed = append(consumersAllowed, NewCollaboratorRef(owner))
	destAllowed = append(destAllowed, NewCollaboratorRef(owner))
	return &TransformationAddress{
		Ref:                 ref,
		Owner:               NewCollaboratorRef(owner),
		ConsumersAllowed:    consumersAllowed,
		DestinationsAllowed: destAllowed,
		Transformation:      t,
		NoiseParams:         noiseParams,
	}
}
func (ta *TransformationAddress) Authorize(collabName AddressRef, _ AddressRef) (bool, error) {
	log.Infof("Collaborator %s is trying to consume transformation %s, Performing authorization", collabName, ta.Ref)
	collabAllowed := false

	for _, allowedCollab := range ta.ConsumersAllowed {
		if collabName == allowedCollab {
			collabAllowed = true
		}
	}
	if collabAllowed {
		log.Infof("Collaborator %s is allowed to consume transformation %s. Authorization Successful", collabName, ta.Ref)
		return true, nil
	}
	return false, fmt.Errorf("collaborator %s is not allowed to consume transformation %s", collabName, ta.Ref)
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

func (da *DestinationAddress) Authorize(collabName AddressRef, _ AddressRef) (bool, error) {
	// Movement from destination is always authorized.
	log.Infof("Collaborator %s is trying to consume destination %s, Performing authorization", collabName, da.Ref)
	return true, nil
}

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

func getTransformationRefSlice(destAllowed []config.SourceDestinationAllowedSpec) []AddressRef {
	addS := make([]AddressRef, 0)
	for _, dest := range destAllowed {
		addS = append(addS, AddressRef(dest.Ref))
	}
	return addS
}
