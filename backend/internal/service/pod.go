package service

import (
	"github.com/yansongwel/kubeops/backend/internal/repository"
)

// PodService Pod业务逻辑层
// 类比Shell函数：list_pods() { pods=$(get_pods_in_namespace); format_output; }
type PodRepositoryInterface interface {
	ListByNamespace(namespace string) ([]interface{}, error)
}

type PodService struct {
	podRepo *repository.PodRepository
}

// NewPodService 创建Pod Service
func NewPodService(repo *repository.PodRepository) *PodService {
	return &PodService{
		podRepo: repo,
	}
}

// ListPodsInNamespace 获取指定命名空间中的Pod列表（带业务规则）
// 业务规则：可以添加过滤、排序等逻辑
// 对应Shell: get_pods_in_namespace | awk '{print $1}'
func (s *PodService) ListPodsInNamespace(namespace string) ([]string, error) {
	// 调用Repository层获取数据
	pods, err := s.podRepo.ListByNamespace(namespace)
	if err != nil {
		return nil, err
	}

	// 业务逻辑：提取Pod名称
	var result []string
	for _, pod := range pods {
		result = append(result, pod.Name)
	}

	return result, nil
}

// GetPod 获取单个Pod信息
func (s *PodService) GetPod(namespace, name string) (string, error) {
	_, err := s.podRepo.GetByName(namespace, name)
	if err != nil {
		return "", err
	}
	return name, nil
}

// ListAllPods 获取所有命名空间的Pod（带命名空间前缀）
func (s *PodService) ListAllPods() ([]map[string]string, error) {
	// 调用Repository层获取数据
	pods, err := s.podRepo.ListAll()
	if err != nil {
		return nil, err
	}

	// 业务逻辑：组合命名空间和Pod名称
	var result []map[string]string
	for _, pod := range pods {
		result = append(result, map[string]string{
			"namespace": pod.Namespace,
			"name":      pod.Name,
			"status":    string(pod.Status.Phase),
		})
	}

	return result, nil
}
