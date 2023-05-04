package address

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

// AddressRef structure:
// /<collaboration_id>/<collaborator_name>/<address_type>/<address_name>
// <collaboration_id>: will not come into play until we enter support multiple collaborations at the same time.
// <collaboration_id>: <collaboration_name>_HASH
type AddressRef string

func (a AddressRef) IsSource() bool {
	return strings.Contains(string(a), string(ADDRESS_TYPE_SOURCE))
}

func (a AddressRef) IsTransformation() bool {
	return strings.Contains(string(a), string(ADDRESS_TYPE_TRANSFORMATION))
}

func (a AddressRef) IsDestination() bool {
	return strings.Contains(string(a), string(ADDRESS_TYPE_DESTINATION))
}

// This function will take in the name of the address and the name of the collaborator and return the absolute addressRef
func Abs(addressName string, collaboratorName string, addType AddressType) AddressRef {
	return AddressRef(ADDRESS_TYPE_ROOT) + AddressRef(collaboratorName) + AddressRef("/") + AddressRef(addType) + AddressRef(addressName)
}

func NewCollaboratorRef(collaboratorName string) AddressRef {
	return AddressRef(ADDRESS_TYPE_ROOT) + AddressRef(collaboratorName)
}
