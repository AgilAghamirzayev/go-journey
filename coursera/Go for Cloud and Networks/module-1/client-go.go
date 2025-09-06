package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Load kubeconfig
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	fmt.Println("✅ Kubernetes client created successfully")

	// Perform CRUD operations
	createPod(clientset)
	readPod(clientset)
	updatePod(clientset)
	deletePod(clientset)
}

func createPod(clientset *kubernetes.Clientset) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:latest",
				},
			},
		},
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create pod: %v", err)
	}

	fmt.Printf("✅ Created Pod %q\n", result.Name)
}

func readPod(clientset *kubernetes.Clientset) {
	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), "my-pod", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to read pod: %v", err)
	}

	fmt.Printf("📌 Pod Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
}

func updatePod(clientset *kubernetes.Clientset) {
	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), "my-pod", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get pod for update: %v", err)
	}

	// Update container image
	pod.Spec.Containers[0].Image = "nginx:1.19"

	updatedPod, err := clientset.CoreV1().Pods("default").Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("Failed to update pod: %v", err)
	}

	fmt.Printf("✅ Updated Pod %q with new image %s\n", updatedPod.Name, updatedPod.Spec.Containers[0].Image)
}

func deletePod(clientset *kubernetes.Clientset) {
	err := clientset.CoreV1().Pods("default").Delete(context.TODO(), "my-pod", metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Failed to delete pod: %v", err)
	}

	fmt.Println("🗑️ Pod deleted successfully")
}
