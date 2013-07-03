package main

import (
	"flag"
	"github.com/sourcegraph/django-docs"
	"log"
)

var docDir = flag.String("docDir", "", "docs directory for django installation")

func max(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func main() {
	flag.Parse()

	docs, errs := django_docs.ExtractDocs(*docDir)

	log.Printf("Docs")
	for _, doc := range docs {
		log.Printf("\nsymbol: %s\nfile: %s\nlocation: %s:%d:%d\nbody: %s\n\n\n\n\n",
			doc.Symbol, doc.SourceFile, doc.SourceFile, doc.Start, doc.End, doc.Body)
	}

	log.Printf("########## Errors ###########")
	for _, err := range errs {
		log.Printf("%v", err)
	}
}
