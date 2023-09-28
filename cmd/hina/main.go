package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nobytesguy/hina/internal/hina"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("no input files")
	}
	if len(args) > 1 {
		panic("hina can only interpret one file at a time")
	}

	filePath := args[0]
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

	env := hina.NewEnv()
	err = hina.EvalTree(jsonContent, env)
	if err != nil {
		panic(err)
	}
}
