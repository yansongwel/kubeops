package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yansongwel/kubeops/backend/internal/client"
	"github.com/yansongwel/kubeops/backend/internal/config"
	"github.com/yansongwel/kubeops/backend/internal/handler"
	"github.com/yansongwel/kubeops/backend/internal/repository"
	"github.com/yansongwel/kubeops/backend/internal/service"
)

func main() {
	// 1. 初始化日志
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()
	logger.Info("KubeOps starting...")

	cfg := loadConfig()

	// 2. 初始化 K8s 客户端
	k8sClient, err := client.InitK8sClient(logger, cfg.Kubeconfig)
	if err != nil {
		logger.Fatal("Failed to initialize Kubernetes client", zap.Error(err))
	}
	logger.Info("Successfully connected to Kubernetes cluster")

	postgresPool, err := client.NewPostgresPool(cfg.Postgres)
	if err != nil {
		logger.Fatal("Failed to initialize Postgres", zap.Error(err))
	}
	defer postgresPool.Close()

	redisClient, err := client.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	defer func() {
		_ = redisClient.Close()
	}()

	// 3. 初始化 Repository 层
	namespaceRepo := repository.NewNamespaceRepository(k8sClient)
	podRepo := repository.NewPodRepository(k8sClient)

	// 4. 初始化 Service 层
	namespaceService := service.NewNamespaceService(namespaceRepo)
	podService := service.NewPodService(podRepo)

	// 5. 初始化 Handler 层
	namespaceHandler := handler.NewNamespaceHandler(namespaceService)
	podHandler := handler.NewPodHandler(podService)
	healthHandler := handler.NewHealthHandler(postgresPool, redisClient)

	// 6. 配置路由
	router := setupRouter(namespaceHandler, podHandler, healthHandler, cfg.Env, logger)

	// 7. 启动 HTTP 服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting KubeOps server",
			zap.String("port", cfg.Port),
			zap.String("environment", cfg.Env),
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

func setupRouter(
	namespaceHandler *handler.NamespaceHandler,
	podHandler *handler.PodHandler,
	healthHandler *handler.HealthHandler,
	env string,
	logger *zap.Logger,
) *gin.Engine {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/health", healthHandler.Health)

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

func loadConfig() config.Config {
	cfg := config.Load()

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	fs.StringVar(&cfg.Port, "port", cfg.Port, "服务端口")
	fs.StringVar(&cfg.Env, "env", cfg.Env, "运行环境")
	fs.StringVar(&cfg.Kubeconfig, "kubeconfig", cfg.Kubeconfig, "kubeconfig 文件路径")
	fs.StringVar(&cfg.Postgres.Host, "postgres-host", cfg.Postgres.Host, "PostgreSQL 地址")
	fs.StringVar(&cfg.Postgres.Port, "postgres-port", cfg.Postgres.Port, "PostgreSQL 端口")
	fs.StringVar(&cfg.Postgres.User, "postgres-user", cfg.Postgres.User, "PostgreSQL 用户")
	fs.StringVar(&cfg.Postgres.Password, "postgres-password", cfg.Postgres.Password, "PostgreSQL 密码")
	fs.StringVar(&cfg.Postgres.Database, "postgres-db", cfg.Postgres.Database, "PostgreSQL 数据库")
	fs.StringVar(&cfg.Postgres.SSLMode, "postgres-sslmode", cfg.Postgres.SSLMode, "PostgreSQL SSL 模式")
	fs.StringVar(&cfg.Redis.Addr, "redis-addr", cfg.Redis.Addr, "Redis 地址")
	fs.StringVar(&cfg.Redis.Password, "redis-password", cfg.Redis.Password, "Redis 密码")
	fs.IntVar(&cfg.Redis.DB, "redis-db", cfg.Redis.DB, "Redis DB 编号")

	fs.Usage = func() {
		_, _ = fmt.Fprintln(os.Stdout, "KubeOps 后端服务")
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "用法:")
		_, _ = fmt.Fprintf(os.Stdout, "  %s [options]\n", os.Args[0])
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "选项:")
		fs.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "示例:")
		_, _ = fmt.Fprintf(os.Stdout, "  %s --postgres-host 192.168.33.100 --postgres-port 5432 --postgres-user kubeops \\\n", os.Args[0])
		_, _ = fmt.Fprintln(os.Stdout, "    --postgres-password kubeops --postgres-db kubeops --redis-addr 192.168.33.100:6379")
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	missing := requiredConfigMissing(cfg)
	if len(missing) > 0 {
		_, _ = fmt.Fprintln(os.Stdout, "缺少必要参数:", strings.Join(missing, ", "))
		fs.Usage()
		os.Exit(2)
	}

	return cfg
}

func requiredConfigMissing(cfg config.Config) []string {
	var missing []string
	if cfg.Postgres.Host == "" {
		missing = append(missing, "--postgres-host/POSTGRES_HOST")
	}
	if cfg.Postgres.Port == "" {
		missing = append(missing, "--postgres-port/POSTGRES_PORT")
	}
	if cfg.Postgres.User == "" {
		missing = append(missing, "--postgres-user/POSTGRES_USER")
	}
	if cfg.Postgres.Password == "" {
		missing = append(missing, "--postgres-password/POSTGRES_PASSWORD")
	}
	if cfg.Postgres.Database == "" {
		missing = append(missing, "--postgres-db/POSTGRES_DB")
	}
	if cfg.Redis.Addr == "" {
		missing = append(missing, "--redis-addr/REDIS_ADDR")
	}
	return missing
}
