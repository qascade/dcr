package config

import (
	"github.com/qascade/dcr/lib/graph"
)

type Spec interface{}

type CollaborationSpec struct {
	Name          string             `yaml:"name"`
	Version       string             `yaml:"version"`
	Purpose       string             `yaml:"purpose,omitempty"`
	ServiceType   string             `yaml:"service_type"`
	Collaborators []CollaboratorSpec `yaml:"collaborators"`
}

type TransformationGroupSpec struct {
	CollaboratorRef string               `yaml:"collaborator"`
	Transformations []TransformationSpec `yaml:"transformations"`
}

type CollaboratorSpec struct {
	Name          string `yaml:"name"`
	GitRepo       string `yaml:"git_repo"` // Name of the git repo.
	CredsLocation string `yaml:"creds"`
}

// For now only supporting Count Query.
// TODO - Make this Spec such that any type of transformation can be parsed using this spec.
type TransformationSpec struct {
	Name            string     `yaml:"name"`
	Count           string     `yaml:"count"`
	From            []FromSpec `yaml:"from"`
	NoiseType       string     `yaml:"noise_type"`
	NoiseParams     []string   `yaml:"noise_parameters"`
	JoinKey         string     `yaml:"join_key"`
	Template        string     `yaml:"template"`
	ConsumerAllowed []string   `yaml:"consumer_allowed"`
}

type SourceGroupSpec struct {
	CollaboratorRef     string       `yaml:"collaborator"`
	Sources             []SourceSpec `yaml:"tables"`
	DestinationsAllowed []DestinationAllowedSpec
}

type SourceSpec struct {
	Name        string `yaml:"name"`
	CSVLocation string `yaml:"csv_location"`
	Description string `yaml:"description"`
	// TODO- Do we need to add addressRef here?
	ConsumersAllowed []string     `yaml:"consumers_allowed"`
	Columns          []ColumnSpec `yaml:"columns"`
}

type ColumnSpec struct {
	Name              string   `yaml:"name"`
	Type              string   `yaml:"type"`
	MaskingType       string   `yaml:"masking_type"`
	Selectable        bool     `yaml:"selectable"`
	AggregatesAllowed []string `yaml:"aggregates_allowed"`
	JoinKey           bool     `yaml:"join_key"`
}

type DestinationGroupSpec struct {
	CollaboratorRef string            `yaml:"collaborator"`
	Destinations    []DestinationSpec `yaml:"destinations"`
}

type DestinationSpec struct {
	Name              string      `yaml:"name"`
	Requestee         string      `yaml:"request"`
	TransformationRef address.Ref `yaml:"transformation_ref"`
}

type DestinationAllowedSpec struct {
	Ref           address.Ref   `yaml:"ref"`
	NoiseParams   []any         `yaml:"noise_parameters"`
}
