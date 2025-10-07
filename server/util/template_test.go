package util_test

import (
	"fmt"
	"ozone/api"
	"ozone/util"
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFetchTemplate(t *testing.T) {
	app, err := api.NewApp()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	tmpl, err := util.FetchTemplate("deployment.yml", "example-uuid")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(tmpl.String())

	data, err := app.K.Unmarshal(tmpl.Bytes())
	if err != nil {
		t.Fatalf("%s", err.Error())

	}
	// create the resource
	err = app.K.CreateResource("lab", data, v1.CreateOptions{})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
