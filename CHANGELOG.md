# KubeOps 更新日志

所有重要的项目变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 新增
- 支持 APISIX 作为 API 网关
- 支持 Higress 作为 API 网关
- 集成 Istio 服务网格
- 添加 mTLS 认证支持
- 添加流量管理（灰度发布、故障注入）
- 添加可观测性（分布式追踪、指标收集）
- 完整的中文文档
- API 网关配置示例（APISIX/Higress）
- Istio 配置示例（认证、流量管理）

### 变更
- 架构更新：使用 APISIX/Higress 替代传统 Ingress
- 组件更新：添加 Istio 作为服务网格层
- 文档中文化：所有文档和注释使用中文

## [0.1.0] - 2026-02-07

### 新增
- 项目初始化
- 微服务架构设计
- 后端服务框架（Go + Gin）
- 前端框架（Vue 3 + TypeScript）
- Helm Chart 配置
- Docker 构建配置
- 基础文档（README、架构、贡献指南）

### 规划
- API 网关完整实现
- Kube 管理器完整实现
- AI 巡检服务
- DevOps 集成
- 日志平台集成
- 监控平台集成

## [0.2.0] - 规划中

### 新增
- 完整的 K8s 资源管理
- AI 智能巡检
- Jenkins 集成
- ArgoCD 集成
- ELK/Loki 日志平台
- Prometheus/VictoriaMetrics 监控
- 多集群支持

### 优化
- 性能优化
- 安全加固
- 可观测性增强

---

## 版本说明

- **主版本号**：不兼容的 API 变更
- **次版本号**：向下兼容的功能新增
- **修订号**：向下兼容的问题修复
