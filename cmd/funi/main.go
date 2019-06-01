package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kawaken/funi"
)

type paths []string

func (p *paths) String() string {
	return strings.Join([]string(*p), ", ")
}

func (p *paths) Set(s string) error {
	*p = append(*p, s)
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && !info.IsDir()
}

func validate(opts funi.Options) error {
	// validation
	if len(opts.TemplateString) == 0 {
		return fmt.Errorf("template is required")
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

func loadTemplate(templatePath paths) (string, error) {
	var buf bytes.Buffer

	for _, path := range templatePath {
		if len(path) == 0 {
			return "", fmt.Errorf("empty path is not available")
		}

		if !fileExists(path) {
			return "", fmt.Errorf("%s does not exist or is a directory", path)
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		_, err = buf.Write(b)
		if err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func run() error {
	var (
		format       string
		inputPath    string
		outputPath   string
		templatePath paths
	)

	flag.StringVar(&format, "f", "json", "data format 'json' or 'yaml'")
	flag.StringVar(&inputPath, "i", "", "input file (default STDIN)")
	flag.StringVar(&outputPath, "o", "", "output file (default STDOUT)")
	flag.Var(&templatePath, "t", "template file path (required if no arg), can use multiple")

	flag.Parse()
	var err error

	templateString := flag.Arg(0)
	if len(templateString) == 0 {
		if len(templatePath) == 0 {
			return fmt.Errorf("template is required")
		}
		templateString, err = loadTemplate(templatePath)
		if err != nil {
			return err
		}
	}

	opts := funi.Options{
		Format:         format,
		InputPath:      inputPath,
		OutputPath:     outputPath,
		TemplateString: templateString,
	}

	return funi.Render(opts)
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
