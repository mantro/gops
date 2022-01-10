package main

import (
	"bytes"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
)

var logger, _ = zap.NewDevelopment()

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("/"),
	jet.InDevelopmentMode(),
	jet.WithDelims("#{{", "}}#"),
)

type dictionary = map[string]interface{}

type meta struct {
	RootDirectory   string
	ConfigDirectory string
}

type viewmodel struct {
	Meta meta
	Data dictionary
}

func generateViewModel() (vm viewmodel) {
	path, err := os.Getwd()
	if err != nil {
		panic("Cannot get current directory")
	}

	pathOverride := os.Getenv("GOOPS_PATH")
	if pathOverride != "" {
		path = pathOverride
		logger.Info("OVERRIDE path: " + path)
	}

	gitRoot, err := GetGitRoot(path)
	if err != nil {
		logger.Error("No git repository found at " + path)
		os.Exit(1)
	}

	vm = viewmodel{}
	vm.Meta.RootDirectory = gitRoot
	vm.Meta.ConfigDirectory = filepath.Join(gitRoot, "ops", "config")
	return vm
}

func generateConfig(vm *viewmodel) {

	vm.Data = dictionary{}

	directories, err := ioutil.ReadDir(vm.Meta.ConfigDirectory)
	if err != nil {
		logger.Error("Cannot iterate " + RelPath(vm.Meta.RootDirectory, vm.Meta.ConfigDirectory))
		os.Exit(1)
	}

	for _, directory := range directories {
		if directory.IsDir() {

			files, _ := filepath.Glob(filepath.Join(vm.Meta.ConfigDirectory, directory.Name(), "*.yaml"))
			merged := LoadYamlFiles(files...)

			vm.Data[directory.Name()] = merged
		}
	}
}

func templateFile(file string, vm viewmodel) {

	relativeFile := RelPath(vm.Meta.RootDirectory, file)
	logger.Info("Templating " + relativeFile)

	info, err := os.Stat(file)
	if err != nil {
		logger.Error("File does not exist: " + relativeFile)
	}

	tmpl, err := views.GetTemplate(file)
	if err != nil {
		logger.Error("Unable to parse template file " + relativeFile)
		panic(err)
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, nil, vm)
	if err != nil {
		logger.Error("Cannot execute template " + relativeFile)
		panic(err)
	}

	outputFile := strings.Replace(file, ".template.", ".generated.", 1)

	err = os.WriteFile(outputFile, buf.Bytes(), info.Mode())
}

func main() {

	args := os.Args[1:]

	command := "templates"
	if len(args) > 0 {
		command = args[0]
	}

	vm := generateViewModel()
	generateConfig(&vm)

	switch command {
	case "config":
		output, _ := yaml.Marshal(&vm)
		fmt.Println(string(output))
	case "templates":
		files := Glob(vm.Meta.RootDirectory, ".template.")
		for _, file := range files {
			templateFile(file, vm)
		}

	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
