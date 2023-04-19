package collaboration

import (
	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/collaboration/config"
	"github.com/qascade/dcr/lib/collaboration/destination"
	"github.com/qascade/dcr/lib/collaboration/source"
	"github.com/qascade/dcr/lib/collaboration/transformation"
)

type Collaboration struct {
	AddressGraph          *address.Graph
	Collaborators         []string
	collaborationConfig   config.CollaborationConfig
	cachedSources         map[string]source.Source
	cachedTransformations map[string]transformation.Transformation
	cachedDestination     map[string]destination.Destination
}

func NewCollaboration(collabConfig config.CollaborationConfig) *Collaboration {
	var collaborators []string
	for _, pkgConfig := range collabConfig.PackagesInfo {
		collaboratorName := pkgConfig.PackageMetadata.CollaboratorName
		collaborators = append(collaborators, collaboratorName)
	}
	graph := address.NewGraph()
	return &Collaboration{
		AddressGraph:        graph,
		Collaborators:       collaborators,
		collaborationConfig: collabConfig,
	}
}
