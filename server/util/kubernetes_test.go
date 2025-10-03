package util

import (
	"os"
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuilder(t *testing.T) {

	client, err := NewOutOfClusterResourceManagerClient("kubeconfig.yml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	b, err := os.ReadFile("deployment.yml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	resource, err := client.Unmarshal(b)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	// err = client.CreateResource("web", resource, v1.CreateOptions{})
	// if err != nil {
	// 	t.Fatalf("%s", err)
	// }
	err = client.DeleteResource("web", resource, v1.DeleteOptions{})
}

func TestKustomize(t *testing.T) {
	client, err := NewOutOfClusterResourceManagerClient("kubeconfig.yml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	b, err := os.ReadFile("deployment.yml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	err = client.Kustomize(b)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
