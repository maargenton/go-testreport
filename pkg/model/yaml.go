package model

import (
	"io"

	"github.com/maargenton/go-fileutils"
	"gopkg.in/yaml.v3"
)

// LoadFromYAML loads the internal representation of a result object with its
// list of packages and their tests from a stream containing the native YAML
// representation.
func LoadFromYAML(r io.Reader) (results *Results, err error) {
	var d = yaml.NewDecoder(r)
	results = &Results{}
	if err := d.Decode(results); err != nil {
		return nil, err
	}
	results.UpdateCounts()
	return
}

// LoadFromYAMLFile loads the internal representation of a result object with
// its list of packages and their tests from a file containing the native YAML
// representation.
func LoadFromYAMLFile(filename string) (results *Results, err error) {
	err = fileutils.ReadFile(filename, func(r io.Reader) error {
		results, err = LoadFromYAML(r)
		return err
	})
	return
}

// SaveToYAML save the internal representation of a result object with its list
// of packages and their tests to a stream in the native YAML format.
func SaveToYAML(w io.Writer, results *Results) error {
	e := yaml.NewEncoder(w)
	return e.Encode(results)
}

// SaveToYAMLFile save the internal representation of a result object with its
// list of packages and their tests to a file in the native YAML format.
func SaveToYAMLFile(filename string, results *Results) error {
	return fileutils.WriteFile(filename, func(w io.Writer) error {
		return SaveToYAML(w, results)
	})
}
