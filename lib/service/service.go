package service

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/qascade/dcr/lib/collaboration"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Service struct {
	collaboration       *collaboration.Collaboration
	orderedCollabEvents []collaboration.Event
	eventStatus         map[collaboration.Event]collaboration.EventStatus
}

func NewService(pkgPath string) (*Service, error) {
	collab, err := collaboration.NewCollaboration(pkgPath)
	if err != nil {
		err = fmt.Errorf("err creating new collaboration with package path: %s", pkgPath)
		log.Error(err)
		return nil, err
	}

	runnableEvents, err := collab.GetOrderedRunnableEvents()
	if err != nil {
		err = fmt.Errorf("err getting runnable events: %s", err)
		log.Error(err)
		return nil, err
	}

	eventStatus := make(map[collaboration.Event]collaboration.EventStatus)
	service := &Service{
		collaboration:       collab,
		orderedCollabEvents: runnableEvents,
		eventStatus:         eventStatus,
	}
	return service, nil
}

func (s *Service) Run() {
	for _, event := range s.orderedCollabEvents {
		event.Run()
	}

}

func (s *Service) Runner() {

}
