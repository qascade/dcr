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

type ResultStore struct {
	Store map[address.AddressRef]string
}

func NewResultStore() *ResultStore {
	return &ResultStore{
		Store: make(map[address.AddressRef]string),
	}
}

type Service struct {
	collaboration       *collaboration.Collaboration
	orderedCollabEvents []Event
	ResultStore         *ResultStore
	eventStatus         map[Event]EventStatus
}

func NewService(pkgPath string) (*Service, error) {
	collab, err := collaboration.NewCollaboration(pkgPath)
	if err != nil {
		err = fmt.Errorf("err creating new collaboration with package path: %s", pkgPath)
		log.Error(err)
		return nil, err
	}
	resultStore := NewResultStore()
	runnableEvents, err := GetOrderedRunnableEvents(collab, resultStore)
	if err != nil {
		err = fmt.Errorf("err getting runnable events: %s", err)
		log.Error(err)
		return nil, err
	}

	eventStatus := make(map[Event]EventStatus)

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
			err := event.Run()
			//s.ResultStore.Store[event.AddressRef()] = output
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
			err := event.Run()
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
