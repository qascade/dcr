// This package will contain transformation types
package transformation

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/qascade/dcr/lib/collaboration/config"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type Transformation interface {
	GetSourcesInfo() []SourceMetadata
	GetPongoInputs() map[string]string
	AppLocation() string
}

type SourceMetadata struct {
	CollaboratorName   string
	SourceName         string
	LocationPongoInput string
	AddressRef         string
}

// A go binary code that takes lists csv's as input and outputs a list of csv's
type GoApp struct {
	CollaboratorName string
	pongoInputs      map[string]string
	appLocation      string
	sourcesInfo      []SourceMetadata
}

func NewGoApp(cName string, tSpec config.TransformationSpec) Transformation {
	pongoInputs := make(map[string]string)
	pongoInputs["uniqueID"] = tSpec.UniqueId
	sources := getSourcesFromSpec(cName, tSpec)
	registerNoiseParams(tSpec, pongoInputs, sources)
	return &GoApp{
		CollaboratorName: cName,
		pongoInputs:      pongoInputs,
		appLocation:      tSpec.AppLocation,
		sourcesInfo:      sources,
	}
}

func (ga *GoApp) GetSourcesInfo() []SourceMetadata {
	return ga.sourcesInfo
}

func (ga *GoApp) GetPongoInputs() map[string]string {
	return ga.pongoInputs
}

func (ga *GoApp) AppLocation() string {
	return ga.appLocation
}
func getSourcesFromSpec(cName string, tSpec config.TransformationSpec) []SourceMetadata {
	var sources []SourceMetadata
	for _, source := range tSpec.From {
		metadata := SourceMetadata{
			CollaboratorName:   cName,
			SourceName:         source.Name,
			LocationPongoInput: source.LocationTag,
			AddressRef:         source.Ref,
		}
		sources = append(sources, metadata)
	}
	return sources
}

func registerNoiseParams(tSpec config.TransformationSpec, pongoInputs map[string]string, sources []SourceMetadata) {
	for _, noiseParam := range tSpec.NoiseParams {
		// These will be populated by transformation runner.
		pongoInputs[noiseParam] = ""
	}
	for _, source := range sources {
		pongoInputs[source.LocationPongoInput] = ""
	}
}
