package main

import (
	"context"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg, err := getConfig("0.0.1", "K8s CronJob to update NGINX Ingress whitelist")
	if err != nil {
		log.Fatal(err)
	}

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err.Error())
	}

	// k8sConfig, err := rest.InClusterConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	k8sClientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatal(err)
	}

	configMap, err := k8sClientSet.CoreV1().ConfigMaps(cfg.ConfigMap.Namespace).Get(context.Background(), cfg.ConfigMap.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	newConfigMap := configMap.DeepCopy()
	data := newConfigMap.Data
	if data == nil {
		data = make(map[string]string)
	}

	cidrs, err := getCIDRList(cfg.SourceURIs)
	if err != nil {
		log.Fatal(err)
	}

	data[cfg.NginxConfigName] = cidrs

	newConfigMap.Data = data

	// if the lists are different update the ConfigMap
	if configMap.Data[cfg.NginxConfigName] != newConfigMap.Data[cfg.NginxConfigName] {
		log.Println("Updating ConfigMap")
		_, err = k8sClientSet.CoreV1().ConfigMaps(cfg.ConfigMap.Namespace).Update(context.Background(), newConfigMap, metav1.UpdateOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
}
