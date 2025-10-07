package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"text/template"
)

type GithubResponse struct {
	res []map[string]interface{}
}

func GithubDownloader(dd string) (string, error) {
	resourceURI := "https://api.github.com/repos/GlennTatum/ozone/contents/server/api/manifests"

	res, err := http.Get(resourceURI)
	if err != nil {
		return "", err
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var gh GithubResponse
	err = json.Unmarshal(buf, &gh.res)
	if err != nil {
		return "", err
	}
	for _, info := range gh.res {
		for k, v := range info {
			t := reflect.TypeOf(v).String()
			if t == "string" {
				val := v.(string)
				if k == "download_url" {
					uri := strings.Split(val, "/")
					basename := uri[len(uri)-1]
					if basename == dd {
						return val, nil
					}
				}
			}
		}
	}
	return "", nil
}

type ResourceMeta struct {
	Name string
	Port int
}

func FetchTemplate(name string, opts string, port int) (*bytes.Buffer, error) {

	remote, err := GithubDownloader(name)
	if err != nil {
		return nil, err
	}
	r, err := http.Get(remote)
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
	meta := ResourceMeta{Name: opts, Port: port}
	b := new(bytes.Buffer)
	tmpl.Execute(b, meta)
	return b, err
}
