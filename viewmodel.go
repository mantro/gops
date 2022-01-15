package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

type Dictionary = map[string]interface{}

type Meta struct {
	GitRoot        string
	ConfigFilePath string
}

type GopsConfig struct {
	ConfigDirectory string
	Target          string
	DelimiterLeft   string
	DelimiterRight  string
}

type ViewModel struct {
	Meta   Meta
	Config GopsConfig
	Data   Dictionary
}

func CreateViewModel() (vm ViewModel) {

	gitRoot, searchPath, err := GetGitRoot()
	if err != nil {
		logrus.Error("No git repository found at " + searchPath)
		os.Exit(1)
	}

	vm = ViewModel{}
	vm.Meta = Meta{}
	vm.Meta.GitRoot = gitRoot
	vm.Meta.ConfigFilePath = filepath.Join(gitRoot, ".gops.yaml")

	return vm
}

func setViperDefaults(vm *ViewModel) {

	viper.SetDefault("ConfigDirectory", "ops/config")
	viper.SetDefault("Target", "")
	viper.SetDefault("DelimiterLeft", "#{{")
	viper.SetDefault("DelimiterRight", "}}#")
	viper.SetConfigFile(vm.Meta.ConfigFilePath)
	viper.SetConfigType("yaml")
}

func getViperConfig() GopsConfig {

	return GopsConfig{
		ConfigDirectory: viper.GetString("ConfigDirectory"),
		Target:          viper.GetString("Target"),
		DelimiterLeft:   viper.GetString("DelimiterLeft"),
		DelimiterRight:  viper.GetString("DelimiterRight"),
	}
}

func WriteGopsConfig(vm *ViewModel) {
	output, _ := yaml.Marshal(vm.Config)

	if err := os.WriteFile(vm.Meta.ConfigFilePath, output, 0644); err != nil {
		logrus.Error("Cannot open " + vm.Meta.ConfigFilePath + " for writing")
		panic(err)
	}

	logrus.Info("Written config to " + vm.Meta.ConfigFilePath)
}

func LoadGopsConfig(vm *ViewModel) {

	setViperDefaults(vm)

	if _, err := os.Stat(vm.Meta.ConfigFilePath); err == nil {

		if err := viper.ReadInConfig(); err != nil {
			logrus.Error("Cannot load config file " + vm.Meta.ConfigFilePath)
			panic(err)
		}

		vm.Config = getViperConfig()

	} else {
		vm.Config = getViperConfig()
		WriteGopsConfig(vm)
	}

	contents := ""

	gitignorePath := filepath.Join(vm.Meta.GitRoot, ".gitignore")
	if gitignore, err := os.ReadFile(gitignorePath); err != nil {
		logrus.Warn(gitignorePath + " does not exist, creating..")
		contents = ".gops.yaml"
	} else {
		contents = string(gitignore)

		if !strings.Contains(contents, ".gops.yaml") {
			logrus.Info("Adding '.gops.yaml' to " + gitignorePath)
			contents = contents + "\n.gops.yaml"
		}
	}

	if err := os.WriteFile(gitignorePath, []byte(contents), 0644); err != nil {
		logrus.Error("Cannot write " + gitignorePath)
		panic(err)
	}
}

func LoadAndMergeConfigDirectory(vm *ViewModel) {

	if vm.Config.Target == "" {
		logrus.Error("No target has been set, please invoke 'gops target' first")
		os.Exit(1)
	}

	vm.Data = Dictionary{}

	directories, err := ioutil.ReadDir(vm.Config.ConfigDirectory)
	if err != nil {
		logrus.Error("Directory does not exist:  " + RelPath(vm.Meta.GitRoot, vm.Config.ConfigDirectory))
		os.Exit(1)
	}

	for _, directory := range directories {
		if directory.IsDir() {

			files, _ := filepath.Glob(filepath.Join(vm.Config.ConfigDirectory, directory.Name(), "*.yaml"))
			merged := LoadYamlFiles(files...)

			vm.Data[directory.Name()] = merged
		}
	}

	if _, ok := vm.Data["target"]; ok {
		logrus.Error("There is a config directory called target, bailing")
		os.Exit(1)
	}

	if val, ok := vm.Data[vm.Config.Target]; ok {
		vm.Data["target"] = val
	} else {
		logrus.Error("Cannot find target " + vm.Config.Target + " in merged configuration")
		os.Exit(1)
	}

	logrus.Info("Current target: " + vm.Config.Target)
}

func ProcessTemplates(vm *ViewModel) {

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader("/"),
		jet.InDevelopmentMode(),
		jet.WithDelims(vm.Config.DelimiterLeft, vm.Config.DelimiterRight),
	)

	files := Glob(vm.Meta.GitRoot, ".template.")
	for _, file := range files {
		TemplateFile(views, file, *vm)
	}
}
