package contract

type ContractSpec struct {
	Name             string             `yaml:"name"`
	Version          string             `yaml:"version"`
	Purpose          string             `yaml:"purpose,omitempty"`
	Collaborators    []CollaboratorSpec `yaml:"collaborators"`
	ComputeWarehouse WarehouseType      `yaml:"compute_warehouse"`
}

type CollaboratorSpec struct {
	Name       string          `yaml:"name"`
	GitRepo    string          `yaml:"git_repo"`
	UserAgents []UserAgentSpec `yaml:"user_agents"`
	Warehouse  WarehouseType   `yaml:"warehouse"`
}

type UserAgentSpec struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type TablesContractSpec struct {
	Name         string            `yaml:"name"`
	SourceTables []SourceTableSpec `yaml:"tables"`
}

type SourceTableSpec struct {
	Name           string       `yaml:"name"`
	Database       string       `yaml:"database"`
	Schema         string       `yaml:"schema"`
	Description    string       `yaml:"description"`
	ColumnsAllowed []ColumnSpec `yaml:"columns_allowed"`
}

type ColumnSpec struct {
	Name           string      `yaml:"name"`
	MaskingType    string      `yaml:"masking_type"`
	QueriesAllowed []QueryType `yaml:"queries_allowed"`
}
