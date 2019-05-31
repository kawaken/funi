# funi

`funi` is template renderer with json/yaml data.

## Install

```bash
go get github.com/kawaken/funi/cmd/funi
```

## Usage

A template is passed via argument.

```bash
$ echo '{"key":"value"}' | funi "key is {{.key}}"
key is value
```

Use template file.

```bash
$ echo "key is {{.key}}" > key.tmpl
$ echo '{"key":"value"}' | funi -t key.tmpl
key is value
```

## Template

Template syntax is from `text/template`.  
[template \- The Go Programming Language](https://golang.org/pkg/text/template/)

`funi` can handle multple template files. Use `-t` option multiple. Template are combined as one template.

```bash
$ echo "key: {{.key}}" > key2.tmpl 
$ echo '{"key":"value"}' | ./funi -t key.tmpl -t key2.tmpl
key is value
key: value
```

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
