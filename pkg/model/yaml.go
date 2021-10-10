package model

import (
	"io"

	"github.com/maargenton/go-fileutils"
	"gopkg.in/yaml.v3"
)

// LoadFromYAML loads the internal representation of a list for packages and
// their tests from a stream containing the native YAML representation.
func LoadFromYAML(r io.Reader) (pkgs []Package, err error) {
	d := yaml.NewDecoder(r)
	if err := d.Decode(pkgs); err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		pkg.LinkTests()
	}
	return
}

// LoadFromYAMLFile loads the internal representation of a list for packages and
// their tests from a file containing the native YAML representation.
func LoadFromYAMLFile(filename string) (pkgs []Package, err error) {
	err = fileutils.ReadFile(filename, func(r io.Reader) error {
		pkgs, err = LoadFromYAML(r)
		return err
	})
	return
}

// SaveToYAML save the internal representation of a list for packages and their
// tests to a stream in the native YAML format.
func SaveToYAML(w io.Writer, pkgs []Package) error {
	e := yaml.NewEncoder(w)
	return e.Encode(pkgs)
}

// SaveToYAMLFile save the internal representation of a list for packages and
// their tests to a file in the native YAML format.
func SaveToYAMLFile(filename string, pkgs []Package) error {
	return fileutils.WriteFile(filename, func(w io.Writer) error {
		return SaveToYAML(w, pkgs)
	})
}
