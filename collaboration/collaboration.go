package collaboration

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/qascade/dcr/collaboration/contract"
	"github.com/qascade/dcr/collaboration/utils"
)

// Collaboration Package structure
// should contain three files
// 1. contract.yaml
// 2. two *_tables.yaml for each collaborator mentioned in contract.yaml
// This function for now, is to validate names across over the contract package

type CollaborationPackage struct {
	collabPkg    fs.FS
	repoName     string
	ContractSpec *contract.ContractSpec
}

// A Type of Collaboration Package must implement CollaborationParser interface
type Collaboration interface {
	Parse(path string) (*contract.ContractSpec, *[]contract.TablesContractSpec, error)
	Verify(path string) error
	UploadToRepo(linkToContractFile string) error
	Terminate() error
	GetContractSpec() *contract.ContractSpec
	GetContractFromRepo(link string) (*string, *error)
	DeleteRepo() error
}

func NewCollaborationPkg(path string) (*CollaborationPackage, error) {
	collabPkg := &CollaborationPackage{}
	cSpec, _, err := collabPkg.Parse(path)
	if err != nil {
		return nil, err
	}
	for _, collaborator := range cSpec.Collaborators {
		if collaborator.Name == cSpec.Name {
			collabPkg.repoName = collaborator.GitRepo
			break
		}
	}
	collabPkg.ContractSpec = cSpec
	return collabPkg, nil
}

func (c *CollaborationPackage) GetContractSpec() *contract.ContractSpec {
	return c.ContractSpec
}

func (c *CollaborationPackage) Parse(path string) (*contract.ContractSpec, *[]contract.TablesContractSpec, error) {
	c.collabPkg = os.DirFS(path)
	var (
		cSpec    contract.Spec
		tSpec    contract.Spec
		cResult  contract.ContractSpec
		tResults []contract.TablesContractSpec
	)
	var tablesRE = regexp.MustCompile(`.*_tables\.yaml`)
	var contractRE = regexp.MustCompile(`contract\.yaml`)
	err := fs.WalkDir(c.collabPkg, ".", func(relpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if contractRE.MatchString(d.Name()) {
			contractB, err := os.ReadFile(path + "/" + d.Name())
			if err != nil {
				return fmt.Errorf("error while reading the file: %w", err)
			}
			cSpec, err = ParseSpec(contractB, contract.ContractSpecType)
			if err != nil {
				return fmt.Errorf("error while reading the file: %w", err)
			}
			cBytes, err := yaml.Marshal(cSpec)
			if err != nil {
				return fmt.Errorf("error while marshalling the file: %w", err)
			}
			err = yaml.Unmarshal(cBytes, &cResult)
			if err != nil {
				return fmt.Errorf("unable to unmarshal to TablesContractSpec, %s", d.Name())
			}
		} else if tablesRE.MatchString(d.Name()) {
			testContractB, err := os.ReadFile(path + "/" + d.Name())
			if err != nil {
				return fmt.Errorf("error while reading the file: %w", err)
			}
			tSpec, err = ParseSpec(testContractB, contract.TablesContractSpecType)
			if err != nil {
				return fmt.Errorf("error while reading the file: %w", err)
			}
			tBytes, err := yaml.Marshal(tSpec)
			if err != nil {
				return fmt.Errorf("error while marshalling the file: %w", err)
			}
			var tResult contract.TablesContractSpec
			err = yaml.Unmarshal(tBytes, &tResult)
			if err != nil {
				return fmt.Errorf("unable to unmarshal to TablesContractSpec, %s", d.Name())
			}
			tResults = append(tResults, tResult)
		} else {
			return fmt.Errorf("invalid file, file type with type name not supported, %s", d.Name())
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while walking the directory: %w", err)
	}
	return &cResult, &tResults, nil
}

func ParseSpec(fileYaml []byte, specType contract.SpecType) (contract.Spec, error) {
	var bs contract.Spec
	if specType == contract.ContractSpecType {
		bs = contract.ContractSpec{}
	} else {
		bs = contract.TablesContractSpec{}
	}

	err := utils.UnmarshalStrict(fileYaml, &bs)
	if err != nil {
		var bs2 contract.Spec
		err2 := yaml.Unmarshal(fileYaml, &bs2)
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

func (c* CollaborationPackage) Terminate() error {
	err := c.DeleteRepo()
	if err != nil {
		return err
	}
	log.Println("Collaboration Terminate")
	return err
}