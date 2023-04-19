// Orchestrator will be a server that will create a collaboration session
// It will take requests to run collaboration events.
// A collaboration event can be of type:
// 		1. Run Transformation.
// 		2. Download Destination

package orchestrator

import (
	"context"
	"github.com/qascade/dcr/lib/collaboration"
)

type Service interface {
}

type CleanRoomService struct {
	SessionId            string
	CollaborationPtr     *collaboration.Collaboration
	CollaborationContext *context.Context
}
