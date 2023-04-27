package collaboration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2/v6"
	log "github.com/sirupsen/logrus"

	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/collaboration/config"
	"github.com/qascade/dcr/lib/utils"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Collaboration struct {
	AddressGraph          *address.Graph
	Collaborators         []string
	collaborationConfig   config.CollaborationConfig
	cachedSources         map[address.AddressRef]address.DcrAddress
	cachedTransformations map[address.AddressRef]address.DcrAddress
	cachedDestination     map[address.AddressRef]address.DcrAddress
}

func NewCollaboration(pkgPath string) (*Collaboration, error) {
	var collaborators []string

	collabConfig, err := config.ConfigParser(config.NewConfigFolder()).Parse(pkgPath)
	if err != nil {
		err = fmt.Errorf("err parsing collaboration package with package path: %s", pkgPath)
		log.Error(err)
		return nil, err
	}

	for _, pkgConfig := range collabConfig.PackagesInfo {
		collaboratorName := pkgConfig.CollaboratorName
		collaborators = append(collaborators, collaboratorName)
	}
	cSources, cTransformations, cDestinations := address.CacheAddresses(*collabConfig)
	graph, err := address.NewGraph(cSources, cTransformations, cDestinations)
	if err != nil {
		return nil, err
	}
	collaboration := &Collaboration{
		AddressGraph:          graph,
		Collaborators:         collaborators,
		collaborationConfig:   *collabConfig,
		cachedSources:         cSources,
		cachedTransformations: cTransformations,
		cachedDestination:     cDestinations,
	}
	return collaboration, nil
}

func (c *Collaboration) AuthorizeCollaborationEvent(collaboratorRef address.AddressRef, root address.AddressRef) (bool, error) {
	// Authorization is permissible only for transformation and destination addresses.
	// Source addresses are not authorized to access any other address.
	// If collaborator wants to run transformation, it should pass the transformation address as root.
	// If collaborator wants to download a destination it should pass the destination address as root.
	visited := make(map[address.AddressRef]bool)
	var parentTransformationRef address.AddressRef
	if strings.Contains(string(root), "transformation") {
		parentTransformationRef = root
	}
	if strings.Contains(string(root), "destination") {
		// As movement from destination to Transformation is always allowed. We need to store its neighbouring transformation.
		dAddress, ok := c.cachedDestination[root].(*address.DestinationAddress)
		if !ok {
			return false, fmt.Errorf("could not cast address to destination address type: %v", string(root))
		}
		parentTransformationRef = address.AddressRef(dAddress.Destination.GetTransformationRef())
	}
	for _, neighbour := range c.AddressGraph.AdjacencyList[root] {
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		movementPermission, err := c.Authorizer(collaboratorRef, neighbour, visited, parentTransformationRef)
		if !movementPermission {
			return false, err
		}
	}
	return true, nil
}

