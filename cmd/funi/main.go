package main

import (
	"flag"
	"fmt"
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

func main() {
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

	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	opts := funi.Options{
		Format:         format,
		InputPath:      inputPath,
		OutputPath:     outputPath,
		TemplatePath:   templatePath,
		TemplateString: flag.Arg(0),
	}
	err = funi.Render(opts)
}
