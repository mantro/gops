package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

func main() {

	args := os.Args[1:]

	command := "templates"
	if len(args) > 0 {
		command = args[0]
	}

	vm := CreateViewModel()
	LoadGopsConfig(&vm)

	switch command {
	case "config":
		output, _ := yaml.Marshal(vm.Config)
		fmt.Println(string(output))

	case "dump":
		LoadAndMergeConfigDirectory(&vm)
		output, _ := yaml.Marshal(&vm)
		fmt.Println(string(output))

	case "help":
		fmt.Println("Available commands: config dump help target templates")

	case "current-target":
		fmt.Println(vm.Meta.Target)

	case "target":

		configDirectory := filepath.Join(vm.Meta.GitRoot, vm.Config.ConfigDirectory)
		files, err := ioutil.ReadDir(configDirectory)
		if err != nil {
			logrus.Error("Cannot enumerate " + configDirectory)
			panic(err)
		}

		var fileNames []string
		for _, file := range files {
			fileNames = append(fileNames, file.Name())
		}

		if len(args) > 1 {
			target := args[1]
			if !SliceContains(fileNames, target) {
				logrus.Error("Cannot find target: " + target)
			} else {
				vm.Meta.Target = target
				WriteGopsConfig(&vm)
			}
		} else {
			logrus.Info("Available targets:")

			for _, f := range files {
				logrus.Info("- " + f.Name())
			}
		}

		if vm.Meta.Target == "" {
			logrus.Error("No target set..")
		} else {
			logrus.Info("Current target: " + vm.Meta.Target)
		}

	case "templates":

		LoadAndMergeConfigDirectory(&vm)
		ProcessTemplates(&vm)

	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
