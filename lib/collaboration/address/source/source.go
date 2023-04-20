// This package will have all the different types of sources
package source

import (
	"fmt"
	"os"

	"github.com/qascade/dcr/lib/collaboration/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

// All source types must implement this interface.
type Source interface {
	Extract() ([]byte, error)
}

type Column struct {
	Name              string
	Type              string
	MaskingType       string
	Selectable        bool
	JoinKey           bool
	AggregatesAllowed []string
}

type LocalSource struct {
	CollaboratorName string
	CsvLocation      string
	Columns          []Column
}

func NewLocalSource(cName string, sPec config.SourceSpec) Source {
	var columns []Column
	for _, col := range sPec.Columns {
		col := Column{
			Name:              col.Name,
			Type:              col.Type,
			MaskingType:       col.MaskingType,
			Selectable:        col.Selectable,
			JoinKey:           col.JoinKey,
			AggregatesAllowed: col.AggregatesAllowed,
		}
		columns = append(columns, col)
	}
	return &LocalSource{
		CollaboratorName: cName,
		CsvLocation:      sPec.CSVLocation,
		Columns:          columns,
	}
}

func (ls *LocalSource) Extract() ([]byte, error) {
	// Read CSV file
	csv, err := os.Open(ls.CsvLocation)
	var csvBytes []byte
	if err != nil {
		return nil, fmt.Errorf("error while opening the csv file: %v", err)
	}
	_, err = csv.Read(csvBytes)
	if err != nil {
		return nil, fmt.Errorf("error while reading the csv file: %v", err)
	}
	return csvBytes, nil
}
