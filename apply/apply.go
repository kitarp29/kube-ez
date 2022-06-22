package apply

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func Main(filename string) string {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print(err.Error())
		return (err.Error())
	}
	log.Printf("%q \n", string(b))

	kubeconfig := os.Getenv("KUBECONFIG")

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Print(err.Error())
		return (err.Error())
	}

	dd, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Print(err.Error())
		return (err.Error())
	}

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			log.Print(err.Error())
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			log.Print(err.Error())
			return (err.Error())
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Print(err.Error())
			return (err.Error())
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(c.Discovery())
		if err != nil {
			log.Print(err.Error())
			return (err.Error())
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Print(err.Error())
			return (err.Error())
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			log.Print(err.Error())
			return (err.Error())
		}
	}
	if err != io.EOF {
		log.Print(err.Error())
		return (err.Error())
	}
	return filename + " Applied!"
}
