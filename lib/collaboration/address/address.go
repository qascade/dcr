package address

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Parent interface {
	ParentAddress() *ParentAddress
	Name() string
}

type ParentAddress struct {
	Parent DcrAddress
	name   string
}

// Address refers to folders in dcr lib. Adress is what will be the node in DCR Graph.
type DcrAddress interface {
	Parent
	// Collaboration() *Collaboration
	//IsRoot() bool
	//AddAddress
	//AddTransformationSpecs
	//DeRefTranformationGroup
	//DeRefTransformation
	//DeRefAddressGroup
	//DeRefAddress

	//ListTransformations()
	//ListAddresses()

}
type Address struct {
	Parent
	name string
}

func (a Address) ParentAddress() *ParentAddress {
	log.Infof("parent function for address needs to be implemented")
	return nil
}

func (a Address) Name() string {
	return a.name
}

func NewAddress() DcrAddress {
	return Address{}
}

const (
	COLLABORATION_FOLDER_NAME = "collaboration"
)

const (
	ADDRESS_TYPE_ROOT           = "/"
	ADDRESS_TYPE_SOURCE         = "source/"
	ADDRESS_TYPE_DESTINATION    = "destination/"
	ADDRESS_TYPE_TRANSFORMATION = "transformation/"
)

// All Address types will implement this interface
type Ref string

//SourceRef : /collaborator_name/source/source_table_name
//DestinationRef : /collaborator_name/destination/destination_table_name
//TransformationRef : /collaborator_name/transformation/transformation_group_name
