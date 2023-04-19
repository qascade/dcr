package address

const (
	COLLABORATION_FOLDER_NAME = "collaboration"
)

type AddressType string

const (
	ADDRESS_TYPE_ROOT           AddressType = "/"
	ADDRESS_TYPE_SOURCE         AddressType = "source/"
	ADDRESS_TYPE_DESTINATION    AddressType = "destination/"
	ADDRESS_TYPE_TRANSFORMATION AddressType = "transformation/"
)

// All Address types will implement this interface
//SourceRef : /collaborator_name/source/source_table_name
//DestinationRef : /collaborator_name/destination/destination_table_name
//TransformationRef : /collaborator_name/transformation/transformation_group_name

type Ref string
