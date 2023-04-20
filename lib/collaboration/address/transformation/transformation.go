// This package will contain transformation types
package transformation

import (
	"github.com/flosch/pongo2/v6"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/qascade/dcr/lib/collaboration/config"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type TransformationType string

type DestinationAllowed interface {
	Transformation
}

type Transformation interface {
	Compile(template pongo2.Template) string
}

// All Supported Differentially Private Queries must comply to this interface.
// We may also define some DPQuery templates.
// But we also need to give support for general DP Query.

type GenericPrivateQuery struct {
	CollaboratorName string
	Query            string
	template         string
	tables           []string
	templateInputs   map[string]string
}

func NewGenericPrivateQuery(cName string, tSpec config.TransformationSpec) Transformation {
	return &GenericPrivateQuery{
		CollaboratorName: cName,
		template:         tSpec.Template,
		tables:           extractTableRefs(tSpec.From),
		templateInputs:   createTemplateInputs(tSpec.NoiseParams...),
	}
}

func (g *GenericPrivateQuery) Compile(template pongo2.Template) string {
	log.Info("Compile query for generic private query yet to be implemented.")
	return ""
}

func extractTableRefs(from []config.FromSpec) []string {
	var tables []string
	for _, table := range from {
		tables = append(tables, table.Ref)
	}
	return tables
}

func createTemplateInputs(options ...string) map[string]string {
	templateInputs := make(map[string]string)
	for _, option := range options {
		templateInputs[option] = option
	}
	return templateInputs
}
