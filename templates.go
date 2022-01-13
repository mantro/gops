package main

import (
	"bytes"
	"github.com/CloudyKit/jet/v6"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func TemplateFile(views *jet.Set, file string, vm ViewModel) {

	relativeFile := RelPath(vm.Meta.GitRoot, file)
	logrus.Info("Templating " + relativeFile)

	info, err := os.Stat(file)
	if err != nil {
		logrus.Error("File does not exist: " + relativeFile)
	}

	tmpl, err := views.GetTemplate(file)
	if err != nil {
		logrus.Error("Unable to parse template file " + relativeFile)
		panic(err)
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, nil, vm)
	if err != nil {
		logrus.Error("Cannot execute template " + relativeFile)
		panic(err)
	}

	outputFile := strings.Replace(file, ".template.", ".generated.", 1)

	err = os.WriteFile(outputFile, buf.Bytes(), info.Mode())
}
