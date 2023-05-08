package address

import (
	"github.com/qascade/dcr/lib/collaboration/address/destination"
	"github.com/qascade/dcr/lib/collaboration/address/source"
	"github.com/qascade/dcr/lib/collaboration/address/transformation"
	"github.com/qascade/dcr/lib/collaboration/config"
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
			//noiseParams := extractNoiseParams(sSpec)
			ref := Abs(sSpec.Name, collaboratorName, ADDRESS_TYPE_SOURCE)
			sAddress := NewSourceAddress(ref, collaboratorName, getAddressRefSlice(sSpec.ConsumersAllowed), getTransformationRefSlice(sSpec.DestinationsAllowed), s, registerSourceNoises(sSpec.DestinationsAllowed))
			cSources[ref] = sAddress
		}

		for _, tSpec := range pkgConfig.TransformationGroupSpec.Transformations {
			t := transformation.NewGoApp(collaboratorName, tSpec)
			ref := Abs(tSpec.Name, collaboratorName, ADDRESS_TYPE_TRANSFORMATION)
			tAddress := NewTransformationAddress(ref, collaboratorName, getAddressRefSlice(tSpec.ConsumerAllowed), getAddressRefSlice(tSpec.DestinationAllowed), t, tSpec.NoiseParams)
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

func registerSourceNoises(sourceNoises []config.SourceDestinationAllowedSpec) map[AddressRef]map[string]string {
	noiseMap := make(map[AddressRef]map[string]string)
	for _, sourceNoise := range sourceNoises {
		noiseMapList := sourceNoise.NoiseParams
		singleNoiseMap := make(map[string]string)
		for _, noiseMap := range noiseMapList {
			for k, v := range noiseMap {
				singleNoiseMap[k] = v
			}
		}
		noiseMap[AddressRef(sourceNoise.Ref)] = singleNoiseMap
	}
	return noiseMap
}
