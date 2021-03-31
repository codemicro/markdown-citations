package parse

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var referenceRegexp = regexp.MustCompile(`(?m)\[r:([^\s]+)\]`)

// TransformReferences finds any references in the form [r:reference-name] in fcont and substitutes them
// for a numbered link with the corresponding URL.
//
// If `<!-- c:footer -->` is found within fcont, a footer is generated of all citations and the associated
// note with them, if one is supplied.
func TransformReferences(fcont *[]byte, citations map[string]*Citation) error {

	numerical := make(map[string]int)
	c := 1

	matches := referenceRegexp.FindAllSubmatch(*fcont, -1)
	if matches == nil {
		// no matches found
		return nil
	}

	for _, match := range matches {
		stringMatchName := strings.ToLower(string(match[1]))

		cit, foundCitation := citations[stringMatchName]
		if !foundCitation {
			return fmt.Errorf("unrecognised reference to citation name '%s'", stringMatchName)
		}

		preC, found := numerical[stringMatchName]
		n := preC
		if !found {
			n = c
			c += 1
		}

		numerical[stringMatchName] = n

		label := fmt.Sprintf("[%d]", n)
		*fcont = bytes.Replace(*fcont, match[0], []byte(makeMarkdownLink(label, cit.URL)), 1)

	}

	createFooter(fcont, numerical, citations)

	return nil
}

func makeMarkdownLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

var footerMarkerRegexp = regexp.MustCompile(`(?m)<!-- ?c:footer ?-->`)

func createFooter(fcont *[]byte, mapping map[string]int, citations map[string]*Citation) error {

	var lines sort.StringSlice

	for key, n := range mapping {

		c := citations[key]

		linkSection := c.Text
		if linkSection != "" {
			linkSection += " - "
		}
		linkSection += makeMarkdownLink(c.URL, c.URL)

		lines = append(lines, fmt.Sprintf("%d: %s", n, linkSection))
	}

	lines.Sort()
	textBlock := strings.Join(lines, "\n")
	*fcont = footerMarkerRegexp.ReplaceAll(*fcont, []byte(textBlock))
	return nil
}
