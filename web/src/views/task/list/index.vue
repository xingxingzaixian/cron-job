<template>
  <div class="flex flex-col items-stretch gap-16px overflow-hidden <sm:overflow-auto">
    <TableSearch />
    <NCard
      :title="$t('page.task.list.title')"
      :bordered="false"
      size="small"
      class="sm:flex-1 sm:overflow-hidden rd-8px shadow-sm"
    >
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
        class="sm:h-full"
      />
      <TaskOperateDrawer
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :dataId="dataId"
        @submitted="getData"
      />
    </NCard>
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

const { columns, filteredColumns, data, loading, pagination, updateSearchParams, getData } = useTable<
  SearchTaskResponse,
  QueryTask
>({
  apiFn: fetchTaskList,
  apiParams: {
    pageNo: 1,
    pageSize: 15
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

const operateType = ref<OperateType>('add');
function handleAdd() {
  operateType.value = 'add';
  openDrawer();
}
</script>

<style scoped></style>
