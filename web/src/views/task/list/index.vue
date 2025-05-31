<template>
  <div class="task-management">
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
            {{ $t('page.task.list.title') }}
          </h1>
          <p class="page-description">管理定时任务、监控执行状态和查看运行日志</p>
        </div>
        <div class="header-actions">
          <n-button type="primary" size="large" @click="handleAdd" class="add-btn">
            <template #icon>
              <n-icon>
                <svg viewBox="0 0 24 24">
                  <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                </svg>
              </n-icon>
            </template>
            新增任务
          </n-button>
        </div>
      </div>
    </div>

    <!-- 搜索区域 -->
    <n-card :bordered="false" class="search-card">
      <TableSearch @search="search" @reset="reset" />
    </n-card>

    <!-- 数据表格 -->
    <n-card :bordered="false" class="table-card">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="filteredColumns"
          :canDelete="false"
          :loading="loading"
          @add="handleAdd"
          @delete="handleBatchDelete"
          @refresh="getData"
        />
      </template>
      <NDataTable
        v-model:checked-row-keys="checkedRowKeys"
        :columns="columns"
        :data="data"
        size="small"
        :scroll-x="962"
        :loading="loading"
        remote
        :pagination="pagination"
        :row-key="(item: any) => item.id"
        class="task-table"
      />
      <TaskOperateDrawer
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :dataId="dataId"
        @submitted="getData"
      />
    </n-card>
  </div>
</template>

<script lang="tsx" setup>
import { ref } from 'vue';
import { $t } from '@/locales';
import { message } from '@/utils/message';
import { useBoolean, useTable } from '@/hooks';
import TableSearch from './modules/task-search.vue';
import TableHeaderOperation from '@/components/custom/TableHeaderOperation.vue';
import TaskOperateDrawer, { type OperateType } from './modules/task-operate-drawer.vue';
import { fetchTaskList, fetchTaskOp } from '@/api/task';
import type { SearchTaskResponse, QueryTask, TaskItemOutput } from '@/api/task/types';
import { NButton, NPopconfirm, NSwitch } from 'naive-ui';
import { TaskProtocol, TaskStatus } from '@/enum/task';
import router from '@/router';

defineOptions({ name: 'TaskList' });

const checkedRowKeys = ref<string[]>([]);
const dataId = ref<number>(0);
const { bool: drawerVisible, setTrue: openDrawer } = useBoolean();

const { columns, filteredColumns, data, loading, pagination, updateSearchParams, resetSearchParams, getData } =
  useTable<SearchTaskResponse, QueryTask>({
    apiFn: fetchTaskList,
    apiParams: {
      pageNo: 1,
      pageSize: 15,
      name: '',
      tag: '',
      protocol: 0
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
        key: 'id',
        title: '任务ID',
        align: 'center',
        width: 100
      },
      {
        key: 'name',
        title: $t('page.task.list.name'),
        minWidth: 100
      },
      {
        key: 'tag',
        title: $t('page.task.list.tag'),
        minWidth: 100
      },
      {
        key: 'spec',
        title: $t('page.task.list.spec'),
        minWidth: 100
      },
      {
        key: 'protocol',
        title: $t('page.task.list.protocol'),
        minWidth: 40,
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
        key: 'status',
        title: $t('page.task.list.status'),
        minWidth: 40,
        render: (row: any) => (
          <div>
            <NSwitch
              value={row.status !== TaskStatus.Disabled}
              on-update:value={(val: boolean) => handleEnable(val, row)}
            />
          </div>
        )
      },
      {
        key: 'operate',
        title: $t('common.operate'),
        align: 'center',
        width: 200,
        render: (row: any) => (
          <div class="flex items-center justify-center gap-8px">
            <NButton type="primary" ghost size="small" class="mr-4px" onClick={() => handleEdit(row.id)}>
              {$t('common.edit')}
            </NButton>
            <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
              {{
                default: () => $t('common.confirmDelete'),
                trigger: () => (
                  <NButton type="error" ghost size="small" class="mr-4px">
                    {$t('common.delete')}
                  </NButton>
                )
              }}
            </NPopconfirm>
            <NPopconfirm onPositiveClick={() => handleRun(row.id)}>
              {{
                default: () => $t('page.task.list.runDesc'),
                trigger: () => (
                  <NButton type="info" ghost size="small">
                    {$t('page.task.list.runTask')}
                  </NButton>
                )
              }}
            </NPopconfirm>
            <NButton ghost size="small" class="mr-4px" onClick={() => handleLog(row.id)}>
              {$t('common.log')}
            </NButton>
          </div>
        )
      }
    ]
  });

async function handleBatchDelete() {
  message.success($t('common.deleteSuccess'));

  checkedRowKeys.value = [];

  getData();
}

function handleEdit(id: number) {
  operateType.value = 'edit';
  dataId.value = id;
  openDrawer();
}

async function handleDelete(id: number) {
  const res = await fetchTaskOp({ id, op: 'delete' });
  if (res.code === 200) {
    message.success($t('task.message.deleteSuccess'));
    getData();
  } else {
    message.error($t('task.message.deleteFailed'));
  }
}

async function handleRun(id: number) {
  await fetchTaskOp({ id, op: 'run' });
}

async function handleEnable(val: boolean, row: TaskItemOutput) {
  const res = await fetchTaskOp({ id: row.id, op: val ? 'start' : 'stop' });
  if (res.code === 200) {
    if (val) {
      message.success($t('task.message.startSuccess'));
    } else {
      message.success($t('task.message.stopSuccess'));
    }
    row.status = val ? TaskStatus.Enabled : TaskStatus.Disabled;
  } else {
    if (val) {
      message.success($t('task.message.startFailed'));
    } else {
      message.success($t('task.message.stopFailed'));
    }
  }
}

async function handleLog(id: number) {
  router.push({ name: 'TaskLog', query: { id } });
}

const search = (model: Omit<QueryTask, 'pageNo' | 'pageSize'>) => {
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

const operateType = ref<OperateType>('add');
function handleAdd() {
  operateType.value = 'add';
  openDrawer();
}
</script>

<style scoped>
/* 登录页面风格的渐变背景色彩 */
.task-management {
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

.add-btn {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  border-radius: 12px;
  padding: 0 24px;
  height: 44px;
  font-weight: 500;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
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

.task-table {
  border-radius: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .task-management {
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
