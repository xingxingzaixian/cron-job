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
            {{ $t('page.task.list.title') }}
          </h1>
          <p class="page-description">管理定时任务、监控执行状态和查看运行日志</p>
        </div>
        <div class="header-actions">
          <n-button type="primary" size="large" class="add-btn" @click="handleAdd">
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

    <!-- 数据表格 -->
    <n-card :bordered="false" class="glass-card table-card">
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
      <TableSearch @search="search" @reset="reset" />
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
        class="data-table"
        :max-height="tableHeight"
      />
    </n-card>

    <!-- 抽屉组件 -->
    <TaskOperateDrawer
      v-model:visible="drawerVisible"
      :operate-type="operateType"
      :dataId="dataId"
      @submitted="getData"
    />
  </div>
</template>

<script lang="tsx" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
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

// 动态计算表格高度
const windowHeight = ref(window.innerHeight);

const tableHeight = computed(() => {
  // 减去页面头部(120px) + 卡片头部(60px) + 搜索区域(80px) + 分页器(60px) + 边距(40px)
  return windowHeight.value - 360;
});

// 监听窗口大小变化
const handleResize = () => {
  windowHeight.value = window.innerHeight;
};

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
/* 任务列表页面特定样式 */
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
