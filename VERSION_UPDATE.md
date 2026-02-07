# KubeOps 版本更新总结

**更新日期**: 2026-02-07
**更新内容**: 统一版本号和仓库地址

## 版本更新清单

### 后端 (Go)

| 组件 | 旧版本 | 新版本 | 说明 |
|------|--------|--------|------|
| Go | 1.21+ | **1.25** | K8s 1.33需要Go 1.25+ |
| Gin | v1.9.1 | **v1.11.0** | 最新稳定版 |
| gRPC | v1.59.0 | **v1.70.0** | 最新稳定版 |
| Protobuf | v1.31.0 | **v1.36.1** | 最新稳定版 |
| K8s API | v0.28.3 | **v0.35.0** | 对应K8s 1.33 |
| K8s client-go | v0.28.3 | **v0.35.0** | 对应K8s 1.33 |
| controller-runtime | v0.16.3 | **v0.20.0** | 最新稳定版 |
| JWT | v5.2.0 | **v5.2.2** | 补丁更新 |
| Redis | v9.3.0 | **v9.7.0** | 最新稳定版 |
| GORM | v1.25.5 | **v1.25.13** | 补丁更新 |
| Prometheus | v1.17.0 | **v1.20.0** | 最新稳定版 |
| Zap (日志) | v1.26.0 | **v1.27.0** | 最新稳定版 |
| Viper | v1.17.0 | **v1.19.0** | 最新稳定版 |

**Kubernetes版本**: 1.33 (client-go v0.35.0)

### 前端 (Node.js)

| 组件 | 旧版本 | 新版本 | 说明 |
|------|--------|--------|------|
| Node.js | 18+ | **24.0+** | LTS最新版 |
| Vue | ^3.4.0 | **^3.5.0** | 最新稳定版 |
| Vue Router | ^4.2.5 | **^4.5.0** | 最新稳定版 |
| Pinia | ^2.1.7 | **^2.2.0** | 最新稳定版 |
| Axios | ^1.6.2 | **^1.7.0** | 最新稳定版 |
| Element Plus | ^2.5.0 | **^2.9.0** | 最新稳定版 |
| ECharts | ^5.4.3 | **^5.6.0** | 最新稳定版 |
| vue-echarts | ^6.6.4 | **^7.0.0** | 主版本升级 |
| TypeScript | ^5.3.0 | **^5.7.0** | 最新稳定版 |
| Vite | ^5.0.0 | **^6.0.0** | 主版本升级 |
| Vitest | ^1.1.0 | **^3.0.0** | 主版本升级 |
| ESLint | ^8.55.0 | **^9.0.0** | 主版本升级 |

## 仓库地址更新

**旧地址**: `github.com/your-org/kubeops`
**新地址**: `https://github.com/yansongwel/kubeops.git`

### 更新的文件

1. ✅ `backend/kube-manager/go.mod` - 模块名和依赖版本
2. ✅ `backend/kube-manager/cmd/server/main.go` - 导入路径
3. ✅ `backend/kube-manager/internal/handler/*.go` - 导入路径
4. ✅ `backend/kube-manager/internal/service/*.go` - 导入路径
5. ✅ `go.mod` (根目录) - 主模块配置
6. ✅ `frontend/package.json` - Node版本和依赖
7. ✅ `CLAUDE.md` - 技术栈版本说明
8. ✅ `README.md` - 前置要求
9. ✅ `PROJECT_STATUS.md` - 技术栈和仓库地址

## 代码重构完成

同时完成了kube-manager的代码重构，实现分层架构：

### 新增文件
- ✅ `backend/kube-manager/internal/repository/namespace.go`
- ✅ `backend/kube-manager/internal/repository/pod.go`
- ✅ `backend/kube-manager/internal/service/namespace.go`
- ✅ `backend/kube-manager/internal/service/pod.go`
- ✅ `backend/kube-manager/internal/handler/namespace.go`
- ✅ `backend/kube-manager/internal/handler/pod.go`

### 重构的文件
- ✅ `backend/kube-manager/cmd/server/main.go` - 使用分层架构

### 代码结构
```
Request → Handler → Service → Repository → K8s API
  ↓         ↓          ↓            ↓          ↓
HTTP     参数验证    业务逻辑    数据访问    K8s集群
```

## 编译验证

✅ kube-manager服务编译成功
- 生成文件: `backend/kube-manager/server.exe` (46MB)
- 所有依赖正确导入
- 模块路径正确: `github.com/yansongwel/kubeops/kube-manager`

## 兼容性说明

### Go版本
- K8s 1.33 (client-go v0.35.0) 要求 Go 1.25+
- 已自动升级到 Go 1.25.0

### Node版本
- 前端要求 Node.js 24.0+
- Vite 6.0 需要较新的Node版本

### Kubernetes版本
- 支持的K8s版本: 1.33.x
- client-go会向后兼容，但建议使用同版本

## 下一步建议

1. **测试启动**: 运行 `server.exe` 测试服务启动
2. **API测试**: 测试命名空间和Pod列表API
3. **前端测试**: 运行 `npm install` 和 `npm run dev`
4. **依赖更新**: 其他服务同样更新版本号

## 命令参考

### 后端
```bash
# 编译
cd backend/kube-manager
go build -o server.exe ./cmd/server

# 运行
./server.exe

# 或使用go run
go run cmd/server/main.go
```

### 前端
```bash
cd frontend
npm install
npm run dev
```

### Docker
```bash
# 构建后端
docker build --build-arg SERVICE=kube-manager -t kubeops/kube-manager:latest .

# 构建前端
docker build -t kubeops/frontend:latest ./frontend
```

---
**更新完成时间**: 2026-02-07 16:00
**更新人**: Claude Code
**状态**: ✅ 全部完成
