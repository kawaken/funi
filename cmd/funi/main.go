package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kawaken/funi"
)

func main() {
	var (
		format       string
		inputPath    string
		outputPath   string
		templatePath string
	)

	flag.StringVar(&format, "f", "json", "data format 'json' or 'yaml'")
	flag.StringVar(&inputPath, "i", "", "input file (default STDIN)")
	flag.StringVar(&outputPath, "o", "", "output file (default STDOUT)")
	flag.StringVar(&templatePath, "t", "", "template file path (required)")

	flag.Parse()
	var err error

	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	if templatePath == "" && len(flag.Args()) == 0 {
		err = fmt.Errorf("pattern is required")
		return
	}

	opts := funi.Options{
		Format:         format,
		InputPath:      inputPath,
		OutputPath:     outputPath,
		TemplatePath:   templatePath,
		TemplateString: flag.Arg(0),
	}
	err = funi.Render(opts)
}
