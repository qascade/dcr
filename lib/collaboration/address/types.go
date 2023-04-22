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

func CacheAddresses(collabConfig config.CollaborationConfig) (map[AddressRef]DcrAddress, map[AddressRef]DcrAddress, map[AddressRef]DcrAddress) {
	cSources := make(map[AddressRef]DcrAddress)
	cTransformations := make(map[AddressRef]DcrAddress)
	cDestinations := make(map[AddressRef]DcrAddress)

	for _, pkgConfig := range collabConfig.PackagesInfo {
		// TODO - Need to create all these address through a AddressFactory
		collaboratorName := pkgConfig.CollaboratorName
		for _, sSpec := range pkgConfig.SourceSpec.Sources {
			s := source.NewLocalSource(collaboratorName, sSpec)
			ref := Abs(sSpec.Name, collaboratorName, ADDRESS_TYPE_SOURCE)
			sAddress := NewSourceAddress(ref, collaboratorName, getAddressRefSlice(sSpec.ConsumersAllowed), getTransformationRefSlice(sSpec.DestinationsAllowed), s)
			cSources[ref] = sAddress
		}

		for _, tSpec := range pkgConfig.TransformationGroupSpec.Transformations {
			t := transformation.NewGoApp(collaboratorName, tSpec)
			ref := Abs(tSpec.Name, collaboratorName, ADDRESS_TYPE_TRANSFORMATION)
			tAddress := NewTransformationAddress(ref, collaboratorName, getAddressRefSlice(tSpec.ConsumerAllowed), getAddressRefSlice(tSpec.DestinationAllowed), t)
			cTransformations[ref] = tAddress
		}

		for _, dSpec := range pkgConfig.DestinationGroupSpec.Destinations {
			d := destination.NewLocalDestination(collaboratorName, dSpec)
			ref := Abs(dSpec.Name, collaboratorName, ADDRESS_TYPE_DESTINATION)
			requester := NewCollaboratorRef(collaboratorName)
			dAddress := NewDestinationAddress(ref, requester, d)
			cDestinations[ref] = dAddress
		}
	}
	return cSources, cTransformations, cDestinations
}
