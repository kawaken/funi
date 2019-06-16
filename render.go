package funi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
)

// Renderer renders data via template.
type Renderer struct {
	t *template.Template
}

// NewRenderer returns Renderer. If tmplText is invalid template, it returns error.
func NewRenderer(tmplText string) (*Renderer, error) {
	t, err := template.New("mail").Funcs(sprig.TxtFuncMap()).Parse(tmplText)
	if err != nil {
		return nil, err
	}
	return &Renderer{
		t: t,
	}, nil
}

// ExecuteJSON unmarshals data as json and applies a template to data.
func (ren *Renderer) ExecuteJSON(in io.Reader, out io.Writer) error {
	return ren.execute(in, out, json.Unmarshal)
}

// ExecuteYAML unmarshals data as yaml and applies a template to data.
func (ren *Renderer) ExecuteYAML(in io.Reader, out io.Writer) error {
	return ren.execute(in, out, yaml.Unmarshal)
}

// ExecuteObject applies a template to object.
func (ren *Renderer) ExecuteObject(obj interface{}, out io.Writer) error {
	return ren.t.Execute(out, obj)
}

func (ren *Renderer) execute(in io.Reader, out io.Writer, uf func([]byte, interface{}) error) error {
	body, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	var data interface{}
	err = uf(body, &data)
	if err != nil {
		return err
	}

	return ren.t.Execute(out, data)
}
