package client

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"go.uber.org/zap"
)

// InitK8sClient 初始化 Kubernetes 客户端
func InitK8sClient(logger *zap.Logger) (*kubernetes.Clientset, error) {
	// 优先使用集群内配置
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		// 回退到 kubeconfig 文件
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			homeDir := os.Getenv("HOME")
			if homeDir == "" {
				homeDir = os.Getenv("USERPROFILE") // Windows
			}
			kubeconfig = homeDir + "/.kube/config"
		}

		logger.Info("Using kubeconfig file", zap.String("kubeconfig", kubeconfig))
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
		}
	}

	// 创建 K8s 客户端
	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return client, nil
}
