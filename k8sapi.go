package main

import (
	"fmt"
	"github.com/ebar-go/ego/http"
	"github.com/ebar-go/ego/utils"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"github.com/ebar-go/ego/http/response"
)

func main() {
	server := http.NewServer()
	// 添加路由
	server.Router.GET("/k8s/getnodes", func(context *gin.Context) {
		servicename := context.Query("servicename")

		ns := context.Query("namespace")
		if (ns == "") {
			ns = "gott"
		}

		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}
		labelPod := labels.SelectorFromSet(labels.Set(map[string]string{"app": servicename}))
		listPodOptions := metav1.ListOptions{
			LabelSelector: labelPod.String(),
		}
		list, _ := clientset.CoreV1().Pods(ns).List(listPodOptions)
		var nodeSlice []string
		for i, node := range list.Items {
			nodeSlice = append(nodeSlice, node.Spec.NodeName)
			fmt.Printf("[%d] %s\n", i, node.Spec.NodeName)
		}
		fmt.Printf("%v\n", nodeSlice)
		response.WrapContext(context).Success(response.Data{"node":nodeSlice})
	})
	utils.FatalError("StartServer", server.Start())
}