package shortener

import "fmt"

type mapSlugToURL struct {
	mappings map[string]string
}

func NewMapSlugToURL(mappings map[string]string) SlugToURL {
	return &mapSlugToURL{
		mappings: mappings,
	}
}

func (m *mapSlugToURL) URLExists(slug string) (bool, error) {
	_, ok := m.mappings[slug]
	return ok, nil
}

func (m *mapSlugToURL) URL(slug string) (string, error) {
	url, ok := m.mappings[slug]
	if !ok {
		return "", fmt.Errorf("URL not found for slug \"%s\"", slug)
	}
	return url, nil
}
