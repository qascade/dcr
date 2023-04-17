package config

// This package will contain all the Config extraction logic.

import (
	"fmt"
	"os"

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
	Parse(path string) (*CollaborationSpec, map[string]PackageConfig, error)
}

// ConfigFolder is the root folder for all the config files and implements ConfigParser.
type ConfigFolder struct{}

// All the Specs associated with a single Collaborator
type PackageConfig struct {
	pkgInfo                 pkgInfo
	SourceSpec              *SourceGroupSpec
	TransformationGroupSpec *TransformationGroupSpec
	DestinationGroupSpec    *DestinationGroupSpec
}

func (c ConfigFolder) Parse(path string) (*CollaborationSpec, map[string]PackageConfig, error) {
	collabConfig := make(map[string]PackageConfig)
	cSpec, err := c.parseCollaborationSpec(path)
	if err != nil {
		return nil, nil, err
	}
	pkgInfos := c.getPkgInfos(cSpec, path)
	for _, pkgInfo := range pkgInfos {
		sSpec, err := c.parseSourceSpec(pkgInfo)
		if err != nil {
			return nil, nil, err
		}
		tSpec, err := c.parseTransformationSpec(pkgInfo)
		if err != nil {
			return nil, nil, err
		}
		dSpec, err := c.parseDestinationSpec(pkgInfo)
		if err != nil {
			return nil, nil, err
		}
		collabConfig[pkgInfo.CollaboratorName] = PackageConfig{
			pkgInfo:                 pkgInfo,
			SourceSpec:              sSpec,
			TransformationGroupSpec: tSpec,
			DestinationGroupSpec:    dSpec,
		}
	}
	return cSpec, collabConfig, nil
}

func (c *ConfigFolder) parseCollaborationSpec(path string) (*CollaborationSpec, error) {
	collaborationYamlPath := path + "/collaboration.yaml"
	collaborationSpecB, err := os.ReadFile(collaborationYamlPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file: %w", err)
	}
	cSpec, err := ParseSpec(collaborationSpecB, CollaborationSpecType)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file: %w", err)
	}
	cBytes, err := yaml.Marshal(cSpec)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling the file: %w", err)
	}
	var cResult CollaborationSpec
	err = yaml.Unmarshal(cBytes, &cResult)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal to Spec, %s", collaborationYamlPath)
	}
	return &cResult, nil
}

func (c *ConfigFolder) parseSourceSpec(pkgI pkgInfo) (*SourceGroupSpec, error) {
	sourceYamlPath := pkgI.PkgPath + "/sources.yaml"
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
		return nil, fmt.Errorf("unable to unmarshal to TablesContractSpec, %s", sourceYamlPath)
	}
	return &sResult, nil
}

func (c *ConfigFolder) parseTransformationSpec(pkgI pkgInfo) (*TransformationGroupSpec, error) {
	transformationYamlPath := pkgI.PkgPath + "/transformations.yaml"
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
		return nil, fmt.Errorf("unable to unmarshal to TablesContractSpec, %s", transformationYamlPath)
	}
	return &tResult, nil
}

func (c *ConfigFolder) parseDestinationSpec(pkgI pkgInfo) (*DestinationGroupSpec, error) {
	destinationYamlPath := pkgI.PkgPath + "/destinations.yaml"
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
		return nil, fmt.Errorf("unable to unmarshal to TablesContractSpec, %s", destinationYamlPath)
	}
	return &dResult, nil
}

type pkgInfo struct {
	CollaboratorName string
	PkgPath          string
}

func (c *ConfigFolder) getPkgInfos(cSpec *CollaborationSpec, pkgDirPath string) []pkgInfo {
	var pkgInfos []pkgInfo
	for _, collaborator := range cSpec.Collaborators {
		path := pkgDirPath + "/" + collaborator.Name + "_pkg"
		pkgInfos = append(pkgInfos, pkgInfo{CollaboratorName: collaborator.Name, PkgPath: path})
	}
	return pkgInfos
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
	case CollaborationSpecType:
		bs = CollaborationSpec{}
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
