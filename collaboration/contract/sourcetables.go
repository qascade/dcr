package contract

// TODO - Should we move these structs to package source?
type TableRegister struct {
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
	QueriesAllowed []QueryType // Mapped with table name
}

func NewTableRegister(tcSpecs *[]TablesContractSpec, collaboratorName string) *TableRegister {
	tableRegister := &TableRegister{
		SourceTables: make(map[string]SourceTable),
	}
	for _, tcSpec := range *tcSpecs {
		if collaboratorName != tcSpec.Name {
			continue
		}
		for _, stSpec := range tcSpec.SourceTables {
			tableRegister.SourceTables[stSpec.Name] = SourceTable {
				Name: stSpec.Name,
				Description: stSpec.Description,
				Database: stSpec.Database,
				Schema: stSpec.Schema,
				CollumnsAllowed: registerColumns(stSpec.ColumnsAllowed, stSpec.Name),
			}
		}
	}
	return tableRegister
}

func registerColumns(columns []ColumnSpec, stName string) map[string]Column {
	columnsAllowed := make(map[string]Column)
	for _, c := range columns {
		columnsAllowed[c.Name] = Column{
			Name: c.Name,
			Masking: MaskingType(c.MaskingType),
			QueriesAllowed: c.queriesAllowed,
		}
	}
	return columnsAllowed
}
