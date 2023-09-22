package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nobytesguy/hina/internal/hina"
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) == 0 {
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

	var jsonContent hina.Object
	err = json.Unmarshal(fileContent, &jsonContent)
	if err != nil {
		panic(err)
	}

	env := hina.NewEnvironment()
	err = hina.EvalTree(jsonContent, env)
	if err != nil {
		panic(err)
	}
}
