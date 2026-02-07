package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yansongwel/kubeops/backend/internal/service"
)

// PodHandler Pod HTTP处理层
// 类比Shell函数：list_pods_handler() { result=$(list_pods); echo "$result"; }
type PodHandler struct {
	podService *service.PodService
}

// NewPodHandler 创建Pod Handler
func NewPodHandler(svc *service.PodService) *PodHandler {
	return &PodHandler{
		podService: svc,
	}
}

// ListPods 处理 GET /api/v1/namespaces/:namespace/pods 请求
// 对应Shell: case "pods" list_pods_handler ;;
func (h *PodHandler) ListPods(c *gin.Context) {
	// 从URL参数中提取namespace
	namespace := c.Param("namespace")

	// 调用Service层获取数据
	pods, err := h.podService.ListPodsInNamespace(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list pods",
			"details": err.Error(),
		})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"data": pods,
		"namespace": namespace,
	})
}

// GetPod 处理 GET /api/v1/namespaces/:namespace/pods/:name 请求
func (h *PodHandler) GetPod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	// 调用Service层获取数据
	pod, err := h.podService.GetPod(namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Pod not found",
			"details": err.Error(),
		})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"data": pod,
		"namespace": namespace,
	})
}

// ListAllPods 处理 GET /api/v1/pods 请求（获取所有命名空间的Pod）
func (h *PodHandler) ListAllPods(c *gin.Context) {
	// 调用Service层获取数据
	pods, err := h.podService.ListAllPods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list all pods",
			"details": err.Error(),
		})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"data": pods,
	})
}
