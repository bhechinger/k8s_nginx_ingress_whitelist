package main

import (
	"context"
	"log"

	"github.com/bhechinger/k8s_nginx_ingress_whitelist/internal/cidr"
	"github.com/bhechinger/k8s_nginx_ingress_whitelist/internal/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	cfg, err := config.New("0.0.1", "K8s CronJob to update NGINX Ingress whitelist")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Maybe make this so you can run outside of the cluster. I'm not sure why anyone would want to do that, however.
	// k8sConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	// if err != nil {
	// 	panic(err.Error())
	// }

	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

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

	cidrs, err := cidr.GetCIDRList(cfg.SourceURIs)
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
