package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
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

	case "target":

		files, err := ioutil.ReadDir(vm.Config.ConfigDirectory)
		if err != nil {
			logrus.Error("Cannot enumerate " + vm.Config.ConfigDirectory)
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
				vm.Config.Target = target
				WriteGopsConfig(&vm)
			}
		}

		if vm.Config.Target == "" {
			logrus.Error("No target set..")
		} else {
			logrus.Info("Current target: " + vm.Config.Target)
		}

		logrus.Info("Available targets:")

		for _, f := range files {
			logrus.Info("- " + f.Name())
		}

	case "templates":

		LoadAndMergeConfigDirectory(&vm)
		ProcessTemplates(&vm)

	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
