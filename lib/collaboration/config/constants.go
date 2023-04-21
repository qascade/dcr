package config

type WarehouseType string
type QueryType string
type MaskingType string
type SpecType string

const (
	SnowFlake WarehouseType = "snowflake"
	BigQuery  WarehouseType = "bigquery"
)

const (
	Select QueryType = "select"
)

const (
	SHA256 MaskingType = "sha256"
)

const (
	CollaborationSpecType  SpecType = "collaboration_spec"
	SourceSpecType         SpecType = "source_spec"
	TransformationSpecType SpecType = "transformation_spec"
	DestinationSpecType    SpecType = "destination_spec"
)
