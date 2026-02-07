## 前端常用组合式函数

### useTable - 表格管理
```typescript
// src/composables/useTable.ts
import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';
import type { TableProps } from 'element-plus';

export function useTable<T>(apiFn: (params: any) => Promise<any>) {
  const data = ref<T[]>([]);
  const loading = ref(false);
  const total = ref(0);

  // 分页
  const pagination = reactive({
    page: 1,
    pageSize: 20,
  });

  // 加载数据
  const load = async () => {
    loading.value = true;
    try {
      const res = await apiFn({
        page: pagination.page,
        pageSize: pagination.pageSize,
      });
      data.value = res.data || [];
      total.value = res.total || 0;
    } catch (error) {
      ElMessage.error('加载数据失败');
    } finally {
      loading.value = false;
    }
  };

  // 刷新
  const refresh = () => {
    pagination.page = 1;
    load();
  };

  return {
    data,
    loading,
    total,
    pagination,
    load,
    refresh,
  };
}
```

### useDialog - 对话框管理
```typescript
// src/composables/useDialog.ts
import { ref } from 'vue';

export function useDialog<T = any>() {
  const visible = ref(false);
  const data = ref<T>({} as T);
  const loading = ref(false);

  const open = (item?: T) => {
    data.value = item || ({} as T);
    visible.value = true;
  };

  const close = () => {
    visible.value = false;
    data.value = {} as T;
  };

  const toggle = (value: boolean) => {
    visible.value = value;
  };

  return {
    visible,
    data,
    loading,
    open,
    close,
    toggle,
  };
}
```

### useWebSocket - WebSocket 连接
```typescript
// src/composables/useWebSocket.ts
import { ref, onUnmounted } from 'vue';
import { ElMessage } from 'element-plus';

export function useWebSocket(url: string) {
  const ws = ref<WebSocket | null>(null);
  const connected = ref(false);
  const message = ref<string>('');
  const messages = ref<string[]>([]);

  const connect = () => {
    ws.value = new WebSocket(url);

    ws.value.onopen = () => {
      connected.value = true;
      ElMessage.success('WebSocket 连接成功');
    };

    ws.value.onmessage = (event) => {
      messages.value.push(event.data);
      message.value = event.data;
    };

    ws.value.onerror = (error) => {
      ElMessage.error('WebSocket 错误');
      console.error('WebSocket error:', error);
    };

    ws.value.onclose = () => {
      connected.value = false;
      ElMessage.warning('WebSocket 连接关闭');
    };
  };

  const send = (data: string | object) => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(typeof data === 'string' ? data : JSON.stringify(data));
    } else {
      ElMessage.error('WebSocket 未连接');
    }
  };

  const close = () => {
    ws.value?.close();
  };

  onUnmounted(() => {
    close();
  });

  return {
    connected,
    message,
    messages,
    connect,
    send,
    close,
  };
}
```

### useRequest - 请求封装
```typescript
// src/composables/useRequest.ts
import { ref } from 'vue';
import { ElMessage } from 'element-plus';

export function useRequest<T = any>(
  apiFn: (...args: any[]) => Promise<T>
) {
  const loading = ref(false);
  const error = ref<Error | null>(null);

  const execute = async (...args: any[]): Promise<T | null> => {
    loading.value = true;
    error.value = null;
    try {
      const result = await apiFn(...args);
      return result;
    } catch (e) {
      error.value = e as Error;
      ElMessage.error((e as Error).message);
      return null;
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    error,
    execute,
  };
}
```

## 表单验证工具

