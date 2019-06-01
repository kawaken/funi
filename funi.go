package funi

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Options is option for rendering
type Options struct {
	Format         string
	InputPath      string
	OutputPath     string
	TemplateString string
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

	renderer, err := NewRenderer(opts.TemplateString)
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

	var unmarshalFunc func([]byte, interface{}) error
	switch opts.Format {
	case "json":
		unmarshalFunc = json.Unmarshal
	case "yaml":
		unmarshalFunc = yaml.Unmarshal
	}

	return renderer.execute(inFile, outFile, unmarshalFunc)
}
