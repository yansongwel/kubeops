package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yansongwel/kubeops/backend/internal/service"
)

// NamespaceHandler 命名空间HTTP处理层
// 类比Shell函数：list_namespaces_handler() { result=$(list_namespaces); echo "$result"; }
type NamespaceHandler struct {
	namespaceService *service.NamespaceService
}

// NewNamespaceHandler 创建命名空间Handler
func NewNamespaceHandler(svc *service.NamespaceService) *NamespaceHandler {
	return &NamespaceHandler{
		namespaceService: svc,
	}
}

// ListNamespaces 处理 GET /api/v1/namespaces 请求
// 对应Shell: case "namespaces" list_namespaces_handler ;;
func (h *NamespaceHandler) ListNamespaces(c *gin.Context) {
	// 调用Service层获取数据
	namespaces, err := h.namespaceService.ListNamespaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list namespaces",
			"details": err.Error(),
		})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"data": namespaces,
	})
}

// GetNamespace 处理 GET /api/v1/namespaces/:name 请求
func (h *NamespaceHandler) GetNamespace(c *gin.Context) {
	name := c.Param("name")

	// 调用Service层获取数据
	namespace, err := h.namespaceService.GetNamespace(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Namespace not found",
			"details": err.Error(),
		})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"data": namespace,
	})
}
