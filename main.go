package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

func main() {
	var (
		format       string
		inputPath    string
		outputPath   string
		templatePath string
	)

	flag.StringVar(&format, "format", "json", "data format 'json' or 'yaml' (default json)")
	flag.StringVar(&format, "f", "json", "data format 'json' or 'yaml' (default json)")

	flag.StringVar(&inputPath, "in", "", "input file (default STDIN)")
	flag.StringVar(&inputPath, "i", "", "input file (default STDIN)")

	flag.StringVar(&outputPath, "out", "", "output file (default STDOUT)")
	flag.StringVar(&outputPath, "o", "", "output file (default STDOUT)")

	flag.StringVar(&templatePath, "template", "", "template file path (required)")
	flag.StringVar(&templatePath, "t", "", "template file path (required)")

	flag.Parse()

	if templatePath == "" {
		fmt.Fprintln(os.Stderr, "template file path is required.")
		flag.Usage()
		return
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tempate file error:", err)
		return
	}

	var inFile io.ReadCloser
	if inputPath == "" {
		inFile = os.Stdin
	} else {
		var err error
		inFile, err = os.Open(inputPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "input file error:", err)
			return
		}
		defer inFile.Close()
	}

	var outFile io.WriteCloser
	if outputPath == "" {
		outFile = os.Stdout
	} else {
		var err error
		outFile, err = os.Create(outputPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "output file error:", err)
			return
		}
		defer outFile.Close()
	}

	b, err := ioutil.ReadAll(inFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "input file error:", err)
		return
	}

	m := make(map[string]interface{})
	switch format {
	case "json":
		err = json.Unmarshal(b, &m)
	case "yaml":
		err = yaml.Unmarshal(b, &m)
	default:
		fmt.Fprintln(os.Stderr, "unknown data format")
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "input file error:", err)
		return
	}

	err = tmpl.Execute(outFile, m)
	if err != nil {
		fmt.Fprintln(os.Stderr, "output file error:", err)
	}
}
