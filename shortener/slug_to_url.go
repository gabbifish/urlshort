package shortener

type SlugToURL interface {
	URLExists(slug string) (bool, error)
	URL(slug string) (string, error)
}
