package repository

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NamespaceRepository 命名空间数据访问层
// 类比Shell函数：get_all_namespaces() { kubectl get namespaces ... }
type NamespaceRepository struct {
	client *kubernetes.Clientset
}

// NewNamespaceRepository 创建命名空间Repository
func NewNamespaceRepository(client *kubernetes.Clientset) *NamespaceRepository {
	return &NamespaceRepository{
		client: client,
	}
}

// ListAll 获取所有命名空间
// 对应Shell: kubectl get namespaces -o json
func (r *NamespaceRepository) ListAll() ([]corev1.Namespace, error) {
	list, err := r.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	return list.Items, nil
}

// GetByName 根据名称获取命名空间
// 对应Shell: kubectl get namespace $NAME
func (r *NamespaceRepository) GetByName(name string) (*corev1.Namespace, error) {
	ns, err := r.client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace %s: %w", name, err)
	}
	return ns, nil
}
