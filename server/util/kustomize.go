package util

import (
	"context"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

/*
 * params:
 * - resource []bytes the yaml manifest of the resource that is to be kustomized
 *
 */
func (c *ResourceManagerClient) Kustomize(resource []byte) error {
	f := filesys.MakeFsInMemory()

	b, err := os.ReadFile("kustomization.yaml") // FIXME: read the file from a networked resource or support options for a typed kustomization using kyaml
	if err != nil {
		return err
	}
	f.WriteFile("/data/kustomization.yaml", b)
	uresource, err := c.Unmarshal(resource)
	if err != nil {
		return err
	}
	f.WriteFile("/data/resource.yaml", resource)
	kopts := krusty.MakeDefaultOptions()
	kustomizer := krusty.MakeKustomizer(kopts)
	resMap, err := kustomizer.Run(f, "./data")
	if err != nil {
		return err
	}
	// kustomizations
	kustomApply := make([]unstructured.Unstructured, len(resMap.Resources()))
	for _, res := range resMap.Resources() {
		y, err := res.AsYAML()
		if err != nil {
			return err
		}
		u, err := c.Unmarshal(y)
		if err != nil {
			return err
		}
		kustomApply = append(kustomApply, *u)
	}
	for _, kustom := range kustomApply {
		// kustomize may emit empty objects
		if len(kustom.Object) == 0 {
			continue
		}
		gvr, _ := c.GroupVersionResource(uresource)
		if gvr.Empty() {
			continue
		}
		_, err := c.Dynamic.Resource(gvr).
			Namespace(uresource.GetNamespace()).
			Apply(context.Background(), uresource.GetName(), &kustom, v1.ApplyOptions{Force: true, FieldManager: "botto"}) // field manager is a required api field to determine the client who edited the manifest
		if err != nil {
			return err
		}
	}
	return nil
}
