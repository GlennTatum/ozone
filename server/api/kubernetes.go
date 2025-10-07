package api

import (
	"context"
	"fmt"
	"ozone/util"

	networkv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
Patches an additional user service onto an ozone ingress controller
*/
func (app *App) OzoneIngressUserResourcePatch(namespace string) error {
	ingress, err := app.K.Clientset.NetworkingV1().Ingresses(namespace).Get(
		context.Background(),
		"ozone-ingress",
		v1.GetOptions{},
	)
	currentPaths := ingress.Spec.Rules[0].HTTP.Paths
	currentPaths = append(currentPaths, networkv1.HTTPIngressPath{})

	if err != nil {
		return err
	}

	return nil
}

func (app *App) CreateResourceFromTeplate(path string, namespace string, resource_id string, port int) error {
	// get the yaml manifest as a template
	// the Account table will now have a resource associated with it in the kubernetes cluster
	tmpl, err := util.FetchTemplate(path, resource_id, port)
	if err != nil {
		return err
	}
	data, err := app.K.Unmarshal(tmpl.Bytes())
	if err != nil {
		return err
	}
	// create the resource
	fmt.Println(data)

	// TODO set environment profile for default kubernetes namespace
	err = app.K.CreateResource(namespace, data, v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
