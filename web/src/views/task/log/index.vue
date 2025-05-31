<template>
  <div class="task-log-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="page-title">
            <n-icon size="24" class="title-icon">
              <svg viewBox="0 0 24 24">
                <path fill="currentColor" d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z"/>
              </svg>
            </n-icon>
            任务执行日志
          </h1>
          <p class="page-description">查看任务执行历史、运行状态和详细日志信息</p>
        </div>
        <div class="header-actions">
          <NPopconfirm @positive-click="batchDelete">
            <template #trigger>
              <n-button type="error" size="large" class="delete-btn">
                <template #icon>
                  <n-icon>
                    <svg viewBox="0 0 24 24">
                      <path fill="currentColor" d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
                    </svg>
                  </n-icon>
                </template>
                {{ $t('common.batchDelete') }}
              </n-button>
            </template>
            {{ $t('common.confirmDelete') }}
          </NPopconfirm>
        </div>
      </div>
    </div>

    <!-- 搜索区域 -->
    <n-card :bordered="false" class="search-card">
      <TaskLogSearch @search="search" @reset="reset" />
    </n-card>

    <!-- 数据表格 -->
    <n-card :bordered="false" class="table-card">
      <NDataTable
        :columns="columns"
        :data="data"
        size="small"
        :scroll-x="962"
        :loading="loading"
        remote
        :pagination="pagination"
        :row-key="(item: any) => item.id"
        @update:checked-row-keys="handleCheck"
        class="log-table"
      />
    </n-card>

    <LogInfo v-model:show="showModal" :task-item="taskItem" />
  </div>
</template>

<script lang="tsx" setup>
import { $t } from '@/locales';
import { ref } from 'vue';
import LogInfo from './modules/info.vue';
import { useTable } from '@/hooks';
import TaskLogSearch from './modules/task-log-search.vue';
import { fetchTaskLogList, fetchTaskLogDelete } from '@/api/task';
import type { TaskLogOutput, TaskLogItemOutput, QueryTaskLog } from '@/api/task/types';
import { useRoute } from 'vue-router';
import { TaskProtocol, TaskStatus } from '@/enum/task';
import { NTag, NButton } from 'naive-ui';
import type { DataTableRowKey } from 'naive-ui';
import { message } from '@/utils/message';

defineOptions({ name: 'ServiceLog' });

const route = useRoute();
const showModal = ref<boolean>(false);
const taskId = Number(route.query.id) || 0;
let checkedRowKeys: DataTableRowKey[] = [];
const taskItem = ref<TaskLogItemOutput | null>(null);
const { columns, data, loading, pagination, updateSearchParams, resetSearchParams, getData } = useTable<
  TaskLogOutput,
  QueryTaskLog
