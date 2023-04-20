package config

type Spec interface{}

type TransformationGroupSpec struct {
	CollaboratorRef string               `yaml:"collaborator"`
	Transformations []TransformationSpec `yaml:"transformations"`
}

// For now only supporting Count Query.
// TODO - Make this Spec such that any type of transformation can be parsed using this spec.
type TransformationSpec struct {
	Name               string     `yaml:"name"`
	Count              string     `yaml:"count"`
	From               []FromSpec `yaml:"from"`
	NoiseType          string     `yaml:"noise_type"`
	NoiseParams        []string   `yaml:"noise_parameters"`
	JoinKey            string     `yaml:"join_key"`
	Template           string     `yaml:"template"`
	ConsumerAllowed    []string   `yaml:"consumer_allowed"`
	DestinationAllowed []string   `yaml:"destination_allowed"`
}

type SourceGroupSpec struct {
	CollaboratorRef string       `yaml:"collaborator"`
	Sources         []SourceSpec `yaml:"sources"`
}

type SourceSpec struct {
	Name        string       `yaml:"name"`
	CSVLocation string       `yaml:"csv_location"`
	Description string       `yaml:"description"`
	Columns     []ColumnSpec `yaml:"columns"`
	// TODO- Do we need to add addressRef here?
	ConsumersAllowed    []string                 `yaml:"consumers_allowed"`
	DestinationsAllowed []DestinationAllowedSpec `yaml:"destinations_allowed"`
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
	Name string `yaml:"name"`
	Ref  string `yaml:"ref"`
}

type SourceDestinationAllowedSpec struct {
	Ref         string `yaml:"ref"`
	NoiseParams []any  `yaml:"noise_parameters"`
}

type FromSpec struct {
	Name string `yaml:"name"`
	Ref  string `yaml:"ref"`
}
