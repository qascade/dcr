package collaboration

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/qascade/dcr/lib/collaboration/address"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func (c *Collaboration) DeRefSource(ref address.AddressRef) (address.DcrAddress, error) {
	if add, ok := c.AddressGraph.CachedSources[ref]; ok {
		return add, nil
	}
	return nil, fmt.Errorf("source address not found. %s", ref)
}

func (c *Collaboration) DeRefTransformation(ref address.AddressRef) (address.DcrAddress, error) {
	if add, ok := c.AddressGraph.CachedTransformations[ref]; ok {
		return add, nil
	}
	return nil, fmt.Errorf("transformation address not found. %s", ref)
}

func (c *Collaboration) DeRefDestination(ref address.AddressRef) (address.DcrAddress, error) {
	if add, ok := c.AddressGraph.CachedDestinations[ref]; ok {
		return add, nil
	}
	return nil, fmt.Errorf("transformation address not found. %s", ref)
}

// This function returns the path, where to put the destination result.
func (c *Collaboration) GetOutputPath(destOwner address.AddressRef) (string, error) {
	owner := string(destOwner)
	owner = owner[1:]
	pkgInfo, ok := c.collaborationConfig.PackagesInfo[owner]
	if !ok {
		err := fmt.Errorf("collaborator with name: %s does not exist", owner)
		log.Error(err)
		return "", err
	}
	return pkgInfo.PkgPath, nil
}

// This is a helper function for the unique email specific example. To be removed later.
func filterResults(output string) string {
	s := strings.Split(output, " ")
	n := len(s)
	return fmt.Sprintf("NonPrivateCount:%s PrivateCount:%s", strings.TrimLeft(s[n-2], "...\n"), strings.Trim(s[n-1], "\n"))
}
