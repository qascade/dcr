package service

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/qascade/dcr/lib/collaboration"
	"github.com/qascade/dcr/lib/collaboration/address"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Service struct {
	collaboration       *collaboration.Collaboration
	orderedCollabEvents []Event
	ResultStore         map[address.AddressRef]string
	eventStatus         map[Event]EventStatus
}

type ResultFetcher interface {
	FetchResult(address.AddressRef) (string, error)
}

func NewService(pkgPath string) (*Service, error) {
	collab, err := collaboration.NewCollaboration(pkgPath)
	if err != nil {
		err = fmt.Errorf("err creating new collaboration with package path: %s", pkgPath)
		log.Error(err)
		return nil, err
	}

	runnableEvents, err := GetOrderedRunnableEvents(collab)
	if err != nil {
		err = fmt.Errorf("err getting runnable events: %s", err)
		log.Error(err)
		return nil, err
	}

	eventStatus := make(map[Event]EventStatus)
	resultStore := make(map[address.AddressRef]string)
	service := &Service{
		collaboration:       collab,
		orderedCollabEvents: runnableEvents,
		ResultStore:         resultStore,
		eventStatus:         eventStatus,
	}
	return service, nil
}

func (s *Service) Run() error {
	// Every event is already authorized to run
	// Run all transformations and store them in ResultStore.
	for _, event := range s.orderedCollabEvents {
		if event.Type() == RUN_TRANSFORMATION_EVENT_TYPE {
			output, err := event.Run()
			s.ResultStore[event.AddressRef()] = output
			if err != nil {
				err = fmt.Errorf("err running event: %s", err)
				log.Error(err)
				s.eventStatus[event] = EventStatus{
					statusType: NOT_READY,
					ErrorMsg:   err.Error(),
				}
				return err
			}
			s.eventStatus[event] = EventStatus{
				statusType: READY,
				ErrorMsg:   "",
			}
		}
	}

	for _, event := range s.orderedCollabEvents {
		if event.Type() == SEND_DESTINATION_EVENT_TYPE {
			_, err := event.Run()
			if err != nil {
				err = fmt.Errorf("err running event: %s", err)
				log.Error(err)
				s.eventStatus[event] = EventStatus{
					statusType: NOT_READY,
					ErrorMsg:   err.Error(),
				}
				return err
			}
			s.eventStatus[event] = EventStatus{
				statusType: READY,
				ErrorMsg:   "",
			}
		}
	}
	return nil
}

func (s *Service) FetchResult(ref address.AddressRef) (string, error) {
	if val, ok := s.ResultStore[ref]; ok {
		return val, nil
	}
	err := fmt.Errorf("err while fetching result for ref: %s", ref)
	log.Error(err)
	return "", err
}
