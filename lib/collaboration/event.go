package collaboration

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/utils"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

// These are status values.
const (
	AUTHORIZED           = iota // Authorized to run
	NOT_AUTHORIZED              // Not authorized to run
	YET_TO_BE_AUTHORIZED        // Yet to be authorized
	READY                       // Event Completed and results are stored.
	NOT_READY                   // Yet to be computed
)

type Event interface {
	Run() error
	Status() EventStatus
}

type EventStatus struct {
	statusType      int
	AuthorityStatus int
	ErrorMsg        error
}

// Transformation event is an event that runs a transformation. It is to be computed if Destination, is triggered.
type TransformationEvent struct {
	goAppLocation string
	Result        string
	status        EventStatus
}

func NewRunTransformationEvent(collab *Collaboration, ref address.AddressRef) (Event, error) {
	// TODO- Make this generic, maybe compiled transformation??
	goAppLocation, err := collab.CompileTransformation(ref)
	if err != nil {
		err = fmt.Errorf("err compiling transformation: %s", err)
		log.Error(err)
		return nil, err
	}

	return &TransformationEvent{
		goAppLocation: goAppLocation,
		status: EventStatus{
			statusType:      NOT_READY,
			AuthorityStatus: YET_TO_BE_AUTHORIZED,
		},
	}, nil
}

//go:embed temp_enclave.json
var newEnclaveContent string

func (te *TransformationEvent) Run() error {
	// TODO- Need to authorize the event

	err := os.Chdir(te.goAppLocation)
	if err != nil {
		err = fmt.Errorf("couldn't change directory path to %s", te.goAppLocation)
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

	// Set Simulation Mode by Default
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
	te.Result = filterResults(output)
	return nil
}

func (te *TransformationEvent) Status() EventStatus {
	return te.status
}

type DestinationEvent struct {
	status                    EventStatus
	ParentTransformationEvent Event
	OutputLocation            string
}

func NewSendDestinationEvent(collab *Collaboration, ref address.AddressRef) (Event, error) {
	destAddI, err := collab.DeRefDestination(ref)
	if err != nil {
		err = fmt.Errorf("err dereferencing destination: %s", err)
		log.Error(err)
		return nil, err
	}
	destAdd, ok := destAddI.(*address.DestinationAddress)
	if !ok {
		err = fmt.Errorf("err dereferencing destination: %s", err)
		log.Error(err)
		return nil, err
	}

	parentTransformationRef := destAdd.Ref
	outputLocation, err := collab.GetOutputPath(destAdd.Owner)
	if err != nil {
		err = fmt.Errorf("err getting output path: %s", err)
		log.Error(err)
		return nil, err
	}
	parentTransformationEvent, err := NewRunTransformationEvent(collab, parentTransformationRef)
	if err != nil {
		err = fmt.Errorf("err creating new transformation event: %s", err)
		log.Error(err)
		return nil, err
	}

	destEvent := &DestinationEvent{
		ParentTransformationEvent: parentTransformationEvent,
		OutputLocation:            outputLocation,
		status: EventStatus{
			statusType:      NOT_READY,
			AuthorityStatus: YET_TO_BE_AUTHORIZED,
		},
	}
	return destEvent, nil
}

func (de *DestinationEvent) Run() error {
	// TODO- Need to authorize the event
	// Set status accordingly
	outputPath := de.OutputLocation
	outputPath = outputPath + "/results.txt"
	// Need to check authority status of parent transformation event
	// If Parent Already ready
	output := de.ParentTransformationEvent.(*TransformationEvent).Result
	output = filterResults(output)
	utils.WriteStringToFile(outputPath, output)
	return nil
}

// This function returns the status of the destination event
func (de *DestinationEvent) Status() EventStatus {
	return de.ParentTransformationEvent.Status()
}
