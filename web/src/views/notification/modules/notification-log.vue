<template>
  <NModal v-model:show="showModal" preset="card" title="通知发送日志" style="width: 800px">
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
    />
  </NModal>
</template>

<script lang="tsx" setup>
import { ref, watch } from 'vue';
import { useTable } from '@/hooks';
import { fetchNotificationLogs } from '@/api/notification';
import type { NotificationLogResponse, NotificationLogInput } from '@/api/notification/types';
import { NTag } from 'naive-ui';

const props = defineProps<{
  show: boolean;
  taskId: number;
}>();

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
}>();

const showModal = ref(props.show);

const { columns, data, loading, pagination, updateSearchParams, resetSearchParams, getData } = useTable<
  NotificationLogResponse,
  NotificationLogInput
>({
  apiFn: fetchNotificationLogs,
  apiParams: {
    pageNo: 1,
    pageSize: 15,
    task_id: 0,
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
      title: 'ID',
      align: 'center',
      width: 70,
    },
    {
      key: 'notification_id',
      title: '通知配置ID',
      align: 'center',
      width: 100,
    },
    {
      key: 'task_id',
      title: '任务ID',
      align: 'center',
      width: 80,
    },
    {
      key: 'task_log_id',
      title: '任务日志ID',
      align: 'center',
      width: 100,
    },
    {
      key: 'status',
      title: '发送状态',
      minWidth: 100,
      render: (row: any) => (
        <div>
          {row.status === 1 ? (
            <NTag type="success" size="small">成功</NTag>
          ) : (
            <NTag type="error" size="small">失败</NTag>
          )}
        </div>
      )
    },
    {
      key: 'result',
      title: '发送结果',
      minWidth: 200,
      ellipsis: { tooltip: true },
    },
    {
      key: 'retry_count',
      title: '重试次数',
      align: 'center',
      width: 80,
    },
    {
      key: 'start_time',
      title: '开始时间',
      minWidth: 160,
    },
    {
      key: 'end_time',
      title: '结束时间',
      minWidth: 160,
    },
  ]
});

function handleRefresh() {
  resetSearchParams();
  updateSearchParams({
    task_id: props.taskId,
    pageNo: 1,
    pageSize: 15,
  });
  getData();
}

watch(
  () => props.show,
  (val) => {
    showModal.value = val;
    if (val) {
      handleRefresh();
    }
  }
);

watch(
  () => showModal.value,
  (val) => {
    emit('update:show', val);
  }
);
</script>

<style scoped>
.data-table {
  border-radius: 12px;
}
</style>
