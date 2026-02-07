package repository

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodRepository Pod数据访问层
// 类比Shell函数：get_pods_in_namespace() { kubectl get pods -n $NAMESPACE ... }
type PodRepository struct {
	client *kubernetes.Clientset
}

// NewPodRepository 创建Pod Repository
func NewPodRepository(client *kubernetes.Clientset) *PodRepository {
	return &PodRepository{
		client: client,
	}
}

// ListByNamespace 获取指定命名空间的所有Pod
// 对应Shell: kubectl get pods -n $NAMESPACE -o json
func (r *PodRepository) ListByNamespace(namespace string) ([]corev1.Pod, error) {
	list, err := r.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}
	return list.Items, nil
}

// GetByName 获取指定命名空间中的某个Pod
// 对应Shell: kubectl get pod $NAME -n $NAMESPACE
func (r *PodRepository) GetByName(namespace, name string) (*corev1.Pod, error) {
	pod, err := r.client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s in namespace %s: %w", name, namespace, err)
	}
	return pod, nil
}

// ListAll 获取所有命名空间的所有Pod
// 对应Shell: kubectl get pods --all-namespaces -o json
func (r *PodRepository) ListAll() ([]corev1.Pod, error) {
	list, err := r.client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list all pods: %w", err)
	}
	return list.Items, nil
}
