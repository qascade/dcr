package config

// This package will contain all the Config extraction logic.

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/qascade/dcr/lib/utils"
)

// Interface that will be implemented by all contract types
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type ConfigParser interface {
	Parse(path string) (*CollaborationConfig, error)
}

// ConfigFolder is the root folder for all the config files and implements ConfigParser.
type ConfigFolder struct{}

// All the Specs associated with a single Collaborator
type CollaborationConfig struct {
	CollaborationFolderPath string
	PackagesInfo            map[string]*PackageConfig
}

type PackageConfig struct {
	CollaboratorName        string
	PkgPath                 string
	OutpuFolderPath         string
	SourceSpec              *SourceGroupSpec
	TransformationGroupSpec *TransformationGroupSpec
	DestinationGroupSpec    *DestinationGroupSpec
}

func (c ConfigFolder) Parse(path string) (*CollaborationConfig, error) {
	log.Infof("Parsing the config folder with path %s", path)
	pkgsInfo := make(map[string]*PackageConfig)
	var pkgPaths []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
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

func (c *ConfigFolder) newPackageConfig(pkgPath string) (*PackageConfig, error) {
	sSpec, err := c.parseSourceSpec(pkgPath)
	if err != nil {
		return nil, err
	}
	tSpec, err := c.parseTransformationSpec(pkgPath)
	if err != nil {
		return nil, err
	}
	dSpec, err := c.parseDestinationSpec(pkgPath)
	if err != nil {
		return nil, err
	}

	log.Infof("Validating the CollaboratorRefs for all the specs, %s", sSpec.CollaboratorRef)
	_, err = ValidateAllCollaboratorRefsEqual(sSpec.CollaboratorRef, tSpec.CollaboratorRef, dSpec.CollaboratorRef)
	if err != nil {
		return nil, err
	}

	pkgConfig := &PackageConfig{
		CollaboratorName:        sSpec.CollaboratorRef,
		PkgPath:                 pkgPath,
		OutpuFolderPath:         filepath.Join(pkgPath, "output"),
		SourceSpec:              sSpec,
		TransformationGroupSpec: tSpec,
		DestinationGroupSpec:    dSpec,
	}
	return pkgConfig, nil
}
func (c *ConfigFolder) parseSourceSpec(path string) (*SourceGroupSpec, error) {
	sourceYamlPath := path + "/sources.yaml"
	sourceSpecB, err := os.ReadFile(sourceYamlPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file: %w", err)
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
	for _, sourceSpec := range sResult.Sources {
		sourceSpec.CSVLocation, err = filepath.Abs(sourceSpec.CSVLocation)
		if err != nil {
			return nil, fmt.Errorf("unable to get the absolute path of the csv file: %w", err)
		}
	}
	return &sResult, nil
}

func (c *ConfigFolder) parseTransformationSpec(path string) (*TransformationGroupSpec, error) {
	transformationYamlPath := path + "/transformations.yaml"
	transformationSpecB, err := os.ReadFile(transformationYamlPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file: %w", err)
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
	err = yaml.Unmarshal(tBytes, &tResult)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal to TransformationGroupSpec, %s", transformationYamlPath)
	}
	return &tResult, nil
}

func (c *ConfigFolder) parseDestinationSpec(path string) (*DestinationGroupSpec, error) {
	destinationYamlPath := path + "/destinations.yaml"
	destinationSpecB, err := os.ReadFile(destinationYamlPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file: %w", err)
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
	return &dResult, nil
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
