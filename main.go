package main

import (
	"context"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
)

const MatchImage = "ghcr.io/shovanmaity/byoc-echo"

type Specification struct {
	TotalChaosDuration int `required:"true" split_words:"true"`
}

func main() {
	client, err := getKubeClient()
	if err != nil {
		logrus.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	spec := Specification{}
	if err := envconfig.Process("", &spec); err != nil {
		logrus.Fatalf("Failed to process environment variables: %v", err)
	}
	pods, err := client.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Fatalf("Failed to list pods: %v", err)
	}
	for _, pod := range pods.Items {
		for _, con := range pod.Spec.Containers {
			if strings.HasPrefix(con.Image, MatchImage) {
				logrus.Infof("Found matching container in pod %s/%s: %s", pod.Namespace, pod.Name, con.Image)
				podYaml, err := pod.Marshal()
				if err != nil {
					logrus.Errorf("Failed to marshal pod %s/%s: %v", pod.Namespace, pod.Name, err)
					return
				}
				logrus.Infof("Pod YAML: \n %s", string(podYaml))
				return
			}
		}
	}
	for i := 0; i < spec.TotalChaosDuration; i++ {
		logrus.Infof("Waiting for %d seconds before checking again...", spec.TotalChaosDuration-i)
		time.Sleep(1 * time.Second)
	}
}

func getKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
