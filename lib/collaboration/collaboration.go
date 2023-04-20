package collaboration

import (
	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/collaboration/address/destination"
	"github.com/qascade/dcr/lib/collaboration/address/source"
	"github.com/qascade/dcr/lib/collaboration/address/transformation"
	"github.com/qascade/dcr/lib/collaboration/config"
)

type Collaboration struct {
	AddressGraph          *address.Graph
	Collaborators         []string
	collaborationConfig   config.CollaborationConfig
	cachedSources         map[address.AddressRef]source.Source
	cachedTransformations map[address.AddressRef]transformation.Transformation
	cachedDestination     map[address.AddressRef]destination.Destination
}

func NewCollaboration(collabConfig config.CollaborationConfig) *Collaboration {
	var collaborators []string
	for _, pkgConfig := range collabConfig.PackagesInfo {
		collaboratorName := pkgConfig.CollaboratorName
		collaborators = append(collaborators, collaboratorName)
	}
	graph := address.NewGraph(collabConfig)
	cSources, cTransformations, cDestinations := address.CacheAddresses(collabConfig)
	return &Collaboration{
		AddressGraph:          graph,
		Collaborators:         collaborators,
		collaborationConfig:   collabConfig,
		cachedSources:         cSources,
		cachedTransformations: cTransformations,
		cachedDestination:     cDestinations,
	}
}
