package parse

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var referenceRegexp = regexp.MustCompile(`(?m)\[r:([^\s]+)\]`)

func TransformReferences(fcont *[]byte, citations Citations) error {

	numerical := make(map[string]int)
	c := 1

	matches := referenceRegexp.FindAllSubmatch(*fcont, -1)
	if matches == nil {
		// no matches found
		return nil
	}

	for _, match := range matches {
		stringMatchName := string(match[1])

		url, foundCitation := citations[stringMatchName]
		if !foundCitation {
			return fmt.Errorf("unrecognised reference to sitation name '%s'", stringMatchName)
		}

		preC, found := numerical[url]
		n := preC
		if !found {
			n = c
			c += 1
		}

		numerical[url] = n

		label := fmt.Sprintf("[%d]", n)
		*fcont = bytes.Replace(*fcont, match[0], []byte(makeMarkdownLink(label, url)), 1)

	}

	createFooter(fcont, numerical)

	return nil
}

func makeMarkdownLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

var footerMarkerRegexp = regexp.MustCompile(`(?m)<!-- ?c:footer ?-->`)

func createFooter(fcont *[]byte, mapping map[string]int) error {

	var lines sort.StringSlice

	for url, n := range mapping {
		lines = append(lines, fmt.Sprintf("%d: %s", n, makeMarkdownLink(url, url)))
	}

	lines.Sort()
	textBlock := strings.Join(lines, "\n")
	*fcont = footerMarkerRegexp.ReplaceAll(*fcont, []byte(textBlock))
	return nil
}
