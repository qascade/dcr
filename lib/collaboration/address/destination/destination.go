// This package will have all the different types of destinations
package destination

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/qascade/dcr/lib/collaboration/config"
)

// All Destination types must implement this interface
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Destination interface {
	GetTransformationRef() string
}

type LocalDestination struct {
	CollaboratorName  string
	DestinationName   string
	transformationRef string
}

func NewLocalDestination(cName string, dSpec config.DestinationSpec) Destination {
	return &LocalDestination{
		CollaboratorName:  cName,
		DestinationName:   dSpec.Name,
		transformationRef: dSpec.Ref,
	}
}

func (ld *LocalDestination) GetTransformationRef() string {
	return ld.transformationRef
}
