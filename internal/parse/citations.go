package parse

import (
	"bytes"
	"fmt"
	"regexp"
)

type Citations map[string]string

var citationParseRegexp = regexp.MustCompile(`(?m)^([^\s]+):\s(.+)$`)

func CitationsFromSources(sources []*Source) (Citations, error) {

	cit := make(Citations)

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

				cit[string(subs[1])] = string(subs[2])

			}
		}

	}

	return cit, nil
}