### 表单验证规则
```typescript
// src/utils/validate.ts
import type { FormRule } from 'element-plus';

// 必填
export const required = (message = '该字段为必填项'): FormRule => ({
  required: true,
  message,
  trigger: 'blur',
});

// 邮箱
export const email = (): FormRule => ({
  type: 'email',
  message: '请输入正确的邮箱地址',
  trigger: ['blur', 'change'],
});

// K8s 命名规范
export const k8sName = (): FormRule => ({
  pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
  message: '只能包含小写字母、数字和连字符，且必须以字母或数字开头和结尾',
  trigger: 'blur',
});

// URL
export const url = (): FormRule => ({
  type: 'url',
  message: '请输入正确的 URL',
  trigger: 'blur',
});

// 自定义正则
export const pattern = (
  regex: RegExp,
  message: string
): FormRule => ({
  pattern: regex,
  message,
  trigger: 'blur',
});

// 使用
const rules = {
  name: [required('请输入名称'), k8sName()],
  email: [required(), email()],
  replicas: [
    required(),
    pattern(/^[1-9]\d*$/, '副本数必须为正整数'),
  ],
};
```

## 格式化工具

### 格式化工具函数
```typescript
// src/utils/format.ts
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';

dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

// 时间格式化
export const formatTime = (time: string | Date, format = 'YYYY-MM-DD HH:mm:ss') => {
  return dayjs(time).format(format);
};

// 相对时间
export const formatRelativeTime = (time: string | Date) => {
  return dayjs(time).fromNow();
};

// 字节格式化
export const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
};

// CPU 格式化
export const formatCpu = (cpu: number): string => {
  if (cpu < 1) {
    return `${(cpu * 1000).toFixed(0)}m`;
  }
  return `${cpu.toFixed(2)}`;
};

// 内存格式化
export const formatMemory = (memory: string): string => {
  const bytes = parseInt(memory);
  return formatBytes(bytes);
};

// YAML 格式化
export const formatYaml = (obj: any): string => {
  // 简单实现，生产环境建议使用 js-yaml
  return JSON.stringify(obj, null, 2);
};

// 状态颜色
export const getStatusColor = (status: string): string => {
  const colorMap: Record<string, string> = {
    Running: 'success',
    Pending: 'warning',
    Failed: 'danger',
    Succeeded: 'info',
  };
  return colorMap[status] || 'info';
};

// 截断文本
export const truncate = (text: string, length = 50): string => {
  if (text.length <= length) return text;
  return text.substring(0, length) + '...';
};
```

## K8s API 封装示例

### Pod 管理
```typescript
// src/api/kube/pods.ts
import request from '@/utils/request';
import type { Pod, PodList, PodQuery, PodCreate } from '@/types/kube/pod';

enum Api {
  List = '/v1/pods',
  Get = '/v1/namespaces/:namespace/pods/:name',
  Create = '/v1/namespaces/:namespace/pods',
  Update = '/v1/namespaces/:namespace/pods/:name',
  Delete = '/v1/namespaces/:namespace/pods/:name',
  Logs = '/v1/namespaces/:namespace/pods/:name/logs',
  Exec = '/v1/namespaces/:namespace/pods/:name/exec',
}

export const podApi = {
  // 列表
  list: (params: PodQuery) => {
    return request.get<PodList>(Api.List, { params });
  },

  // 详情
  get: (namespace: string, name: string) => {
    return request.get<Pod>(Api.Get, {
      params: { namespace, name },
    });
  },

  // 创建
  create: (namespace: string, data: PodCreate) => {
    return request.post<Pod>(Api.Create, data, {
      params: { namespace },
    });
  },

  // 更新
  update: (namespace: string, name: string, data: Partial<Pod>) => {
    return request.put<Pod>(Api.Update, data, {
      params: { namespace, name },
    });
  },

  // 删除
  delete: (namespace: string, name: string) => {
    return request.delete(Api.Delete, {
      params: { namespace, name },
    });
  },

  // 日志（返回 EventSource 用于实时流）
  logs: (namespace: string, name: string) => {
    const baseURL = import.meta.env.VITE_API_BASE_URL || '';
    const url = `${baseURL}${Api.Logs}`.replace(':namespace', namespace).replace(':name', name);
    return new EventSource(url);
  },

  // 执行命令（返回 WebSocket）
  exec: (namespace: string, name: string, command: string[]) => {
    const baseURL = import.meta.env.VITE_API_BASE_URL || '';
    const url = `${baseURL}${Api.Exec}`.replace(':namespace', namespace).replace(':name', name);
    return new WebSocket(`${url}?command=${command.join(',')}`);
  },
};
```

