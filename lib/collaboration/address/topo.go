package address

import (
	"fmt"
	"os"

	"github.com/edwingeng/deque"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type TopoOrder struct {
	List         []AddressRef
	IndegreeList map[AddressRef]int
}

func NewTopoOrder(adjList map[AddressRef][]AddressRef) *TopoOrder {
	indegreeList := make(map[AddressRef]int)
	// Create indegreeList
	for v, neighbourList := range adjList {
		if _, ok := indegreeList[v]; !ok {
			indegreeList[v] = 0
		}
		for _, neighbour := range neighbourList {
			if _, ok := indegreeList[neighbour]; !ok {
				indegreeList[neighbour] = 0
			}
			indegreeList[neighbour]++
		}
	}
	return &TopoOrder{
		List:         make([]AddressRef, 0),
		IndegreeList: indegreeList,
	}
}

func (g *Graph) GetOrderedRunnableRefs() ([]AddressRef, error) {
	topoOrder := NewTopoOrder(g.AdjacencyList)
	return topoOrder.AuthorizedSort(g)
}

// A Destination is the ultimate requester of authorization.
// If a transformation is requesting authorization, it needs to have a list of associated parent destinations along with it.
// Otherwise, its an error.
func (g *Graph) AuthorizeAddress(root AddressRef) (bool, error) {
	if root.IsSource() {
		return true, nil
	}
	visited := make(map[AddressRef]bool)
	parents := make([]AddressRef, 0)
	parents = append(parents, root)
	for _, neighbour := range g.AdjacencyList[root] {
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		movementPermission, err := g.Authorizer(neighbour, parents, visited)
		if !movementPermission {
			return false, err
		}
	}
	g.AuthorityStatus[root] = true
	return true, nil
}

// Helper function for AuthorizeAddress.
func (g *Graph) Authorizer(root AddressRef, parents []AddressRef, visited map[AddressRef]bool) (bool, error) {
	var err error
	var movementPermission bool
	if root.IsSource() {
		sAddress, ok := g.CachedSources[root]
		if !ok {
			return false, fmt.Errorf("source address with given address ref not found. %s", root)
		}
		log.Infof("Authorizing root %s for source address. %s", root, sAddress.(*SourceAddress).Ref)
		return sAddress.Authorize(parents, root)
	}
	for _, neighbour := range g.AdjacencyList[root] {
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		if neighbour.IsTransformation() {
			neighbourAddress, ok := g.CachedTransformations[neighbour]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing transformation %s for address. %s", root, neighbour)
			movementPermission, err = neighbourAddress.Authorize(parents, root)
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else if neighbour.IsDestination() {
			neighbourAddress, ok := g.CachedDestinations[neighbour]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing root %s for destination address. %s", root, neighbour)
			movementPermission, err = neighbourAddress.Authorize(parents, root)
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else if neighbour.IsSource() {
			neighbourAddress, ok := g.CachedSources[neighbour]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing root %s for source address. %s", root, neighbour)
			movementPermission, err = neighbourAddress.Authorize(parents, root)
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else {
			err := fmt.Errorf("invalid address type. %s", root)
			log.Error(err)
			return false, err
		}
		if !movementPermission {
			return false, err
		}
		//parents = append(parents, neighbour)
		//movementPermission, err = g.Authorizer(neighbour, parents, visited)
		//g.AuthorityStatus[neighbour] = movementPermission
	}
	return movementPermission, err
}

// This function returns the topological order of all the addresses that are runnable in the current graph.
func (t *TopoOrder) AuthorizedSort(g *Graph) ([]AddressRef, error) {
	// Create a queue and enqueue all vertices with indegree 0
	var queue deque.Deque = deque.NewDeque()
	for k, v := range t.IndegreeList {
		isAuthorized, err := g.AuthorizeAddress(k)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if v == 0 && isAuthorized {
			queue.PushBack(k)
		}
	}
	var topoOrder []AddressRef

	for queue.Len() != 0 {
		// Dequeue a vertex from queue and add it to topoOrder
		v, ok := queue.PopFront().(AddressRef)
		if !ok {
			err := fmt.Errorf("could not convert to AddressRef, %v", v)
			log.Error(err)
			return nil, err
		}
		topoOrder = append(topoOrder, v)
		// Iterate through all its neighbouring nodes of dequeued node u and decrease their in-degree by 1
		for _, neighbour := range g.AdjacencyList[v] {
			t.IndegreeList[neighbour]--
			// If in-degree becomes zero, add it to queue
			if t.IndegreeList[neighbour] == 0 {
				queue.PushBack(neighbour)
			}
		}
	}
	// Check if there was a cycle
	if len(topoOrder) != g.Count {
		err := fmt.Errorf("there exists a cycle in the graph")
		log.Error(err)
		return nil, err
	}
	return topoOrder, nil
}
