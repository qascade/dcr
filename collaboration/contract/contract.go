package contract

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/qascade/dcr/collaboration/utils"
	log "github.com/sirupsen/logrus"
)

// Interface that will be implemented by all contract types
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

// Contracter is the interface to which all contracts must implement
type Contracter interface {
	Parse(path string) (*ContractSpec, *[]TablesContractSpec, error)
	Hash() string   // Hash the yamls
}

// Validator will validate yamls
type Validator interface {
	Validate() error // Validate the yamls
}

type Contract struct {
	Name          string
	Version       string
	Purpose       string
	collabPkgPath    fs.FS
	Collaborators map[string]Collaborator
}

func (c *Contract) Parse(path string) (*ContractSpec, *[]TablesContractSpec, error) {
	c.collabPkgPath = os.DirFS(path)
	var (
		cSpec    Spec
		tSpec    Spec
		cResult  ContractSpec
		tResults []TablesContractSpec
	)
	var tablesRE = regexp.MustCompile(`.*_tables\.yaml`)
	var contractRE = regexp.MustCompile(`contract\.yaml`)
	err := fs.WalkDir(c.collabPkgPath, ".", func(relpath string, d fs.DirEntry, err error) error {
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
			cSpec, err = ParseSpec(contractB, ContractSpecType)
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
			tSpec, err = ParseSpec(testContractB, TablesContractSpecType)
			if err != nil {
				return fmt.Errorf("error while reading the file: %w", err)
			}
			tBytes, err := yaml.Marshal(tSpec)
			if err != nil {
				return fmt.Errorf("error while marshalling the file: %w", err)
			}
			var tResult TablesContractSpec
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

func NewContract(spec *ContractSpec, tcSpecs *[]TablesContractSpec, cPkgPath string) (*Contract, error) {
	collaborators := make(map[string]Collaborator)
	for _, c := range spec.Collaborators {
		userAgents := make(map[string]UserAgent)
		for _, u := range c.UserAgents {
			userAgents[u.Name] = UserAgent{
				Name:  u.Name,
				Email: u.Email,
			}
		}
		collaborators[c.Name] = Collaborator{
			Name:       c.Name,
			gitRepo:    c.GitRepo,
			UserAgents: userAgents,
			warehouse:  c.Warehouse,
			tableRegister: NewTableRegister(tcSpecs, c.Name),
		}
	}
	contract := &Contract{
		Name:          spec.Name,
		Version:       spec.Version,
		Purpose:       spec.Purpose,
		collabPkgPath: os.DirFS(cPkgPath),
		Collaborators: collaborators,
	}
	return contract, nil
}

func (c *Contract) Validate() error {
	log.Debug("Contract Validation not Implemented")
	return nil
}

func (c *Contract) Hash() (hashstr string) {
	hashStr := ""
	hashStr += utils.HashString(c.Name)
	hashStr += utils.HashString(c.Version)
	hashStr += utils.HashString(c.Purpose)
	for _, collaborator := range c.Collaborators {
		hashStr += collaborator.Hash()
	}
	hashStr = utils.HashString(hashStr)
	return hashStr
}

type Collaborator struct {
	Name       string
	gitRepo    string
	UserAgents map[string]UserAgent
	// TODO - warehouse to be changed to source.SourceWarehouse later
	warehouse WarehouseType
	tableRegister *TableRegister
}

func (c *Collaborator) Hash() (hashstr string) {
	hashStr := ""
	hashStr += utils.HashString(c.Name)
	hashStr += utils.HashString(c.gitRepo)
	for _, userAgent := range c.UserAgents {
		hashStr += userAgent.Hash()
	}
	hashStr += utils.HashString(string(c.warehouse))
	hashStr = utils.HashString(hashStr)
	return hashStr
}
// Function to create instance of Contract from ContractSpec

func ParseSpec(fileYaml []byte, specType SpecType) (Spec, error) {
	var bs Spec
	if specType == ContractSpecType {
		bs = ContractSpec{}
	} else {
		bs = TablesContractSpec{}
	}

	err := utils.UnmarshalStrict(fileYaml, &bs)
	if err != nil {
		var bs2 Spec
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

type UserAgent struct {
	Name  string
	Email string
}

func (u *UserAgent) Hash() (hashstr string) {
	hashStr := ""
	hashStr += utils.HashString(u.Name)
	hashStr += utils.HashString(u.Email)
	hashStr = utils.HashString(hashStr)
	return hashStr
}
