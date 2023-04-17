package address

type Parent interface {
	Parent() *ParentAddress
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
	IsLevelRoot() bool
	//AddAddress
	//AddTransformationSpecs
	//DeRefTranformationGroup
	//DeRefTransformation
	//DeRefAddressGroup
	//DeRefAddress

	//ListTransformations()
	//ListAddresses()

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
