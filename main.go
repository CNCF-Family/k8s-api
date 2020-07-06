package main

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

func main() {
	router := gin.Default()

	router.GET("/verify/deploy", func(c *gin.Context) {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}
		deployments, err := clientset.AppsV1().Deployments("gott").List(metav1.ListOptions{})
		if err != nil {
			log.Fatalln("failed to get deployments:", err)
		}
		error_pod := ""
		for _, deploy := range deployments.Items {
			if deploy.Status.UpdatedReplicas != *(deploy.Spec.Replicas) ||
				deploy.Status.AvailableReplicas != *(deploy.Spec.Replicas) {
				error_pod = deploy.GetName()
				break;
			}
		}
		if error_pod == "" {
			c.String(http.StatusOK, "success")
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error_pod": error_pod})
		}
	})
	router.Run(":9999")

}