func (c *Collaboration) Authorizer(collaboratorRef address.AddressRef, root address.AddressRef, visited map[address.AddressRef]bool, parentTransformationRef address.AddressRef) (bool, error) {
	var err error
	var movementPermission bool
	if strings.Contains(string(root), "source") {
		sAddress, ok := c.cachedSources[root]
		if !ok {
			return false, fmt.Errorf("address with given address ref not found. %s", root)
		}
		log.Infof("Authorizing collaborator for source address. %s", root)
		return sAddress.Authorize(collaboratorRef, parentTransformationRef)
	}
	for _, neighbour := range c.AddressGraph.AdjacencyList[root] {
		fmt.Println(c.AddressGraph.AdjacencyList[root])
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		if strings.Contains(string(root), "transformation") {
			neighbourAddress, ok := c.cachedTransformations[root]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing collaborator for transformation address. %s", root)
			movementPermission, err = neighbourAddress.Authorize(collaboratorRef, "")
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else if strings.Contains(string(root), "destination") {
			neighbourAddress, ok := c.cachedDestination[root]
			if !ok {
				return false, fmt.Errorf("address with given address ref not found. %s", root)
			}
			log.Infof("Authorizing collaborator for destination address. %s", root)
			movementPermission, err = neighbourAddress.Authorize(collaboratorRef, "")
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else {
			err = fmt.Errorf("invalid address type. %s", root)
			log.Error(err)
			return false, err
		}
		if !movementPermission {
			return false, err
		}
		movementPermission, err = c.Authorizer(collaboratorRef, neighbour, visited, parentTransformationRef)
	}
	return movementPermission, err
}

func (c *Collaboration) DeRefTransformation(ref address.AddressRef) (*address.DcrAddress, error) {
	if add, ok := c.cachedTransformations[ref]; ok {
		return &add, nil
	}
	return nil, fmt.Errorf("transformation address not found. %s", ref)
}

// Compile Transformation will prepare a go_app package that will return the path for the same, Also the path of the output folder on where to put the results.
func (c *Collaboration) CompileTransformation(tRef address.AddressRef) (string, error) {
	tDcrAdd, ok := c.cachedTransformations[tRef]
	if !ok {
		return "", fmt.Errorf("address with given address ref not found. %s", tRef)
	}

	if tDcrAdd.Type() != address.ADDRESS_TYPE_TRANSFORMATION {
		return "", fmt.Errorf("invalid address type. %s. Should be of type transformation", tDcrAdd)
	}

	tAdd, ok := tDcrAdd.(*address.TransformationAddress)
	if !ok {
		return "", fmt.Errorf("could not cast address to transformation address type: %v", tDcrAdd)
	}
	// TODO - Add Authorizer code here.
	t := tAdd.Transformation
	appLocation := t.AppLocation()
	sourceInfo := t.GetSourcesInfo()
	pongoInputs := t.GetPongoInputs()
	log.Info("Noise Validation yet to be implemented.")
	// Fill rest of the pongo inputs
	for _, source := range sourceInfo {
		sAddI := c.cachedSources[address.AddressRef(source.AddressRef)]
		sAdd := sAddI.(*address.SourceAddress)
		// Fill CSVLocations
		for k := range pongoInputs {
			if pongoInputs[k] == "" {
				noiseParams := sAdd.SourceNoises[tAdd.Ref]
				if _, ok := noiseParams[k]; ok {
					pongoInputs[k] = noiseParams[k]
				}
			}
			if k == source.LocationPongoInput {
				// TODO - fill absolute path only from yaml.
				pongoInputs[k] = sAdd.Source.Extract()
			}
		}
	}
	//pongoInputs["uniqueId"] =
	_, err := prepareGoApp(appLocation, pongoInputs)
	if err != nil {
		return "", err
	}
	// TODO- HardCoding outputPath will have to populate later.
	return appLocation, nil
}

func prepareGoApp(appLocation string, pongoInputs map[string]string) (string, error) {
	// TODO - Add the code to prepare the go app.
	//	tpl, err := pongo2.FromFile(appLocation)
	ctx := pongo2.Context{}
	for k, v := range pongoInputs {
		ctx[k] = v
	}
	// Hardcoding the csv into pongo inputs
	ctx["csvLocation1"] = "./test1.csv"
	ctx["csvLocation2"] = "./test2.csv"

	mainFilePath := filepath.Join(appLocation, "main.tpl")
	tpl, err := pongo2.FromFile(mainFilePath)
	if err != nil {
		return "", err
	}
	output, err := tpl.Execute(ctx)
	if err != nil {
		return "", fmt.Errorf("error while executing the template: %v", err)
	}
	compiledMainPath := filepath.Join(appLocation, "main.go")
	compiledMain, err := os.Create(compiledMainPath)
	if err != nil {
		return "", fmt.Errorf("error while creating the main.go file: %v", err)
	}
	log.Infof("Writing the compiled main.go file to %s", compiledMainPath)

	_, err = compiledMain.WriteString(output)
	if err != nil {
		return "", fmt.Errorf("error while writing the main.go file: %v", err)
	}

	csvLocation1 := pongoInputs["csvLocation1"]
	csvLocation2 := pongoInputs["csvLocation2"]

	// Copying the csv's to the go_app folder.
	newCsV1Path := filepath.Join(appLocation, "test1.csv")
	err = utils.CopyFile(newCsV1Path, csvLocation1)
	if err != nil {
		return "", err
	}

	newCsV2Path := filepath.Join(appLocation, "test2.csv")
	err = utils.CopyFile(newCsV2Path, csvLocation2)
	if err != nil {
		return "", err
	}
	return compiledMainPath, nil
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
