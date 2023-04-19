// This package will have all the different types of destinations
package destination

// All Destination types must implement this interface
type Destination interface{}

type LocalDestination struct {
	OutputFolder   string
	DownloadAccess []string
}

func (l *LocalDestination) Authorize() bool {

	return false
}
