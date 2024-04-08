<template>
  <div class="flex flex-col items-stretch gap-16px overflow-hidden <sm:overflow-auto">
    <TaskLogSearch @search="search" @reset="reset" />
    <NCard
      :title="$t('page.task.list.title')"
      :bordered="false"
      size="small"
      class="sm:flex-1 sm:overflow-hidden rd-8px shadow-sm"
    >
      <template #header-extra>
        <NPopconfirm @positive-click="batchDelete">
          <template #trigger>
            <NButton size="small" ghost type="error">
              <template #icon>
                <icon-ic-round-delete />
              </template>
              {{ $t('common.batchDelete') }}
            </NButton>
          </template>
          {{ $t('common.confirmDelete') }}
        </NPopconfirm>
      </template>
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
        class="sm:h-full"
      />
    </NCard>

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

<style scoped></style>
