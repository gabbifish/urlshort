package shortener

//go:generate mockery -name SlugToURL

type SlugToURL interface {
	URLExists(slug string) (bool, error)
	URL(slug string) (string, error)
}
