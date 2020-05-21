package config

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLoader_Load(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cl := NewConfigLoader()
	cfg, err := cl.Load(filepath.Join(cwd, "testdata", "test.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Tasks["task1"] == nil || cfg.Tasks["task1"].Commands[0] != "echo true" {
		t.Error("yaml parsing failed")
	}

	cl = NewConfigLoader()
	cfg, err = cl.Load(filepath.Join(cwd, "testdata", "test.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Tasks["task1"] == nil || cfg.Tasks["task1"].Commands[0] != "echo true" {
		t.Error("yaml parsing failed")
	}

	_, err = cl.LoadGlobalConfig()
	if err != nil {
		t.Fatal()
	}
}

func TestLoader_resolveDefaultConfigFile(t *testing.T) {
	cl := NewConfigLoader()

	cl.dir = filepath.Join(cl.dir, "testdata")
	file, err := cl.resolveDefaultConfigFile()
	if err != nil {
		t.Fatal(err)
	}

	if filepath.Base(file) != "tasks.yaml" {
		t.Error()
	}

	cl.dir = "/"
	file, err = cl.resolveDefaultConfigFile()
	if err == nil || file != "" {
		t.Error()
	}
}

func TestLoader_loadDir(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cl := NewConfigLoader()
	m, err := cl.loadDir(filepath.Join(cwd, "testdata"))
	if err != nil {
		t.Fatal(err)
	}

	tasks := m["tasks"].(map[interface{}]interface{})
	if len(tasks) != 5 {
		t.Error()
	}
}

func TestLoader_readURL(t *testing.T) {
	var r int
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "")
		if r == 0 {
			writer.Header().Set("Content-Type", "application/json")
			r++
		}
		fmt.Fprintln(writer, "{\"tasks\": {\"task1\": {\"command\": [\"true\"]}}}")
	}))

	cl := NewConfigLoader()
	m, err := cl.readURL(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	tasks := m["tasks"].(map[string]interface{})
	if len(tasks) != 1 {
		t.Error()
	}

	_, err = cl.readURL(srv.URL)
	if err != nil {
		t.Fatal()
	}
}
