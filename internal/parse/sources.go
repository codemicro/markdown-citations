package parse

import (
	"os"
	"regexp"
)

type Source struct {
	Filename string
	Content  []byte
}

var sourceImportRegexp = regexp.MustCompile(`(?m)<!-- ?c:source=(.+\S) ?-->`)

// FindSources finds any references to a citation file in a source document. The return value
// is a slice of all these sources. Each source with a given name will appear only once.
func FindSources(fcont []byte) []*Source {

	matches := sourceImportRegexp.FindAllStringSubmatch(string(fcont), -1)

	var o []*Source
	u := make(map[string]struct{})

	for _, match := range matches {
		sourceName := match[1]

		if _, alreadyInSlice := u[sourceName]; !alreadyInSlice {
			o = append(o, &Source{Filename: sourceName})
			u[sourceName] = struct{}{}
		}

	}

	return o
}

// LoadSources takes a slice of *Source and reads the associated files from disk
func LoadSources(sources []*Source) error {
	for _, source := range sources {

		// TODO: make this filepath relative to the source file path
		file, err := os.Open(source.Filename)
		defer file.Close()
		if err != nil {
			return err
		}

		fstat, err := file.Stat()
		if err != nil {
			return err
		}

		source.Content = make([]byte, fstat.Size())

		if _, err = file.Read(source.Content); err != nil {
			return err
		}

	}
	return nil
}
