package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	postgres *pgxpool.Pool
	redis    *redis.Client
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(postgres *pgxpool.Pool, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{postgres: postgres, redis: redisClient}
}

// Health 处理 GET /health 请求
func (h *HealthHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	status := "healthy"
	details := gin.H{}

	if h.postgres != nil {
		if err := h.postgres.Ping(ctx); err != nil {
			status = "unhealthy"
			details["postgres"] = "unreachable"
		} else {
			details["postgres"] = "ok"
		}
	}

	if h.redis != nil {
		if err := h.redis.Ping(ctx).Err(); err != nil {
			status = "unhealthy"
			details["redis"] = "unreachable"
		} else {
			details["redis"] = "ok"
		}
	}

	code := http.StatusOK
	if status != "healthy" {
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, gin.H{
		"status":  status,
		"service": "kubeops",
		"details": details,
	})
}
