<template>
  <div class="flex flex-col items-stretch gap-16px overflow-hidden <sm:overflow-auto">
    <NCard
      :title="$t('page.service.list.title')"
      :bordered="false"
      size="small"
      class="sm:flex-1 sm:overflow-hidden rd-8px shadow-sm"
    >
      <NDataTable
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
    </NCard>
  </div>
</template>

<script lang="ts" setup>
import { $t } from '@/locales';
import { useTable } from '@/hooks';
import { fetchServiceLogList } from '@/api/service';
import type { ServiceRequestLogList, ServiceSearchParam } from '@/api/service/types';

defineOptions({ name: 'ServiceLog' });

const { columns, data, loading, pagination, updateSearchParams, getData } = useTable<
  ServiceRequestLogList,
  ServiceSearchParam
>({
  apiFn: fetchServiceLogList,
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
      title: $t('page.service.log.serviceName'),
      minWidth: 100
    },
    {
      key: 'originUrl',
      title: $t('page.service.log.originUrl'),
      minWidth: 100
    },
    {
      key: 'proxyUrl',
      title: $t('page.service.log.proxyUrl'),
      minWidth: 100
    }
  ]
});
</script>

<style scoped></style>
