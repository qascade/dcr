package address

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// AddressRef structure:
// /<collaboration_id>/<collaborator_name>/<address_type>/<address_name>
// <collaboration_id>: will not come into play until we enter support multiple collaborations at the same time.
// <collaboration_id>: <collaboration_name>_HASH
type AddressRef string

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

// This function will take in the name of the address and the name of the collaborator and return the absolute addressRef
func Abs(addressName string, collaboratorName string, addType AddressType) AddressRef {
	return AddressRef(ADDRESS_TYPE_ROOT) + AddressRef(collaboratorName) + AddressRef(addType) + AddressRef(addressName)
}

func NewCollaboratorRef(collaboratorName string) AddressRef {
	return AddressRef(ADDRESS_TYPE_ROOT) + AddressRef(collaboratorName)
}
