package main

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type Dep struct {
	Name    string
	Path    string
	Deps    []string
	changed bool
}

func main() {
	first := "b2ec1fe58aa1b94fdd5c4a4c5ee1f22e2ecf2c0f"
	second := "6a5e1d74a101f2015242793fe3a93acbf8321416"
	depMap := make(map[string]Dep)
	// TODO: make this an argument
	yamlFile := "test_data/test1/deps.yaml"
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(data), &depMap)
	if err != nil {
		panic(err)
	}
	baseDir := path.Dir(yamlFile)
	for name, dep := range depMap {
		dep.Name = name

	}
}
