package parse

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Citation struct {
	URL  string
	Text string
}

var citationParseRegexp = regexp.MustCompile(`(?m)^([^\s]+):\s(.+)$`)

// CitationsFromSources takes a slice of Source pointers and extracts all citations from them
func CitationsFromSources(sources []*Source) (map[string]*Citation, error) {

	cit := make(map[string]*Citation)

	for _, source := range sources {

		lines := bytes.Split(source.Content, []byte("\n"))
		for i, line := range lines {

			line = bytes.ReplaceAll(line, []byte("\r"), nil) // compensate for CRLF

			// ignore blank lines and comments
			if !(bytes.Equal(line, nil) || bytes.HasPrefix(line, []byte("#"))) {

				subs := citationParseRegexp.FindSubmatch(line)
				if subs == nil {
					return nil, fmt.Errorf("cannot parse '%s' at %s:%d", line, source.Filename, i+1)
				}

				var c Citation
				{
					x := strings.Split(string(subs[2]), " ")
					c.URL = x[0]
					c.Text = strings.Join(x[1:], " ")
				}

				cit[string(subs[1])] = &c

			}
		}

	}

	return cit, nil
}