>({
  apiFn: fetchTaskLogList,
  apiParams: {
    pageNo: 1,
    pageSize: 15,
    taskId,
    taskName: '',
    status: -1
  },
  transformer: (res: any) => {
    const { list = [], total = 0 } = res.data || {};

    return {
      data: list || [],
      pageNum: pagination.page,
      pageSize: pagination.pageSize,
      total
    };
  },
  onPaginationChanged(pg) {
    const { page, pageSize } = pg;

    updateSearchParams({
      pageNo: page,
      pageSize: pageSize
    });

    getData();
  },
  columns: () => [
    {
      type: 'selection'
    },
    {
      key: 'id',
      title: 'ID',
      align: 'center',
      width: 100
    },
    {
      key: 'taskId',
      title: $t('page.task.log.taskId'),
      minWidth: 100
    },
    {
      key: 'taskName',
      title: $t('page.task.log.taskName'),
      minWidth: 100
    },
    {
      key: 'protocol',
      title: $t('page.task.log.protocol'),
      minWidth: 100,
      render: (row: any) => (
        <div>
          {row.protocol === TaskProtocol.HTTP ? (
            <span>HTTP</span>
          ) : row.protocol === TaskProtocol.Shell ? (
            <span>Shell</span>
          ) : (
            <span>Grpc</span>
          )}
        </div>
      )
    },
    {
      key: 'retryTimes',
      title: $t('page.task.log.retryTimes'),
      minWidth: 100
    },
    {
      key: 'totalTime',
      title: $t('page.task.log.runTime'),
      minWidth: 100
    },
    {
      key: 'startTime',
      title: $t('page.task.log.startTime'),
      minWidth: 100
    },
    {
      key: 'status',
      title: $t('page.task.log.runStatus'),
      minWidth: 100,
      render: (row: any) => (
        <div>
          {row.status === TaskStatus.Failure ? (
            <NTag type="error" size="small">
              {$t('page.task.log.failure')}
            </NTag>
          ) : row.status === TaskStatus.Running ? (
            <NTag type="warning" size="small">
              {$t('page.task.log.running')}
            </NTag>
          ) : row.status === TaskStatus.Finish ? (
            <NTag type="success" size="small">
              {$t('page.task.log.success')}
            </NTag>
          ) : (
            <NTag type="info" size="small">
              {$t('page.task.log.timeout')}
            </NTag>
          )}
        </div>
      )
    },
    {
      key: 'operate',
      title: $t('page.task.log.runResult'),
      minWidth: 100,
      render: (row: any) => (
        <div>
          <NButton type="primary" ghost size="small" onClick={() => handleView(row)}>
            {$t('page.task.log.runResult')}
          </NButton>
        </div>
      )
    }
  ]
});

const handleView = async (item: TaskLogItemOutput) => {
  if (item) {
    taskItem.value = item;
    showModal.value = true;
  }
};

const search = (model: Omit<QueryTaskLog, 'pageNo' | 'pageSize'>) => {
  updateSearchParams({
    pageNo: 1,
    pageSize: 15,
    ...model
  });

  getData();
};

const reset = () => {
  resetSearchParams();
  getData();
};

const handleCheck = (rowKeys: DataTableRowKey[]) => {
  checkedRowKeys = rowKeys;
};

const batchDelete = async () => {
  await fetchTaskLogDelete(checkedRowKeys as number[]);
  message.success($t('common.deleteSuccess'));

  reset();
};
</script>

<style scoped>
/* 登录页面风格的渐变背景色彩 */
.task-log-management {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
}

/* 页面头部 */
.page-header {
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 24px 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.header-left {
  flex: 1;
}

.page-title {
  display: flex;
  align-items: center;
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.title-icon {
  margin-right: 12px;
  color: #667eea;
}

.page-description {
  margin: 0;
  color: #64748b;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.delete-btn {
  border-radius: 12px;
  padding: 0 24px;
  height: 44px;
  font-weight: 500;
  box-shadow: 0 4px 16px rgba(239, 68, 68, 0.3);
  transition: all 0.3s ease;
}

.delete-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(239, 68, 68, 0.4);
}

/* 搜索卡片 */
.search-card {
  margin-bottom: 24px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

/* 表格卡片 */
.table-card {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.log-table {
  border-radius: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .task-log-management {
    padding: 16px;
  }

  .header-content {
    flex-direction: column;
    gap: 16px;
    padding: 20px;
  }

  .page-title {
    font-size: 24px;
  }
}

/* 动画效果 */
.search-card,
.table-card {
  animation: fadeInUp 0.6s ease-out;
}

.page-header {
  animation: fadeInDown 0.6s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 表格行悬停效果 */
:deep(.n-data-table-tbody .n-data-table-tr:hover) {
  background-color: rgba(102, 126, 234, 0.05);
}

/* 按钮样式优化 */
:deep(.n-button--primary-type) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
}

:deep(.n-button--primary-type:hover) {
  background: linear-gradient(135deg, #5a6fd8, #6a4190);
}
</style>
