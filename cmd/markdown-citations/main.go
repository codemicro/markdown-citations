package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codemicro/markdown-citations/internal/parse"
)

var errLog = log.New(os.Stderr, "", 0)

func main() {

	if len(os.Args) < 2 {
		errLog.Fatal("missing source file name")
	}

	// read source markdown file
	fcont, err := os.ReadFile(os.Args[1])
	if err != nil {
		errLog.Fatalf("unable to read source file: %s\n", err.Error())
	}

	sourceFiles := parse.FindSources(fcont)
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

	if err = ioutil.WriteFile(os.Args[1], fcont, 0644); err != nil {
		errLog.Fatal(err)
	}

	fmt.Println("done")
}
