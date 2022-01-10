package main

import (
	"github.com/imdario/mergo"
	"io/fs"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
)

func RelPath(parent string, child string) string {

	if len(child) <= len(parent) {
		panic("Child path cannot be shorter that parent path")
	}
	return "." + child[len(parent):]
}

func Glob(root string, extension string) []string {

	var files []string

	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if strings.Contains(s, extension) {
			files = append(files, s)
		}
		return nil
	})

	return files
}

func LoadYamlFile(file string) (result dictionary, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	yamlData := dictionary{}
	err = yaml.Unmarshal(data, &yamlData)

	return yamlData, nil
}

func LoadYamlFiles(files ...string) (result dictionary) {

	merge := dictionary{}

	for _, file := range files {
		contents, _ := LoadYamlFile(file)
		mergo.Merge(&merge, contents)
	}

	return merge
}

func GetGitRoot(searchRoot string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = searchRoot

	path, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}
