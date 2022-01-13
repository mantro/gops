package main

import (
	"fmt"
	"os"
	"sigs.k8s.io/yaml"
)

func main() {

	args := os.Args[1:]

	command := "templates"
	if len(args) > 0 {
		command = args[0]
	}

	vm := CreateViewModel()

	switch command {
	case "config":
		LoadGoopsConfig(&vm)
		output, _ := yaml.Marshal(vm.Config)
		fmt.Println(string(output))

	case "dump":
		LoadGoopsConfig(&vm)
		LoadAndMergeConfigDirectory(&vm)
		output, _ := yaml.Marshal(&vm)
		fmt.Println(string(output))

	case "help":
		fmt.Println("Available commands: config dump init help target templates")

	case "init":
		LoadGoopsConfig(&vm)
		InitializeGoopsConfig(&vm)

	case "target":

	case "templates":

		LoadGoopsConfig(&vm)
		LoadAndMergeConfigDirectory(&vm)
		ProcessTemplates(&vm)

	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
