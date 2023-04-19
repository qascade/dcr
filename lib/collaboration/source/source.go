// This package will have all the different types of sources
package source

// All source types must implement this interface.
type Source interface {
	Extract() error
}

type LocalSource struct {
	CollaboratorRef string
	CsvLocation     string
}
