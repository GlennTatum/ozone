package util

import (
	"context"

	"go.yaml.in/yaml/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourceManagerClient struct {
	Clientset *kubernetes.Clientset
	Dynamic   *dynamic.DynamicClient
}

func (c *ResourceManagerClient) Unmarshal(in []byte) (*unstructured.Unstructured, error) {
	m := map[string]interface{}{}
	err := yaml.Unmarshal(in, m)
	if err != nil {
		return &unstructured.Unstructured{}, err
	}

	return &unstructured.Unstructured{Object: m}, nil
}

func NewInClusterResourceManagerClient() (*ResourceManagerClient, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	c := &ResourceManagerClient{
		Clientset: clientset,
		Dynamic:   dynamicClient,
	}
	return c, nil
}

func NewOutOfClusterResourceManagerClient(kubeconfig string) (*ResourceManagerClient, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	c := &ResourceManagerClient{
		Clientset: clientset,
		Dynamic:   dynamicClient,
	}
	return c, nil
}

func (c *ResourceManagerClient) GroupVersionResource(data *unstructured.Unstructured) (schema.GroupVersionResource, error) {

	groupresources, err := restmapper.GetAPIGroupResources(c.Clientset.DiscoveryClient)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(groupresources)
	mapping, err := mapper.RESTMappings(data.GroupVersionKind().GroupKind())
	if err != nil {
		return schema.GroupVersionResource{}, nil
	}
	if len(mapping) > 0 {
		return mapping[0].Resource, nil
	}
	return schema.GroupVersionResource{}, nil
}

func (c *ResourceManagerClient) CreateResource(namespace string, data *unstructured.Unstructured, opts v1.CreateOptions) error {
	gvr, err := c.GroupVersionResource(data)
	if err != nil {
		return err
	}
	res := c.Dynamic.Resource(gvr)
	_, err = res.Namespace(namespace).Create(
		context.Background(),
		data,
		opts,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *ResourceManagerClient) GetResource(namespace string, name string, data *unstructured.Unstructured, opts v1.GetOptions) error {
	gvr, err := c.GroupVersionResource(data)
	if err != nil {
		return err
	}
	res := c.Dynamic.Resource(gvr)
	_, err = res.Namespace(namespace).Get(
		context.Background(),
		name,
		opts,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *ResourceManagerClient) DeleteResource(namespace string, data *unstructured.Unstructured, opts v1.DeleteOptions) error {
	gvr, err := c.GroupVersionResource(data)
	if err != nil {
		return err
	}
	res := c.Dynamic.Resource(gvr)
	err = res.Namespace(namespace).Delete(
		context.Background(),
		data.GetName(),
		opts,
	)
	if err != nil {
		return err
	}
	return nil
}
