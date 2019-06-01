package funi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && !info.IsDir()
}

// Options is option for rendering
type Options struct {
	Format         string
	InputPath      string
	OutputPath     string
	TemplatePath   []string
	TemplateString string
}

func loadTemplate(opts Options) (*template.Template, error) {
	var buf bytes.Buffer

	if len(opts.TemplateString) > 0 {
		_, err := buf.WriteString(opts.TemplateString)
		if err != nil {
			return nil, err
		}
		return template.New("main").Parse(buf.String())
	}

	for _, path := range opts.TemplatePath {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}

	return template.New("main").Parse(buf.String())
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

func validate(opts Options) error {
	// validation
	if len(opts.TemplateString) == 0 && len(opts.TemplatePath) == 0 {
		return fmt.Errorf("template is required")
	}

	for _, p := range opts.TemplatePath {
		if len(p) == 0 {
			return fmt.Errorf("empty path is not available")
		}

		if !fileExists(p) {
			return fmt.Errorf("%s does not exist or is a directory", p)
		}
	}

	if len(opts.InputPath) > 0 && !fileExists(opts.InputPath) {
		return fmt.Errorf("%s does not exist or is a directory", opts.InputPath)
	}

	if len(opts.OutputPath) > 0 {
		info, err := os.Stat(opts.OutputPath)
		if err == nil {
			if info.IsDir() {
				return fmt.Errorf("%s is a directory", opts.OutputPath)
			}
		} else {
			if os.IsNotExist(err) {
				// output file is created if not exists
			} else {
				return err
			}
		}
	}

	switch opts.Format {
	case "json", "yaml":
	default:
		return fmt.Errorf("unknown data format. use 'json' or 'yaml")
	}

	return nil
}

// Render renders template
func Render(opts Options) error {

	err := validate(opts)
	if err != nil {
		return err
	}

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
	var unmarshalFunc func([]byte, interface{}) error
	switch opts.Format {
	case "json":
		unmarshalFunc = json.Unmarshal
	case "yaml":
		unmarshalFunc = yaml.Unmarshal
	}
	err = unmarshalFunc(b, &data)
	if err != nil {
		return err
	}

	err = tmpl.Execute(outFile, data)

	return err
}
