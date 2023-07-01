package config

// This package will contain all the Config extraction logic.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/qascade/dcr/lib/utils"
)

// Interface that will be implemented by all contract types
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Parser interface {
	Parse(path string) (*CollaborationConfig, error)
}

// CollaborationConfig is the root folder for all the config files and implements Parser.

func NewCollaborationConfig() CollaborationConfig {
	return CollaborationConfig{}
}

// All the Specs associated with a single Collaborator
type CollaborationConfig struct {
	CollaborationFolderPath string
	PackagesInfo            map[string]*PackageConfig
}

type PackageConfig struct {
	CollaboratorName        string
	PkgPath                 string
	OutputFolderPath        string
	SourceSpec              *SourceGroupSpec
	TransformationGroupSpec *TransformationGroupSpec
	DestinationGroupSpec    *DestinationGroupSpec
}

func (c CollaborationConfig) Parse(path string) (*CollaborationConfig, error) {
	log.Infof("Parsing the config folder with path %s", path)
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		err = fmt.Errorf("err unable to get absolute path of %s, %s", path, err)
		log.Error(err)
		return nil, err
	}
	pkgsInfo := make(map[string]*PackageConfig)
	var pkgPaths []string
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.Contains(path, "go_app") {
				return filepath.SkipDir
			}
			pkgPaths = append(pkgPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while walking through the path: %v", err)
	}
	pkgPaths = pkgPaths[1:]
	for _, pkgPath := range pkgPaths {
		var pkgConfig *PackageConfig
		pkgConfig, err = c.newPackageConfig(pkgPath)
		if err != nil {
			return nil, err
		}
		pkgsInfo[pkgConfig.CollaboratorName] = pkgConfig
	}
	collabConfig := &CollaborationConfig{
		CollaborationFolderPath: path,
		PackagesInfo:            pkgsInfo,
	}
	return collabConfig, nil
}

func (c *CollaborationConfig) newPackageConfig(pkgPath string) (*PackageConfig, error) {
	log.Infof("Creating new Package Config for pkgPath %s", pkgPath)
	sSpec, err := c.parseSourceSpec(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("error while parsing the source spec: %v, with pkgPath %s", err, pkgPath)
	}

	tSpec, err := c.parseTransformationSpec(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("error while parsing the transformation spec: %v, with pkgPath %s", err, pkgPath)
	}

	dSpec, err := c.parseDestinationSpec(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("error while parsing the destination spec: %v, with pkgPath %s", err, pkgPath)
	}

	cName, err := getCollaboratorNameFromConfig(sSpec, tSpec, dSpec)
	if err != nil {
		return nil, err
	}

	pkgConfig := &PackageConfig{
		CollaboratorName:        cName,
		PkgPath:                 pkgPath,
		OutputFolderPath:        filepath.Join(pkgPath, "output"),
		SourceSpec:              sSpec,
		TransformationGroupSpec: tSpec,
		DestinationGroupSpec:    dSpec,
	}
	return pkgConfig, nil
}
func (c *CollaborationConfig) parseSourceSpec(path string) (*SourceGroupSpec, error) {
	sourceYamlPath := path + "/sources.yaml"
	sourceSpecB, err := os.ReadFile(sourceYamlPath)
	if err != nil {
		// If the file does not exist, return nil
		return &SourceGroupSpec{}, nil
	}
	sSpec, err := ParseSpec(sourceSpecB, SourceSpecType)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file: %w", err)
	}
	sBytes, err := yaml.Marshal(sSpec)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling the file: %w", err)
	}
	var sResult SourceGroupSpec
	err = yaml.Unmarshal(sBytes, &sResult)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal to SourceGroupSpec, %s", sourceYamlPath)
	}
	for i, sourceSpec := range sResult.Sources {
		sResult.Sources[i].CSVLocation = filepath.Join(path, sourceSpec.CSVLocation)
	}
	return &sResult, nil
}

func (c *CollaborationConfig) parseTransformationSpec(path string) (*TransformationGroupSpec, error) {
	transformationYamlPath := filepath.Join(path, "transformations.yaml")
	transformationSpecB, err := os.ReadFile(transformationYamlPath)
	if err != nil {
		// If the file is not present, then return nil
		return &TransformationGroupSpec{}, nil
	}
	tSpec, err := ParseSpec(transformationSpecB, TransformationSpecType)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file: %w", err)
	}
	tBytes, err := yaml.Marshal(tSpec)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling the file: %w", err)
	}
	var tResult TransformationGroupSpec
	err = utils.UnmarshalStrict(tBytes, &tResult)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal to TransformationGroupSpec %s", transformationYamlPath)
	}
	for i, tSpec := range tResult.Transformations {
		tResult.Transformations[i].AppLocation = resolveRelativeAddress(path, tSpec.AppLocation)
	}
	return &tResult, nil
}

func (c *CollaborationConfig) parseDestinationSpec(path string) (*DestinationGroupSpec, error) {
	destinationYamlPath := filepath.Join(path, "destinations.yaml")
	destinationSpecB, err := os.ReadFile(destinationYamlPath)
	if err != nil {
		// If the file is not present, then return nil
		return &DestinationGroupSpec{}, nil
	}
	dSpec, err := ParseSpec(destinationSpecB, DestinationSpecType)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file: %w", err)
	}
	dBytes, err := yaml.Marshal(dSpec)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling the file: %w", err)
	}
	var dResult DestinationGroupSpec
	err = yaml.Unmarshal(dBytes, &dResult)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal to DestinationGroupSpec, %s", destinationYamlPath)
	}
	for i, dSpec := range dResult.Destinations {
		dResult.Destinations[i].SomeField = resolveRelativeAddress(path, dSpec.SomeField)
		// Modify other fields as necessary
	}
	return &dResult, nil
}

func resolveRelativeAddress(basePath, address string) string {
	if strings.HasPrefix(address, "/") {
		// Absolute address, no need for resolution
		return address
	}
	return filepath.Join(basePath, address)
}

func ParseSpec(yamlBytes []byte, specType SpecType) (Spec, error) {
	var bs Spec
	switch specType {
	case SourceSpecType:
		bs = SourceGroupSpec{}
	case TransformationSpecType:
		bs = TransformationGroupSpec{}
	case DestinationSpecType:
		bs = DestinationGroupSpec{}
	}
	err := utils.UnmarshalStrict(yamlBytes, &bs)
	if err != nil {
		var bs2 Spec
		err2 := yaml.Unmarshal(yamlBytes, &bs2)
		if err2 != nil {
			return bs, fmt.Errorf("error parsing yaml: %w", err2)
		}
		partialSpecYaml, err3 := yaml.Marshal(bs2)
		if err3 != nil {
			return bs, fmt.Errorf("error marshaling partial build spec: %w", err3)
		}
		return bs, fmt.Errorf("error parsing yaml.  Parse result:\n%s\nParse error:%s", partialSpecYaml, err)
	}
	return bs, err
}

func getCollaboratorNameFromConfig(sSpec *SourceGroupSpec, tSpec *TransformationGroupSpec, dSpec *DestinationGroupSpec) (string, error) {
	if sSpec != nil && sSpec.CollaboratorRef != "" {
		return sSpec.CollaboratorRef, nil
	}
	if tSpec != nil && tSpec.CollaboratorRef != "" {
		return tSpec.CollaboratorRef, nil
	}
	if dSpec != nil && dSpec.CollaboratorRef != "" {
		return dSpec.CollaboratorRef, nil
	}
	return "", fmt.Errorf("unable to find the collaborator name in the config")
}
