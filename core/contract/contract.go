package contract

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// Interface that will be implemented by all contract types
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Contracter interface {
	Validate() error // Validate the yamls
	Verify() error   // Verify through github repo's
}

type Contract struct {
	Name          string
	Version       string
	Purpose       string
	Collaborators map[string]Collaborator
}

type Collaborator struct {
	Name       string
	gitRepo    string
	UserAgents map[string]UserAgent
	// TODO - warehouse to be changed to source.SourceWarehouse later
	warehouse WarehouseType
}

type UserAgent struct {
	Name  string
	Email string
}

// Function to create instance of Contract from ContractSpec
func NewContract(spec ContractSpec) (*Contract, error) {
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
			gitRepo:    c.gitRepo,
			UserAgents: userAgents,
			warehouse:  c.Warehouse,
		}
	}
	contract := &Contract{
		Name:          spec.Name,
		Version:       spec.Version,
		Purpose:       spec.Purpose,
		Collaborators: collaborators,
	}
	return contract, nil
}

// Contract Package structure
// should contain three files
// 1. contract.yaml
// 2. two *_tables.yaml for each collaborator mentioned in contract.yaml
// This function for now, is to validate names across over the contract package

func ParseContract(path string) (*Contract, error) {
	log.Debug("Parse not Implemented")
	var cntract *Contract
	cntract.Validate()
	return nil, nil
}

func (c *Contract) Verify() error {
	// TODO
	log.Debug("Verify not Implemented")
	return nil
}

func (c *Contract) Validate() error {
	log.Debug("Validate not Implemented")
	return nil
}
