package main

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

type Dictionary = map[string]interface{}

type Meta struct {
	GitRoot        string
	ConfigFilePath string
	TargetFilePath string
	Target         string
}

type GopsConfig struct {
	ConfigDirectory string
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
	vm.Meta.TargetFilePath = filepath.Join(gitRoot, ".gops-target.yaml")

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

	if err := os.WriteFile(vm.Meta.TargetFilePath, []byte(vm.Meta.Target), 0644); err != nil {
		logrus.Error("Cannot open " + vm.Meta.TargetFilePath + " for writing")
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

	if _, err := os.Stat(vm.Meta.TargetFilePath); err == nil {
		var contents []byte
		if contents, err = os.ReadFile(vm.Meta.TargetFilePath); err != nil {
			logrus.Error("Cannot read " + vm.Meta.TargetFilePath)
			panic(err)
		}

		vm.Meta.Target = string(contents)

		targetPath := filepath.Join(vm.Meta.GitRoot, vm.Config.ConfigDirectory, vm.Meta.Target)
		if _, err := os.Stat(targetPath); err != nil {
			logrus.Error("Cannot find " + targetPath)
			vm.Meta.Target = ""
		}
	}

	EnsureLineInGitIgnore(vm.Meta.GitRoot, ".gops-target.yaml")
}

func LoadAndMergeConfigDirectory(vm *ViewModel) {

	if vm.Meta.Target == "" {
		logrus.Error("No target has been set, please invoke 'gops target' first")
		os.Exit(1)
	}

	vm.Data = Dictionary{}

	configDirectory := filepath.Join(vm.Meta.GitRoot, vm.Config.ConfigDirectory)

	directories, err := ioutil.ReadDir(configDirectory)
	if err != nil {
		logrus.Error("Directory does not exist:  " + RelPath(vm.Meta.GitRoot, configDirectory))
		os.Exit(1)
	}

	for _, directory := range directories {
		if directory.IsDir() {

			files, _ := filepath.Glob(filepath.Join(configDirectory, directory.Name(), "*.yaml"))
			merged := LoadYamlFiles(files...)

			vm.Data[directory.Name()] = merged
		}
	}

	if _, ok := vm.Data["target"]; ok {
		logrus.Error("There is a config directory called target, bailing")
		os.Exit(1)
	}

	if val, ok := vm.Data[vm.Meta.Target]; ok {
		vm.Data["target"] = val
	} else {
		logrus.Error("Cannot find target " + vm.Meta.Target + " in merged configuration")
		os.Exit(1)
	}

	logrus.Info("Current target: " + vm.Meta.Target)
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
