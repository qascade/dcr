package collaboration

import (
	"fmt"
	"os"

	"github.com/qascade/dcr/collaboration/contract"
	log "github.com/sirupsen/logrus"
)

// Collaboration Package structure
// should contain three files
// 1. contract.yaml
// 2. two *_tables.yaml for each collaborator mentioned in contract.yaml
// This function for now, is to validate names across over the contract package

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type CollaborationPackage struct {
	Contract	*contract.Contract
}

// A Type of Collaboration Package must implement Collaboration interface 
type Collaboration interface {
	Verify(path string) error // Verify checks Both Contracts from Git Repos
	Terminate() error // Terminate Clean Room Service Terminates 
}

func NewCollaborationPkg(path string) (*CollaborationPackage, error) {
	var contracter contract.Contracter
	cSpec, tSpecs, err := contracter.Parse(path)
	if err != nil {
		log.Error("error while parsing the contract: %w", err)
		return nil, fmt.Errorf("error while parsing the contract: %w", err)
	}

	contrct, err := contract.NewContract(cSpec, tSpecs, path)
	if err != nil {
		log.Error("error while creating the contract: %w", err)
		return nil, fmt.Errorf("error while creating the contract: %w", err)
	}

	collabPkg := &CollaborationPackage{
		Contract: contrct,
	}
	return collabPkg, nil
}

func (c *CollaborationPackage) Verify(path string) error {
	log.Debug("Verify yet to be implemented")
	return nil
}