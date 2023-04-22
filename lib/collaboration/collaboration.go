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
	for _, neighbour := range c.AddressGraph.AdjacencyList[root] {
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		movementPermission, err := c.Authorizer(collaboratorRef, neighbour, visited)
		if !movementPermission {
			return false, err
		}
	}
	return true, nil
}

func (c *Collaboration) Authorizer(collaboratorRef address.AddressRef, root address.AddressRef, visited map[address.AddressRef]bool) (bool, error) {
	var err error
	for _, neighbour := range c.AddressGraph.AdjacencyList[root] {
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		var movementPermission bool
		if strings.Contains(string(root), "source") {
			neighbourAddress, ok := c.cachedSources[root]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing collaborator for source address. %s", root)
			movementPermission, err = neighbourAddress.Authorize(collaboratorRef)
		} else if strings.Contains(string(root), "transformation") {
			neighbourAddress, ok := c.cachedTransformations[root]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing collaborator for transformation address. %s", root)
			movementPermission, err = neighbourAddress.Authorize(collaboratorRef)
		} else if strings.Contains(string(root), "destination") {
			neighbourAddress, ok := c.cachedDestination[root]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing collaborator for destination address. %s", root)
			movementPermission, err = neighbourAddress.Authorize(collaboratorRef)
		} else {
			err = fmt.Errorf("Invalid address type. %s", root)
			log.Error(err)
			return false, err
		}
		if !movementPermission {
			return false, fmt.Errorf("Collaborator not authorised. %s", err)
		}
		return c.Authorizer(collaboratorRef, neighbour, visited)
	}
	return true, nil
}
