package parse

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

// Source represents a text file that contains citations
type Source struct {
	Filename string
	Content  []byte
}

var sourceImportRegexp = regexp.MustCompile(`(?m)<!-- ?c:source=(.+\S) ?-->`)

// FindSources finds any references to a citation file in a source document. The return value
// is a slice of all these sources. Each source with a given name will appear only once.
func FindSources(fcont *[]byte) []*Source {

	matches := sourceImportRegexp.FindAllStringSubmatch(string(*fcont), -1)
	*fcont = sourceImportRegexp.ReplaceAll(*fcont, nil)

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
func LoadSources(sources []*Source, basepath string) error {
	for _, source := range sources {

		x := filepath.Join(basepath, source.Filename)
		fc, err := ioutil.ReadFile(x)
		if err != nil {
			return err
		}

		source.Content = fc
	}
	return nil
}
