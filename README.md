# funi

`funi` is template renderer with json/yaml data.

## Install

```bash
go get github.com/kawaken/funi/cmd/funi
```

## Usage

```bash
$ echo '{"key":"value"}' | funi "key is {{.key}}"
key is value
```

Template syntax is from `text/template`.  
[template \- The Go Programming Language](https://golang.org/pkg/text/template/)

### Options

```bash
$ funi -h
Usage of funi:
  -f string
    	data format 'json' or 'yaml' (default "json")
  -i string
    	input file (default STDIN)
  -o string
    	output file (default STDOUT)
  -t string
    	template file path (required)
```

`-t` is required when command line argument is nothing.
