package funi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Options is option for rendering
type Options struct {
	Format         string
	InputPath      string
	OutputPath     string
	TemplatePath   []string
	TemplateString string
}

func loadTemplate(opts Options) (*template.Template, error) {
	var tmpl *template.Template
	var err error

	switch len(opts.TemplatePath) {
	case 0:
		tmpl, err = template.New("main").Parse(opts.TemplateString)
	default:
		tmpl, err = template.ParseFiles(opts.TemplatePath...)
	}

	return tmpl, err
}

func loadInputFile(path string, r io.ReadCloser) (io.ReadCloser, error) {
	var inFile io.ReadCloser
	var err error

	if path == "" {
		inFile = r
	} else {
		inFile, err = os.Open(path)
	}

	return inFile, err
}

func openOutputFile(path string, w io.WriteCloser) (io.WriteCloser, error) {
	var outFile io.WriteCloser
	var err error

	if path == "" {
		outFile = w
	} else {
		outFile, err = os.Create(path)
	}

	return outFile, err
}

// Render renders template
func Render(opts Options) error {

	tmpl, err := loadTemplate(opts)
	if err != nil {
		return err
	}

	inFile, err := loadInputFile(opts.InputPath, os.Stdin)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := openOutputFile(opts.OutputPath, os.Stdout)
	if err != nil {
		return err
	}
	defer outFile.Close()

	b, err := ioutil.ReadAll(inFile)
	if err != nil {
		return err
	}

	var data interface{}
	switch opts.Format {
	case "json":
		err = json.Unmarshal(b, &data)
	case "yaml":
		err = yaml.Unmarshal(b, &data)
	default:
		err = fmt.Errorf("unknown data format")
	}
	if err != nil {
		return err
	}

	err = tmpl.Execute(outFile, data)

	return err
}
