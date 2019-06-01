package funi

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRenderer(t *testing.T) {
	cases := []struct {
		t      string
		data   string
		expect string
	}{
		{ t: "key is {{.key}}", data: `{"key":"value"}`, expect: "key is value"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("t: %s, data: %s", c.t, c.data), func(t *testing.T) {
			r, _ := NewRenderer(c.t)

			in := bytes.NewReader([]byte(c.data))
			var buf bytes.Buffer

			err := r.ExecuteJSON(in, &buf)
			if err != nil {
				t.Error(err)
				return
			}

			act := buf.String()
			if c.expect != buf.String() {
				t.Errorf("act: %q expect: %q", act, c.expect)
			}
		})
	}
}
