package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/yansongwel/kubeops/backend/internal/handler"
	"github.com/yansongwel/kubeops/backend/internal/repository"
	"github.com/yansongwel/kubeops/backend/internal/service"
)

func main() {
	// 1. 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("KubeOps starting...")

	// 2. 初始化 K8s 客户端
	k8sClient, err := initK8sClient(logger)
	if err != nil {
		logger.Fatal("Failed to initialize Kubernetes client", zap.Error(err))
	}
	logger.Info("Successfully connected to Kubernetes cluster")

	// 3. 初始化 Repository 层
	namespaceRepo := repository.NewNamespaceRepository(k8sClient)
	podRepo := repository.NewPodRepository(k8sClient)

	// 4. 初始化 Service 层
	namespaceService := service.NewNamespaceService(namespaceRepo)
	podService := service.NewPodService(podRepo)

	// 5. 初始化 Handler 层
	namespaceHandler := handler.NewNamespaceHandler(namespaceService)
	podHandler := handler.NewPodHandler(podService)

	// 6. 配置路由
	router := setupRouter(namespaceHandler, podHandler, logger)

	// 7. 启动 HTTP 服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting KubeOps server",
			zap.String("port", port),
			zap.String("environment", os.Getenv("ENV")),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// 8. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// initK8sClient 初始化 Kubernetes 客户端
func initK8sClient(logger *zap.Logger) (*kubernetes.Clientset, error) {
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

// setupRouter 配置 HTTP 路由
func setupRouter(
	namespaceHandler *handler.NamespaceHandler,
	podHandler *handler.PodHandler,
	logger *zap.Logger,
) *gin.Engine {
	// 设置 Gin 模式
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "kubeops",
		})
	})

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 测试端点
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		// 命名空间相关路由
		v1.GET("/namespaces", namespaceHandler.ListNamespaces)
		v1.GET("/namespaces/:name", namespaceHandler.GetNamespace)

		// Pod 相关路由
		v1.GET("/namespaces/:namespace/pods", podHandler.ListPods)
		v1.GET("/namespaces/:namespace/pods/:name", podHandler.GetPod)
		v1.GET("/pods", podHandler.ListAllPods)
	}

	logger.Info("Routes registered successfully")
	return router
}
