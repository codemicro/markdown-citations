package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/codemicro/markdown-citations/internal/parse"
)

var errLog = log.New(os.Stderr, "", 0) // zero so it just prints the message and nothing else

func main() {

	if len(os.Args) < 2 {
		errLog.Fatal("missing source file name")
	}

	// read source markdown file
	fcont, err := os.ReadFile(os.Args[1])
	if err != nil {
		errLog.Fatalf("unable to read source file: %s\n", err.Error())
	}

	sourceFiles := parse.FindSources(&fcont)
	err = parse.LoadSources(sourceFiles)
	if err != nil {
		errLog.Fatal(err)
	}

	citations, err := parse.CitationsFromSources(sourceFiles)
	if err != nil {
		errLog.Fatal(err)
	}

	if err = parse.TransformReferences(&fcont, citations); err != nil {
		errLog.Fatal(err)
	}

	var outfile string
	{
		x := strings.Split(os.Args[1], ".")
		outfile = strings.Join(append(x[:len(x)-1], "gen", x[len(x)-1]), ".")
	}

	if err = ioutil.WriteFile(outfile, bytes.TrimSpace(fcont), 0644); err != nil {
		errLog.Fatal(err)
	}

	fmt.Println("Written to", outfile)
}
