<template>
  <div class="flex flex-col items-stretch gap-16px overflow-hidden <sm:overflow-auto">
    <NCard
      :title="$t('page.service.list.title')"
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
import TableHeaderOperation from '@/components/custom/TableHeaderOperation.vue';
import TaskOperateDrawer, { type OperateType } from './modules/task-operate-drawer.vue';
import { fetchTaskList } from '@/api/task';
import type { SearchTaskResponse, QueryTask } from '@/api/task/types';
import { NButton, NPopconfirm } from 'naive-ui';
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
      data: list,
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
      width: 64
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
      key: 'cron表达式',
      title: $t('page.task.list.spec'),
      minWidth: 100
    },
    {
      key: 'protocol',
      title: $t('page.task.list.protocol'),
      minWidth: 100
    },
    {
      key: 'status',
      title: $t('page.task.list.status'),
      minWidth: 100
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 130,
      render: (row: any) => (
        <div class="flex-center gap-8px">
          <NButton type="primary" ghost size="small" class="mr-4px" onClick={() => handleEdit(row.id)}>
            {$t('common.edit')}
          </NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{
              default: () => $t('common.confirmDelete'),
              trigger: () => (
                <NButton type="error" ghost size="small">
                  {$t('common.delete')}
                </NButton>
              )
            }}
          </NPopconfirm>
        </div>
      )
    }
  ]
});

async function handleBatchDelete() {
  // request
  console.log(checkedRowKeys.value);
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
  // request
  console.log(id);
  message.success($t('common.deleteSuccess'));

  getData();
}

const operateType = ref<OperateType>('add');
function handleAdd() {
  operateType.value = 'add';
  openDrawer();
}
</script>

<style scoped></style>
