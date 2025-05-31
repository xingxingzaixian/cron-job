<template>
  <div class="management-container">
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

    <!-- 数据表格 -->
    <n-card :bordered="false" class="glass-card table-card">
      <TaskLogSearch @search="search" @reset="reset" />
      <NDataTable
        :columns="columns"
        :data="data"
        size="small"
        :scroll-x="962"
        :loading="loading"
        remote
        :pagination="pagination"
        :row-key="(item: any) => item.id"
        class="data-table"
        @update:checked-row-keys="handleCheck"
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
/* 任务日志页面特定样式 */
/* 大部分样式已提取到 management-page.css 公共样式文件中 */

/* 头部布局优化 */
.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.header-left {
  flex-shrink: 0;
  min-width: 0;
}

.header-actions {
  flex-shrink: 0;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .header-content {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .header-content {
    gap: 12px;
  }

  .page-title {
    font-size: 20px;
  }

  .page-description {
    font-size: 13px;
  }
}
</style>
