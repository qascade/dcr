package source

type Source interface {
	Extract() error
}

type LocalSource struct {
	Location string
}

func NewLocalSource(location string) Source {
	return &LocalSource{
		Location: location,
	}
}

func  (*LocalSource) Extract() error {
	return nil
}