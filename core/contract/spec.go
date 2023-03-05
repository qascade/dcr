package contract

type WarehouseType string

const (
	SnowFlake WarehouseType = "snowflake"
	BigQuery  WarehouseType = "bigquery"
)

type ContractSpec struct {
	Name             string             `yaml:"name"`
	Version          string             `yaml:"version"`
	Purpose          string             `yaml:"omitempty,purpose"`
	Collaborators    []CollaboratorSpec `yaml:"collaborators"`
	ComputeWarehouse WarehouseType      `yaml:"compute_warehouse"`
}

type CollaboratorSpec struct {
	Name       string          `yaml:"name"`
	gitRepo    string          `yaml:"git_repo"`
	UserAgents []UserAgentSpec `yaml:"user_agents"`
	Warehouse  WarehouseType   `yaml:"warehouse"`
}

type UserAgentSpec struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}
