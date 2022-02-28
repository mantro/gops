package main

import (
	"github.com/imdario/mergo"
	"github.com/oleiade/reflections"
	"github.com/sirupsen/logrus"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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

	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if strings.Contains(s, extension) {
			files = append(files, s)
		}
		return nil
	})

	if err != nil {
		logrus.Error("Cannot enumerate " + root + " for files with " + extension)
		panic(err)
	}

	return files
}

func LoadYamlFile(file string) (result Dictionary, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	yamlData := Dictionary{}
	err = yaml.Unmarshal(data, &yamlData)

	return yamlData, nil
}

func EnsureLineInGitIgnore(gitRoot string, line string) {

	contents := ""

	gitignorePath := filepath.Join(gitRoot, ".gitignore")

	if gitignore, err := os.ReadFile(gitignorePath); err != nil {
		logrus.Warn(gitignorePath + " does not exist, creating..")
		contents = line
	} else {
		contents = string(gitignore)

		if !strings.Contains(contents, line) {
			logrus.Info("Adding '" + line + "' to " + gitignorePath)
			contents = contents + "\n" + line
		} else {
			// nothing to do
			return
		}
	}

	if err := os.WriteFile(gitignorePath, []byte(contents), 0644); err != nil {
		logrus.Error("Cannot write " + gitignorePath)
		panic(err)
	}
}

func LoadYamlFiles(files ...string) (result Dictionary) {

	merge := Dictionary{}

	for _, file := range files {
		contents, _ := LoadYamlFile(file)
		if err := mergo.Merge(&merge, contents); err != nil {
			logrus.Error("Cannot merge yaml file " + file)
			panic(err)
		}
	}

	return merge
}

func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetGitRoot() (string, string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		panic("Cannot get current directory")
	}

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = cwd

	path, err := cmd.Output()
	if err != nil {
		return "", cwd, err
	}

	return strings.TrimSpace(string(path)), cwd, nil
}

func Get(obj interface{}, prop string) (interface{}, error) {

	// Get the array access
	arr := strings.Split(prop, ".")

	var err error

	for _, key := range arr {
		obj, err = getProperty(obj, key)
		if err != nil {
			return nil, err
		}
		if obj == nil {
			return nil, nil
		}
	}
	return obj, nil
}

// Loop through this to get properties via dot notation
func getProperty(obj interface{}, prop string) (interface{}, error) {

	if reflect.TypeOf(obj).Kind() == reflect.Map {

		val := reflect.ValueOf(obj)

		valueOf := val.MapIndex(reflect.ValueOf(prop))

		if valueOf == reflect.Zero(reflect.ValueOf(prop).Type()) {
			return nil, nil
		}

		idx := val.MapIndex(reflect.ValueOf(prop))

		if !idx.IsValid() {
			return nil, nil
		}
		return idx.Interface(), nil
	}

	prop = strings.Title(prop)
	return reflections.GetField(obj, prop)
}
