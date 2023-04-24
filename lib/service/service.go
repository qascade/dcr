// Service will be a server that will create a collaboration session
// It will take requests to run collaboration events.
// A collaboration event can be of type:
// 		1. Run Transformation.
// 		2. Download Destination

package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/qascade/dcr/lib/collaboration"
	"github.com/qascade/dcr/lib/collaboration/address"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type CollaborationEvent interface {
	Execute() error
}

type RunTransformationEvent struct {
	Collaboration         *collaboration.Collaboration
	transformationAddress address.DcrAddress
}

func NewRunTransformationEvent(collab *collaboration.Collaboration, tAddress address.DcrAddress) CollaborationEvent {
	return &RunTransformationEvent{
		Collaboration:         collab,
		transformationAddress: tAddress,
	}
}

func (rt *RunTransformationEvent) Execute() error {
	path := "../samples/init_collaboration/go_app/private_count.go"
	buildCmd := exec.Command("ego-go", "build", "-o", path)
	err := buildCmd.Run()
	if err != nil {
		log.Error("Error running ego-go")
		return fmt.Errorf("error running ego-go")
	}

	// Need to strip.go from end
	binPath := strings.TrimSuffix(path, filepath.Ext(path))
	signCmd := exec.Command("ego", "sign", binPath)
	err = signCmd.Run()
	if err != nil {
		log.Error("Error signing ego")
		return fmt.Errorf("error signing ego")
	}

	runCmd := exec.Command("ego", "run", binPath)
	err = runCmd.Run()
	if err != nil {
		err = fmt.Errorf("error running %s", binPath)
		log.Error(err)
		return err
	}
	return nil
}

type DownloadDestinationEvent struct {
}
