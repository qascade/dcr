package collaboration

import (
	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/collaboration/config"
)

type Collaboration struct {
	RootAddress         *address.DcrAddress
	Collaborators       []string
	collaborationConfig config.CollaborationConfig
}

func NewCollaboration(collabConfig config.CollaborationConfig) *Collaboration {
	var collaborators []string
	for _, pkgConfig := range collabConfig.PackagesInfo {
		collaboratorName := pkgConfig.PackageMetadata.CollaboratorName
		collaborators = append(collaborators, collaboratorName)
	}
	rootAddress := address.NewAddress()
	return &Collaboration{
		RootAddress:   &rootAddress,
		Collaborators: collaborators,
	}
}
