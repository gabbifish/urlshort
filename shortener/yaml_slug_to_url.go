package shortener

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type yamlSlugToURL struct {
	mappings map[string]string
}

func NewYAMLSlugToURL(yamlFilename string) (SlugToURL, error) {
	yml, readYamlErr := ioutil.ReadFile(yamlFilename)
	if readYamlErr != nil {
		return nil, readYamlErr
	}
	return NewYAMLSlugToURLFromBytes(yml)
}

func NewYAMLSlugToURLFromBytes(yml []byte) (SlugToURL, error) {
	type URLExpand struct {
		Slug string `yaml:"path"`
		Full string `yaml:"url"`
	}

	var entries []URLExpand

	unmarshalErr := yaml.Unmarshal(yml, &entries)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	// Turn list of yaml entries into map
	m := make(map[string]string)
	for _, entry := range entries {
		m[entry.Slug] = entry.Full
	}
	return &yamlSlugToURL{m}, nil
}

func (y *yamlSlugToURL) URLExists(slug string) (bool, error) {
	_, ok := y.mappings[slug]
	return ok, nil
}

func (y *yamlSlugToURL) URL(slug string) (string, error) {
	url, ok := y.mappings[slug]
	if !ok {
		return "", fmt.Errorf("URL not found for slug \"%s\"", slug)
	}
	return url, nil
}
