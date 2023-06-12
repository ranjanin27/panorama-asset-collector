// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package tracer

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type TestData struct {
	t       *testing.T
	dataDir string
}

func (td *TestData) MustGetContents(filename string) []byte {
	path := filepath.Join(td.dataDir, filename)
	file, err := os.Open(path)
	if err != nil {
		td.t.Fatal(err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			td.t.Fatal(closeErr)
		}
	}()
	byteSequence, err := io.ReadAll(file)
	if err != nil {
		td.t.Fatal(err)
	}
	return byteSequence
}

func (td *TestData) MustGetJSON(filename string, contents interface{}) {
	byteSequence := td.MustGetContents(filename)
	err := json.Unmarshal(byteSequence, contents)
	if err != nil {
		td.t.Fatal(err)
	}
}

func NewTestDataInterface(dir string, t *testing.T) *TestData {
	return &TestData{
		t:       t,
		dataDir: dir,
	}
}
