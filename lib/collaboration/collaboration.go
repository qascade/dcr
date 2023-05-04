package collaboration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
	log "github.com/sirupsen/logrus"

	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/collaboration/address/transformation"
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
	cachedDestinations    map[address.AddressRef]address.DcrAddress
}

func NewCollaboration(pkgPath string) (*Collaboration, error) {
	var collaborators []string

	collabConfig, err := config.Parser(config.NewCollaborationConfig()).Parse(pkgPath)
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
		err = fmt.Errorf("err while topologically sorting the address graph: %s", err)
		log.Error(err)
		return nil, err
	}
	collaboration := &Collaboration{
		AddressGraph:          graph,
		Collaborators:         collaborators,
		collaborationConfig:   *collabConfig,
		cachedSources:         cSources,
		cachedTransformations: cTransformations,
		cachedDestinations:    cDestinations,
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
	if root.IsTransformation() {
		parentTransformationRef = root
	}
	if root.IsDestination() {
		// As movement from destination to Transformation is always allowed. We need to store its neighbouring transformation.
		dAddress, ok := c.cachedDestinations[root].(*address.DestinationAddress)
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
	if root.IsSource() {
		sAddress, ok := c.cachedSources[root]
		if !ok {
			return false, fmt.Errorf("address with given address ref not found. %s", root)
		}
		log.Infof("Authorizing collaborator for source address. %s", root)
		return sAddress.Authorize(collaboratorRef, parentTransformationRef)
	}
	for _, neighbour := range c.AddressGraph.AdjacencyList[root] {
		if visited[neighbour] {
			continue
		}
		visited[neighbour] = true
		if root.IsTransformation() {
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
		} else if root.IsDestination() {
			neighbourAddress, ok := c.cachedDestinations[root]
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
	log.Info("Validating Noises as per trust Group Policy")
	err := validateNoises(sourceInfo)
	if err != nil {
		err = fmt.Errorf("err source noises not compliant to trust group policy, %s", err)
		log.Error(err)
		return "", err
	}
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
				pongoInputs[k] = sAdd.Source.Extract()
			}
		}
	}
	_, err = prepareGoApp(appLocation, pongoInputs)
	if err != nil {
		return "", err
	}
	// TODO- HardCoding outputPath will have to populate later.
	return appLocation, nil
}

// This function validate noises for the members in the trust group.
// A trust group is a set of sources who have given permission to the same transformation.
func validateNoises(sourceInfo []transformation.SourceMetadata) error {
	// This validateNoises will also need all the list of collaborators who gives permission to same destination.
	// This list will be fetched from address graph. All these collaborators will form a Trust Group
	// After this we will have three options for noise Validation/Propagation
	// 1. Only one collaborator from the trust Group is allowed to define noises.
	// 		a. This validation can be simplified in yaml where other callaborators can acknowledge that by refering the noise parameters which can introduced as a address_type.
	// 2. All collaborators that form a trust group have to give same noises at source level. If the noises mismatch, it will result in an error.
	// 3. There is no such thing as a trust group everybody is free to define whatever amount of noise they want. We will have to define a mechanism such that from all the lists of noises.
	//    that contributes the largest noise in the result will end up getting selected.
	log.Info("Noise Validation yet to be implemented")
	return nil

}

func prepareGoApp(appLocation string, pongoInputs map[string]string) (string, error) {
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
