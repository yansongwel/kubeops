package service

import (
	"github.com/yansongwel/kubeops/backend/internal/repository"
)

// NamespaceService 命名空间业务逻辑层
// 类比Shell函数：list_namespaces() { all=$(get_all_namespaces); filter; echo; }
type NamespaceService struct {
	namespaceRepo *repository.NamespaceRepository
}

// NewNamespaceService 创建命名空间Service
func NewNamespaceService(repo *repository.NamespaceRepository) *NamespaceService {
	return &NamespaceService{
		namespaceRepo: repo,
	}
}

// ListNamespaces 获取命名空间列表（带业务规则过滤）
// 业务规则：过滤系统命名空间（kube-system, kube-public, kube-node-lease）
// 对应Shell: get_all_namespaces | grep -v "kube-system" | grep -v "kube-public"
func (s *NamespaceService) ListNamespaces() ([]string, error) {
	// 调用Repository层获取数据
	allNamespaces, err := s.namespaceRepo.ListAll()
	if err != nil {
		return nil, err
	}

	// 业务逻辑：过滤系统命名空间
	var result []string
	for _, ns := range allNamespaces {
		if !isSystemNamespace(ns.Name) {
			result = append(result, ns.Name)
		}
	}

	return result, nil
}

// GetNamespace 获取单个命名空间信息
func (s *NamespaceService) GetNamespace(name string) (string, error) {
	_, err := s.namespaceRepo.GetByName(name)
	if err != nil {
		return "", err
	}
	return name, nil
}

// isSystemNamespace 检查是否为系统命名空间
// 对应Shell: if [[ "$ns" == "kube-system" ]] || [[ "$ns" == "kube-public" ]]; then
func isSystemNamespace(name string) bool {
	systemNamespaces := []string{
		"kube-system",
		"kube-public",
		"kube-node-lease",
	}

	for _, ns := range systemNamespaces {
		if name == ns {
			return true
		}
	}
	return false
}
