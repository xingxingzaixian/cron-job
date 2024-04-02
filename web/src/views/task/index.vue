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
      <ServiceOperateDrawer
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
import ServiceOperateDrawer, { type OperateType } from './modules/service-operate-drawer.vue';
import { fetchServiceList } from '@/api/service';
import type { ServiceListOutput, ServiceSearchParam } from '@/api/service/types';
import { NButton, NPopconfirm } from 'naive-ui';
defineOptions({ name: 'ServiceList' });

const checkedRowKeys = ref<string[]>([]);
const dataId = ref<number>(0);
const { bool: drawerVisible, setTrue: openDrawer } = useBoolean();

const { columns, filteredColumns, data, loading, pagination, updateSearchParams, getData } = useTable<
  ServiceListOutput,
  ServiceSearchParam
>({
  apiFn: fetchServiceList,
  apiParams: {
    pageNo: 1,
    pageSize: 10,
    info: ''
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
      title: 'ID',
      align: 'center',
      width: 64
    },
    {
      key: 'serviceName',
      title: $t('page.service.list.serviceName'),
      minWidth: 100
    },
    {
      key: 'serviceDesc',
      title: $t('page.service.list.serviceDesc'),
      minWidth: 100
    },
    {
      key: 'rule',
      title: $t('page.service.list.rule'),
      minWidth: 100
    },
    {
      key: 'urlRewrite',
      title: $t('page.service.list.urlRewrite'),
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
