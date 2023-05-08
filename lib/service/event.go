package service

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/qascade/dcr/lib/collaboration"
	"github.com/qascade/dcr/lib/collaboration/address"
	"github.com/qascade/dcr/lib/utils"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type EventType string

var (
	RUN_TRANSFORMATION_EVENT_TYPE EventType = "/transformation/run"
	SEND_DESTINATION_EVENT_TYPE   EventType = "/destination/send"
)

// These are status values.
const (
	READY     = iota // Event Completed and results are stored.
	NOT_READY        // Yet to be computed
)

type Event interface {
	Run() (string, error)
	Status() EventStatus
	Type() EventType
	AddressRef() address.AddressRef
}

// This function returns the list ordered runnable events with the event increasing graph depth.
// These events are yet to be authorized and are to be done by Authorizer when triggered by Service.
// Runnable events are already authorized. All the unauthorized addresses and their corresponding dependent addresses should not show up in the topo Order.
func GetOrderedRunnableEvents(collab *collaboration.Collaboration) ([]Event, error) {
	runnableRefs, err := collab.AddressGraph.GetOrderedRunnableRefs()
	if err != nil {
		err := fmt.Errorf("err while getting ordered runnable refs: %s", err)
		log.Error(err)
		return nil, err
	}

	events := make([]Event, 0)
	for _, ref := range runnableRefs {
		if ref.IsDestination() {
			dAddI, ok := collab.AddressGraph.CachedDestinations[ref]
			if !ok {
				err = fmt.Errorf("err while getting cached destination: %s", ref)
				log.Error(err)
				return nil, err
			}
			parentTRef := dAddI.(*address.DestinationAddress).Destination.GetTransformationRef()
			dEvent, err := NewSendDestinationEvent(collab, ref, address.AddressRef(parentTRef))
			if err != nil {
				err = fmt.Errorf("err creating new destination event: %s", err)
				log.Error(err)
				return nil, err
			}
			events = append(events, dEvent)

		}
		if ref.IsTransformation() {
			tEvent, err := NewRunTransformationEvent(collab, ref)
			if err != nil {
				err = fmt.Errorf("err creating new transformation event: %s", err)
				log.Error(err)
				return nil, err
			}
			events = append(events, tEvent)
		}
	}
	return events, nil
}

type EventStatus struct {
	statusType int
	ErrorMsg   string
}

// Transformation event is an event that runs a transformation. It is to be computed if Destination, is triggered.
type TransformationEvent struct {
	ref           address.AddressRef
	eventType     EventType
	goAppLocation string
	Result        string
	status        EventStatus
}

func NewRunTransformationEvent(collab *collaboration.Collaboration, ref address.AddressRef) (Event, error) {
	// TODO- Make this generic, maybe compiled transformation??
	goAppLocation, err := collab.CompileTransformation(ref)
	if err != nil {
		err = fmt.Errorf("err compiling transformation: %s", err)
		log.Error(err)
		return nil, err
	}
	return &TransformationEvent{
		ref:           ref,
		eventType:     RUN_TRANSFORMATION_EVENT_TYPE,
		goAppLocation: goAppLocation,
		status: EventStatus{
			statusType: NOT_READY,
		},
	}, nil
}

//go:embed temp_enclave.json
var newEnclaveContent string

func (te *TransformationEvent) Run() (string, error) {
	err := os.Chdir(te.goAppLocation)
	if err != nil {
		err = fmt.Errorf("couldn't change directory path to %s", te.goAppLocation)
		log.Error(err)
		return "", err
	}
	buildCmd := exec.Command("ego-go", "build", "main.go")
	_, err = utils.RunCmd(buildCmd)
	if err != nil {
		return "", err
	}

	signCmd := exec.Command("ego", "sign", "main")
	_, err = utils.RunCmd(signCmd)
	if err != nil {
		return "", err
	}

	// Put harcoded csv names to enclave.json
	oldEnclave := "./enclave.json"
	err = utils.Remove(oldEnclave)
	if err != nil {
		return "", err
	}

	err = utils.WriteStringToFile("./enclave.json", newEnclaveContent)
	if err != nil {
		return "", err
	}

	// Set Simulation Mode by Default
	err = os.Setenv("OE_SIMULATION", "1")
	if err != nil {
		err = fmt.Errorf("unable to set env variable %s", "OE_SIMULATION")
		log.Error(err)
		return "", err
	}

	mainRunCmd := exec.Command("ego", "run", "main")
	output, err := utils.RunCmd(mainRunCmd)
	if err != nil {
		return "", err
	}
	output = filterResults(output)
	te.Result = output
	return output, nil
}

func (te *TransformationEvent) Status() EventStatus {
	return te.status
}

func (te *TransformationEvent) Type() EventType {
	return te.eventType
}

func (te *TransformationEvent) AddressRef() address.AddressRef {
	return te.ref
}

type DestinationEvent struct {
	ref                     address.AddressRef
	eventType               EventType
	status                  EventStatus
	parentTransformationRef address.AddressRef
	OutputLocation          string
	ResultFetcher           ResultFetcher
}

func NewSendDestinationEvent(collab *collaboration.Collaboration, ref address.AddressRef, parentTRef address.AddressRef) (Event, error) {
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

	outputLocation, err := collab.GetOutputPath(destAdd.Owner)
	if err != nil {
		err = fmt.Errorf("err getting output path: %s", err)
		log.Error(err)
		return nil, err
	}

	destEvent := &DestinationEvent{
		status: EventStatus{
			statusType: NOT_READY,
		},
		ref:                     ref,
		parentTransformationRef: parentTRef,
		OutputLocation:          outputLocation,
		eventType:               SEND_DESTINATION_EVENT_TYPE,
	}
	return destEvent, nil
}

func (de *DestinationEvent) Run() (string, error) {
	outputPath := de.OutputLocation
	outputPath = outputPath + "/results.txt"

	output, err := de.ResultFetcher.FetchResult(de.parentTransformationRef)
	if err != nil {
		return "", err
	}
	output = filterResults(output)
	utils.WriteStringToFile(outputPath, output)
	return "", nil
}

// This function returns the status of the destination event
func (de *DestinationEvent) Status() EventStatus {
	return de.status
}

func (de *DestinationEvent) Type() EventType {
	return de.eventType
}

func (de *DestinationEvent) AddressRef() address.AddressRef {
	return de.ref
}

// This is a helper function for the unique email specific example. To be removed later.
func filterResults(output string) string {
	s := strings.Split(output, " ")
	n := len(s)
	return fmt.Sprintf("NonPrivateCount:%s PrivateCount:%s", strings.TrimLeft(s[n-2], "...\n"), strings.Trim(s[n-1], "\n"))
}
