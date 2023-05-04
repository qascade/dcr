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
	if add, ok := c.cachedSources[ref]; ok {
		return add, nil
	}
	return nil, fmt.Errorf("source address not found. %s", ref)
}

func (c *Collaboration) DeRefTransformation(ref address.AddressRef) (address.DcrAddress, error) {
	if add, ok := c.cachedTransformations[ref]; ok {
		return add, nil
	}
	return nil, fmt.Errorf("transformation address not found. %s", ref)
}

func (c *Collaboration) DeRefDestination(ref address.AddressRef) (address.DcrAddress, error) {
	if add, ok := c.cachedDestinations[ref]; ok {
		return add, nil
	}
	return nil, fmt.Errorf("transformation address not found. %s", ref)
}

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

func (c *Collaboration) GetOrderedRunnableRefs() []address.AddressRef {
	runnableRefs := make([]address.AddressRef, 0)
	graph := *c.AddressGraph
	for i := len(graph.TopoOrder) - 1; i >= 0; i-- {
		if graph.TopoOrder[i].IsTransformation() {
			runnableRefs = append(runnableRefs, graph.TopoOrder[i])
		}
		if graph.TopoOrder[i].IsDestination() {
			runnableRefs = append(runnableRefs, graph.TopoOrder[i])
		}
	}
	return runnableRefs
}

func filterResults(output string) string {
	s := strings.Split(output, " ")
	n := len(s)
	return fmt.Sprintf("NonPrivateCount:%s PrivateCount:%s", strings.TrimLeft(s[n-2], "...\n"), strings.Trim(s[n-1], "\n"))
}
