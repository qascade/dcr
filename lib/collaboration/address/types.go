package address

import (
	"github.com/qascade/dcr/lib/collaboration/address/destination"
	"github.com/qascade/dcr/lib/collaboration/address/source"
	"github.com/qascade/dcr/lib/collaboration/address/transformation"
	"github.com/qascade/dcr/lib/collaboration/config"
)

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

func CacheAddresses(collabConfig config.CollaborationConfig) (map[AddressRef]source.Source, map[AddressRef]transformation.Transformation, map[AddressRef]destination.Destination) {
	cSources := make(map[AddressRef]source.Source)
	cTransformations := make(map[AddressRef]transformation.Transformation)
	cDestinations := make(map[AddressRef]destination.Destination)

	for _, pkgConfig := range collabConfig.PackagesInfo {
		// TODO - Need to create all these address through a AddressFactory
		collaboratorName := pkgConfig.CollaboratorName
		for _, sSpec := range pkgConfig.SourceSpec.Sources {
			s := source.NewLocalSource(collaboratorName, sSpec)
			ref := Abs(sSpec.Name, collaboratorName, ADDRESS_TYPE_SOURCE)
			cSources[ref] = s
		}
		for _, tSpec := range pkgConfig.TransformationGroupSpec.Transformations {
			t := transformation.NewGenericPrivateQuery(collaboratorName, tSpec)
			ref := Abs(tSpec.Name, pkgConfig.CollaboratorName, ADDRESS_TYPE_TRANSFORMATION)
			cTransformations[ref] = t
		}
		for _, dSpec := range pkgConfig.DestinationGroupSpec.Destinations {
			d := destination.NewLocalDestination(collaboratorName, dSpec)
			ref := Abs(dSpec.Name, pkgConfig.CollaboratorName, ADDRESS_TYPE_DESTINATION)
			cDestinations[ref] = d
		}
	}
	return cSources, cTransformations, cDestinations
}
