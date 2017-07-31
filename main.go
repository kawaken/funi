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

	flag.StringVar(&format, "f", "json", "data format 'json' or 'yaml'")
	flag.StringVar(&inputPath, "i", "", "input file (default STDIN)")
	flag.StringVar(&outputPath, "o", "", "output file (default STDOUT)")
	flag.StringVar(&templatePath, "t", "", "template file path (required)")

	flag.Parse()

	var tmpl *template.Template
	var err error

	switch templatePath {
	case "":
		if len(flag.Args()) == 0 {
			err = fmt.Errorf("pattern is required")
			break
		}
		tmpl, err = template.New("main").Parse(flag.Arg(0))
	default:
		tmpl, err = template.ParseFiles(templatePath)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "tempate error:", err)
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

	var data interface{}
	switch format {
	case "json":
		err = json.Unmarshal(b, &data)
	case "yaml":
		err = yaml.Unmarshal(b, &data)
	default:
		fmt.Fprintln(os.Stderr, "unknown data format")
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "input file error:", err)
		return
	}

	err = tmpl.Execute(outFile, data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "output file error:", err)
	}
}
