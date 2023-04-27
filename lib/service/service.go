package service

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/qascade/dcr/lib/collaboration"
	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/utils"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Service struct {
	collaboration *collaboration.Collaboration
	collabEvent   *CollaborationEvent
}

func NewService(pkgPath string, runner string, destOwner string, tRef string, destRef string) (*Service, error) {
	collab, err := collaboration.NewCollaboration(pkgPath)
	if err != nil {
		err = fmt.Errorf("err creating new collaboration with package path: %s", pkgPath)
		log.Error(err)
		return nil, err
	}

	runnerRef := address.NewCollaboratorRef(runner)
	isRunnerAuthorized, err := collab.AuthorizeCollaborationEvent(runnerRef, address.AddressRef(tRef))
	if err != nil {
		err = fmt.Errorf("err while authorizing collaborator %s with ref %s: %s", runner, tRef, err)
		log.Error(err)
		return nil, err
	}
	if !isRunnerAuthorized {
		err = fmt.Errorf("the collaborator %s is not authorized to run transformation %s", runner, tRef)
		log.Error(err)
		return nil, err
	}

	destOwnerRef := address.NewCollaboratorRef(destOwner)
	isDestinationAuthorized, err := collab.AuthorizeCollaborationEvent(destOwnerRef, address.AddressRef(destRef))
	if err != nil {
		err = fmt.Errorf("err while authorizing collaborator %s with ref %s: %s", destOwner, destRef, err)
		return nil, err
	}
	if !isDestinationAuthorized {
		err = fmt.Errorf("err while authorizing collaborator %s with ref %s: %s", destOwnerRef, destRef, err)
		log.Error(err)
		return nil, err
	}

	collabEvent, err := NewCollaborationEvent(collab, runnerRef, address.AddressRef(tRef), destOwnerRef, address.AddressRef(destRef))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	service := &Service{
		collaboration: collab,
		collabEvent:   collabEvent,
	}
	return service, err
}

func (s *Service) RunCollaborationEvent() error {
	return s.collabEvent.Run()
}

type CollaborationEvent struct {
	Collaboration     *collaboration.Collaboration
	Runner            address.AddressRef
	TransformationRef address.AddressRef
	DestinationOwner  address.AddressRef
	DestinationRef    address.AddressRef
}

func NewCollaborationEvent(collab *collaboration.Collaboration, runner address.AddressRef, tRef address.AddressRef, destOwner address.AddressRef, destRef address.AddressRef) (*CollaborationEvent, error) {
	if !strings.Contains(string(tRef), "transformation") {
		err := fmt.Errorf("ref of Invalid type %s. Should be of type /transformation", tRef)
		log.Error(err)
		return nil, err
	}
	if !strings.Contains(string(destRef), "destination") {
		err := fmt.Errorf("ref of Invalid type %s. Should be of type /destination", destRef)
		log.Error(err)
		return nil, err
	}

	collabEvent := &CollaborationEvent{
		Collaboration:     collab,
		Runner:            runner,
		TransformationRef: tRef,
		DestinationOwner:  destOwner,
		DestinationRef:    destRef,
	}
	return collabEvent, nil
}

//go:embed temp_enclave.json
var newEnclaveContent string

func (ce *CollaborationEvent) Run() error {
	goAppPath, err := ce.Collaboration.CompileTransformation(ce.TransformationRef)
	if err != nil {
		log.Error(err)
		return err
	}

	outputPath, err := ce.Collaboration.GetOutputPath(ce.DestinationOwner)
	if err != nil {
		log.Error(err)
		return err
	}

	err = ce.SendDestination(goAppPath, outputPath)
	if err != nil {
		err = fmt.Errorf("err sending Destination to %s for transformation: %s, %s", ce.DestinationOwner, ce.TransformationRef, err)
		log.Error(err)
		return err
	}

	return nil
}

func (ce *CollaborationEvent) SendDestination(appPath string, outputPath string) error {
	err := os.Chdir(appPath)
	if err != nil {
		err = fmt.Errorf("couldn't change directory path to %s", appPath)
		log.Error(err)
		return err
	}
	buildCmd := exec.Command("ego-go", "build", "main.go")
	_, err = utils.RunCmd(buildCmd)
	if err != nil {
		return err
	}

	signCmd := exec.Command("ego", "sign", "main")
	_, err = utils.RunCmd(signCmd)
	if err != nil {
		return err
	}

	// Put harcoded csv names to enclave.json
	oldEnclave := "./enclave.json"
	err = utils.Remove(oldEnclave)
	if err != nil {
		return err
	}

	err = utils.WriteStringToFile("./enclave.json", newEnclaveContent)
	if err != nil {
		return err
	}

	err = os.Setenv("OE_SIMULATION", "1")
	if err != nil {
		err = fmt.Errorf("unable to set env variable %s", "OE_SIMULATION")
		log.Error(err)
		return err
	}

	mainRunCmd := exec.Command("ego", "run", "main")
	output, err := utils.RunCmd(mainRunCmd)
	if err != nil {
		return err
	}
	//fmt.Println(output)

	outputPath = outputPath + "/results.txt"
	output = filterResults(output)
	utils.WriteStringToFile(outputPath, output)
	return nil
}

func filterResults(output string) string {
	s := strings.Split(output, " ")
	n := len(s)
	return fmt.Sprintf("NonPrivateCount:%s PrivateCount:%s", strings.TrimLeft(s[n-2], "...\n"), strings.Trim(s[n-1], "\n"))
}
