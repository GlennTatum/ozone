package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

type ResourceMeta struct {
	Metadata string
}

func FetchTemplate(name string, opts string) (*bytes.Buffer, error) {

	resource_uri := "http://localhost:8001" // TODO switch to env profile
	r, err := http.Get(fmt.Sprintf("%s/%s", resource_uri, name))
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("resource").Parse(string(content))
	if err != nil {
		return nil, err
	}
	meta := ResourceMeta{Metadata: opts}
	b := new(bytes.Buffer)
	tmpl.Execute(b, meta)
	return b, err
}
