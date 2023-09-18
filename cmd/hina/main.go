package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nobytesguy/hina/internal/hina"
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) < 1 {
		panic("no input files")
	}
	if filepath.Ext(arguments[0]) != ".json" {
		panic("file format not recognized")
	}

	// TODO: interpret multiple files
	fileContent, err := os.ReadFile(arguments[0])
	if err != nil {
		panic(err)
	}

	var jsonContent map[string]interface{}
	err = json.Unmarshal(fileContent, &jsonContent)
	if err != nil {
		panic(err)
	}

	err = hina.WalkTree(jsonContent)
	if err != nil {
		panic(err)
	}
}
