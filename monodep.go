package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type Dep struct {
	Name    string
	Path    string
	Deps    []string
	changed bool
}

type DepMap map[string]*Dep

func main() {
	// TODO: make these into os.Args
	if len(os.Args) < 4 {
		fmt.Println("Expected 3 arguments, but got", len(os.Args)-1)
		fmt.Println("monodep deps.yaml hash1 hash2")
		os.Exit(1)
	}
	yamlFile := os.Args[1]
	first := os.Args[2]
	second := os.Args[3]
	// first := "766b95f345432a33d594d4291748a6af53e682b3"
	// second := "fac51d26fad6f2ea535cd6b2fc70ad65288542aa"
	depMap := make(DepMap)
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		fmt.Println("can't find the specified yaml file:", yamlFile)
		os.Exit(1)
	}
	err = yaml.Unmarshal([]byte(data), &depMap)
	if err != nil {
		fmt.Println("can't parse the yaml file:", err)
		os.Exit(1)
	}
	baseDir := path.Dir(yamlFile)
	for name, dep := range depMap {
		dep.Name = name
		dep.Path = path.Join(baseDir, dep.Path)
	}

	depMap.checkDeps()

	cmd := exec.Command("git", "diff", "--name-only", first, second)
	cmdOutput, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Start()
	reader := bufio.NewReader(cmdOutput)
	line, _, err := reader.ReadLine()
	for err == nil {
		fileDir := path.Dir(string(line))
		for _, dep := range depMap {
			if strings.Contains(fileDir, dep.Path) {
				dep.changed = true
			}
		}
		line, _, err = reader.ReadLine()
	}

	for key := range depMap {
		if depMap.shouldRecompile(key) {
			fmt.Println(key)
		}
	}

}

func (m DepMap) checkDeps() {
	for _, d := range m {
		for _, dep := range d.Deps {
			if _, ok := m[dep]; !ok {
				fmt.Println(d.Name, "is depending on", dep, "which is not declared")
				os.Exit(1)
			}
		}
		if _, err := os.Stat(d.Path); os.IsNotExist(err) {
			fmt.Println(d.Name, "has a path which can't be found:", d.Path)
			os.Exit(1)
		}
	}
}

func (m DepMap) shouldRecompile(key string) bool {
	dep := m[key]
	if dep == nil {
		return false
	}
	if dep.changed {
		return true
	}
	for _, d := range dep.Deps {
		if m.shouldRecompile(d) {
			return true
		}
	}
	return false
}