## TypeScript 类型定义模板

### K8s 资源类型
```typescript
// src/types/kube/common.ts
export interface Metadata {
  name: string;
  namespace?: string;
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
  creationTimestamp?: string;
}

export interface ObjectReference {
  kind: string;
  name: string;
  namespace?: string;
}

export interface ContainerState {
  waiting?: { reason?: string; message?: string };
  running?: { startedAt: string };
  terminated?: { exitCode: number; reason?: string; message?: string };
}

export interface ContainerStatus {
  name: string;
  ready: boolean;
  restartCount: number;
  image: string;
  state: ContainerState;
}

export interface ResourceRequirements {
  limits?: {
    cpu?: string;
    memory?: string;
  };
  requests?: {
    cpu?: string;
    memory?: string;
  };
}

export interface Port {
  containerPort: number;
  protocol?: 'TCP' | 'UDP';
  name?: string;
}

// src/types/kube/pod.ts
import type { Metadata, ContainerStatus, ResourceRequirements, Port } from './common';

export interface Pod {
  metadata: Metadata;
  spec: PodSpec;
  status: PodStatus;
}

export interface PodSpec {
  containers: Container[];
  restartPolicy?: 'Always' | 'OnFailure' | 'Never';
  nodeName?: string;
}

export interface Container {
  name: string;
  image: string;
  ports?: Port[];
  resources?: ResourceRequirements;
  env?: Array<{ name: string; value?: string; valueFrom?: any }>;
  volumeMounts?: any[];
}

export interface PodStatus {
  phase: 'Pending' | 'Running' | 'Succeeded' | 'Failed' | 'Unknown';
  podIP?: string;
  hostIP?: string;
  startTime?: string;
  containerStatuses?: ContainerStatus[];
  conditions?: Condition[];
}

export interface Condition {
  type: string;
  status: 'True' | 'False' | 'Unknown';
  reason?: string;
  message?: string;
}

export interface PodList {
  items: Pod[];
  continue?: string;
  remainingItemCount?: number;
}

export interface PodQuery {
  namespace?: string;
  labelSelector?: string;
  fieldSelector?: string;
  limit?: number;
  continue?: string;
}

export interface PodCreate {
  metadata: Metadata;
  spec: PodSpec;
}
```

## 页面模板 - Pod 列表

