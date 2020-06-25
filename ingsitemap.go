package main

import (
    // "fmt"
    // "time"
    "flag"

    // "sort"
    // // "encoding/json"
    // "sync"

    // "io/ioutil"
    // "net/http"
	// "bytes"
	// "os"
    
    "k8s.io/client-go/rest"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/apis/meta/v1" 

    log "github.com/sirupsen/logrus"
)

var (
    //read kubeconfig
    // kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	kubeconfig = "/_ext/Development/Project/devcn.fun/g-dev2/fk-kubernetes-auto-ingress/kubeconfig-vm23.203"

)

func getDatas()([]*data){
	flag.Parse()

    var err error
    var config *rest.Config

    if kubeconfig != "" {
        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
    } else {
        //get config when running inside Kubernetes
        config, err = rest.InClusterConfig()
    }

    if err != nil {
        log.Errorln(err.Error())
        return nil
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Errorln(err.Error())
        return nil
	}	
	
	pages:=getIngs(clientset)
	return pages
}

func getIngs(clientset *kubernetes.Clientset)([]*data) {
	// pods, err := clientset.Core().Pods("").List(v1.ListOptions{})
	ings, err := clientset.ExtensionsV1beta1().Ingresses("").List(v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// var rows Rows
	var pages []*data
	for _, ing := range ings.Items {
		log.Info("Ing added: ", ing.Name, " > ", ing.Spec.Rules[0].Host)
		d:= new(data)
		d.Title=ing.Name
		d.Url=ing.Spec.Rules[0].Host+":31714"

		pages= append(pages, d)
	}

	return pages
}