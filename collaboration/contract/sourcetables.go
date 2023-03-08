package contract

// TODO - Should we move these structs to package source? 
type TableRegister struct {
	collaborator *Collaborator
	SourceTables map[string]SourceTable // Mapped with schema
}

type SourceTable struct {
	Name            string
	Database        string
	Schema          string
	Description     string
	CollumnsAllowed map[string]Column // Mapped with table Name
}

type Column struct {
	Name           string
	Masking        MaskingType
	QueriesAllowed map[string]QueryType // Mapped with table name
}
