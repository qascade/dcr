package collaboration

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/collaboration/config"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Collaboration struct {
	AddressGraph          *address.Graph
	Collaborators         []string
	collaborationConfig   config.CollaborationConfig
	cachedSources         map[address.AddressRef]address.DcrAddress
	cachedTransformations map[address.AddressRef]address.DcrAddress
	cachedDestination     map[address.AddressRef]address.DcrAddress
}

func NewCollaboration(collabConfig config.CollaborationConfig) (*Collaboration, error) {
	var collaborators []string
	for _, pkgConfig := range collabConfig.PackagesInfo {
		collaboratorName := pkgConfig.CollaboratorName
		collaborators = append(collaborators, collaboratorName)
	}
	cSources, cTransformations, cDestinations := address.CacheAddresses(collabConfig)
	graph, err := address.NewGraph(cSources, cTransformations, cDestinations)
	if err != nil {
		return nil, err
	}
	collaboration := &Collaboration{
		AddressGraph:          graph,
		Collaborators:         collaborators,
		collaborationConfig:   collabConfig,
		cachedSources:         cSources,
		cachedTransformations: cTransformations,
		cachedDestination:     cDestinations,
	}
	return collaboration, nil
}

func (c *Collaboration) AuthorizeCollaborationEvent(collaboratorRef address.AddressRef, root address.AddressRef) (bool, error) {
	// Authorization is permissible only for transformation and destination addresses.
	// Source addresses are not authorized to access any other address.
	// If collaborator wants to run transformation, it should pass the transformation address as root.
	// If collaborator wants to download a destination it should pass the destination address as root.
	visited := make(map[address.AddressRef]bool)
	var parentTransformationRef address.AddressRef
	if strings.Contains(string(root), "transformation") {
		parentTransformationRef = root
	}
	if strings.Contains(string(root), "destination") {
		// As movement from destination to Transformation is always allowed. We need to store its neighbouring transformation.
		dAddress, ok := c.cachedDestination[root].(*address.DestinationAddress)
		if !ok {
			return false, fmt.Errorf("could not cast address to destination address type: %v", string(root))
		}
		parentTransformationRef = address.AddressRef(dAddress.Destination.GetTransformationRef())
	}
	for _, neighbour := range c.AddressGraph.AdjacencyList[root] {
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		movementPermission, err := c.Authorizer(collaboratorRef, neighbour, visited, parentTransformationRef)
		if !movementPermission {
			return false, err
		}
	}
	return true, nil
}

func (c *Collaboration) Authorizer(collaboratorRef address.AddressRef, root address.AddressRef, visited map[address.AddressRef]bool, parentTransformationRef address.AddressRef) (bool, error) {
	var err error
	var movementPermission bool
	if strings.Contains(string(root), "source") {
		sAddress, ok := c.cachedSources[root]
		if !ok {
			return false, fmt.Errorf("address with given address ref not found. %s", root)
		}
		log.Infof("Authorizing collaborator for source address. %s", root)
		return sAddress.Authorize(collaboratorRef, parentTransformationRef)
	}
	for _, neighbour := range c.AddressGraph.AdjacencyList[root] {
		fmt.Println(c.AddressGraph.AdjacencyList[root])
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		if strings.Contains(string(root), "transformation") {
			neighbourAddress, ok := c.cachedTransformations[root]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing collaborator for transformation address. %s", root)
			movementPermission, err = neighbourAddress.Authorize(collaboratorRef, "")
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else if strings.Contains(string(root), "destination") {
			neighbourAddress, ok := c.cachedDestination[root]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing collaborator for destination address. %s", root)
			movementPermission, err = neighbourAddress.Authorize(collaboratorRef, "")
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else {
			err = fmt.Errorf("invalid address type. %s", root)
			log.Error(err)
			return false, err
		}
		if !movementPermission {
			return false, err
		}
		movementPermission, err = c.Authorizer(collaboratorRef, neighbour, visited, parentTransformationRef)
	}
	return movementPermission, err
}

func (c *Collaboration) DeRefTransformation(ref address.AddressRef) (*address.DcrAddress, error) {
	if add, ok := c.cachedTransformations[ref]; ok {
		return &add, nil
	}
	return nil, fmt.Errorf("transformation address not found. %s", ref)
}

func (c *Collaboration) CompileTransformation(sAddRef1 address.AddressRef, sAddRef2 address.AddressRef, tAddress address.AddressRef) error {
	sAdd1I, ok := c.cachedSources[sAddRef1]
	if !ok {
		return fmt.Errorf("source address not found. %s", sAddRef1)
	}
	sAdd2I, ok := c.cachedSources[sAddRef2]
	if !ok {
		return fmt.Errorf("source address not found. %s", sAddRef2)
	}

	sAdd1, ok := sAdd1I.(*address.SourceAddress)
	if !ok {
		return fmt.Errorf("could not cast address to source address type: %v", sAddRef1)
	}
	sAdd2, ok := sAdd2I.(*address.SourceAddress)
	if !ok {
		return fmt.Errorf("could not cast address to source address type: %v", sAddRef2)
	}

	noiseParams, err := matchNoiseParameters(sAdd1, sAdd2, &tAddress)
	if err != nil {
		return err
	}

	tAddI, ok := c.cachedTransformations[tAddress]
	if !ok {
		return fmt.Errorf("transformation address not found. %s", tAddress)
	}
	tAdd, ok := tAddI.(*address.TransformationAddress)
	if !ok {
		return fmt.Errorf("could not cast address to transformation address type: %v", tAddress)
	}
	pongoInputs := tAdd.Transformation.GetPongoInputs()
	//_ := tAdd.Transformation.AppLocation()
	for k, v := range noiseParams {
		vS := fmt.Sprintf("%v", v)
		pongoInputs[k] = vS
	}

	return nil
}

func matchNoiseParameters(sAdd1 *address.SourceAddress, sAdd2 *address.SourceAddress, tAdd *address.AddressRef) (map[string]interface{}, error) {
	m1 := sAdd1.NoiseParams[*tAdd]
	m2 := sAdd2.NoiseParams[*tAdd]

	for k, v := range m1 {
		if v != m2[k] {
			return nil, fmt.Errorf("noise parameters do not match. %s", *tAdd)
		}
	}
	return m1, nil
}