```vue
<!-- src/views/kube/pods/index.vue -->
<template>
  <div class="pods-page">
    <!-- 页面头部 -->
    <el-page-header @back="goBack" title="Pod 管理">
      <template #content>
        <el-descriptions :column="3" size="small">
          <el-descriptions-item label="命名空间">
            {{ namespace || '全部' }}
          </el-descriptions-item>
          <el-descriptions-item label="Pod 数量">
            {{ pagination.total }}
          </el-descriptions-item>
        </el-descriptions>
      </template>
    </el-page-header>

    <!-- 工具栏 -->
    <div class="toolbar">
      <el-select v-model="namespace" placeholder="选择命名空间" clearable>
        <el-option label="全部" value="" />
        <el-option
          v-for="ns in namespaces"
          :key="ns"
          :label="ns"
          :value="ns"
        />
      </el-select>
      <el-button type="primary" @click="handleCreate">
        创建 Pod
      </el-button>
      <el-button @click="loadPods">刷新</el-button>
    </div>

    <!-- Pod 列表 -->
    <el-table
      v-loading="loading"
      :data="pods"
      stripe
      @row-click="showDetail"
    >
      <!-- 展开：容器信息 -->
      <el-table-column type="expand">
        <template #default="{ row }">
          <el-table :data="row.status.containerStatuses">
            <el-table-column prop="name" label="容器名称" />
            <el-table-column label="状态">
              <template #default="{ row: container }">
                <el-tag :type="container.ready ? 'success' : 'danger'">
                  {{ container.ready ? '就绪' : '未就绪' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="restartCount" label="重启次数" />
            <el-table-column prop="image" label="镜像" />
          </el-table>
        </template>
      </el-table-column>

      <el-table-column prop="metadata.name" label="名称" />
      <el-table-column prop="metadata.namespace" label="命名空间" />
      <el-table-column prop="status.phase" label="状态">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status.phase)">
            {{ row.status.phase }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间">
        <template #default="{ row }">
          {{ formatTime(row.metadata.creationTimestamp) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250">
        <template #default="{ row }">
          <el-button size="small" link @click.stop="showLogs(row)">
            日志
          </el-button>
          <el-button size="small" link @click.stop="execPod(row)">
            终端
          </el-button>
          <el-button size="small" link @click.stop="showYaml(row)">
            YAML
          </el-button>
          <el-button
            size="small"
            link
            type="danger"
            @click.stop="deletePod(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      @current-change="loadPods"
    />

    <!-- 详情抽屉 -->
    <el-drawer v-model="detailVisible" title="Pod 详情" size="50%">
      <PodDetail v-if="currentPod" :pod="currentPod" />
    </el-drawer>

    <!-- 日志抽屉 -->
    <el-drawer v-model="logsVisible" title="Pod 日志" size="60%">
      <PodLogs
        v-if="currentPod"
        :namespace="currentPod.metadata.namespace"
        :name="currentPod.metadata.name"
      />
    </el-drawer>

    <!-- YAML 对话框 -->
    <el-dialog v-model="yamlVisible" title="YAML 配置" width="60%">
      <YamlEditor v-model="yamlContent" readonly />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { podApi } from '@/api/kube/pods';
import type { Pod } from '@/types/kube/pod';
import { formatTime } from '@/utils/format';
import PodDetail from './components/PodDetail.vue';
import PodLogs from './components/PodLogs.vue';
import YamlEditor from '@/components/YamlEditor.vue';

const router = useRouter();

// 状态
const namespace = ref('');
const pods = ref<Pod[]>([]);
const loading = ref(false);
const namespaces = ref<string[]>([]);

const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
});

const detailVisible = ref(false);
const logsVisible = ref(false);
const yamlVisible = ref(false);
const currentPod = ref<Pod | null>(null);
const yamlContent = ref('');

// 方法
const loadPods = async () => {
  loading.value = true;
  try {
    const res = await podApi.list({
      namespace: namespace.value,
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
    });
    pods.value = res.items;
    pagination.value.total = res.total;
  } catch (error) {
    ElMessage.error('加载 Pod 列表失败');
  } finally {
    loading.value = false;
  }
};

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    Running: 'success',
    Pending: 'warning',
    Failed: 'danger',
    Succeeded: 'info',
  };
  return map[status] || 'info';
};

const showDetail = (pod: Pod) => {
  currentPod.value = pod;
  detailVisible.value = true;
};

const showLogs = (pod: Pod) => {
  currentPod.value = pod;
  logsVisible.value = true;
};

const showYaml = (pod: Pod) => {
  yamlContent.value = JSON.stringify(pod, null, 2);
  yamlVisible.value = true;
};

const execPod = (pod: Pod) => {
  // 打开终端
  router.push({
    name: 'PodExec',
    params: {
      namespace: pod.metadata.namespace,
      name: pod.metadata.name,
    },
  });
};

const deletePod = async (pod: Pod) => {
  try {
    await ElMessageBox.confirm(
      `确定删除 Pod "${pod.metadata.name}" 吗？`,
      '确认删除'
    );
    await podApi.delete(pod.metadata.namespace!, pod.metadata.name);
    ElMessage.success('删除成功');
    loadPods();
  } catch (error) {
    // 用户取消
  }
};

const goBack = () => {
  router.back();
};

const handleCreate = () => {
  router.push({ name: 'PodCreate' });
};

// 生命周期
onMounted(() => {
  loadPods();
});
</script>

<style scoped>
.pods-page {
  padding: 20px;
}

.toolbar {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}
</style>
```

这些模板可以直接使用，帮助您快速上手开发！
