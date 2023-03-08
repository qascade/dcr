package contract 

type WarehouseType string
type QueryType string
type MaskingType string
type SpecType string

type Spec interface {}

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
	ContractSpecType SpecType = "contract"
	TablesContractSpecType SpecType = "table"
)