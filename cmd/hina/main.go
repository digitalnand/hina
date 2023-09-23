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
	if len(arguments) > 1 {
		panic("hina can only interpret one file at a time")
	}

	filePath := arguments[0]
	if filepath.Ext(filePath) != ".json" {
		panic("file format not recognized")
	}
	fileContent, err := os.ReadFile(filePath)
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
