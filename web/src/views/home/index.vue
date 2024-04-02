<template>
  <NSpace vertical :size="16">
    <NGrid :x-gap="16" :y-gap="16" responsive="screen" item-responsive>
      <NGi span="6">
        <Header
          :title="$t('page.home.dayTitle')"
          icon="ant-design:bar-chart-outlined"
          :count="dataList.day.reduce((accumulator, currentValue) => accumulator + currentValue.value, 0)"
        />
      </NGi>
      <NGi span="6">
        <Header
          :title="$t('page.home.weekTitle')"
          icon="mdi:chart-line"
          :count="dataList.week.reduce((accumulator, currentValue) => accumulator + currentValue.value, 0)"
        />
      </NGi>
      <NGi span="6">
        <Header
          :title="$t('page.home.interfaceTitle')"
          icon="mdi:chart-bell-curve"
          :count="dataList.interface.reduce((accumulator, currentValue) => accumulator + currentValue.value, 0)"
        />
      </NGi>
    </NGrid>

    <NGrid responsive="screen" item-responsive>
      <NGi span="12">
        <LineChart
          :title="$t('page.home.dayTitle')"
          :x-axis="$t('page.home.hour')"
          :data="dataList.day"
          color="#ec4786"
        />
      </NGi>
      <NGi span="12">
        <LineChart
          :title="$t('page.home.weekTitle')"
          :x-axis="$t('page.home.day')"
          :data="dataList.week"
          color="#865ec0"
        />
      </NGi>
      <NGi span="8">
        <LineChart
          :title="$t('page.home.interfaceTitle')"
          :x-axis="$t('page.home.interface')"
          :data="dataList.interface"
          color="#fcbc25"
        />
      </NGi>
    </NGrid>
  </NSpace>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
import { NGrid } from 'naive-ui';
import Header from './modules/header.vue';
import LineChart from './modules/line-chart.vue';
import { fetchServiceStatistics } from '@/api/service';
import type { ServiceStatisticsRes } from '@/api/service/types';

const dataList = reactive<{
  day: ServiceStatisticsRes[];
  week: ServiceStatisticsRes[];
  interface: ServiceStatisticsRes[];
}>({
  day: [],
  week: [],
  interface: []
});

fetchServiceStatistics('day').then((res) => {
  if (res.data) {
    dataList.day.push(...res.data);
  }
});

fetchServiceStatistics('week').then((res) => {
  if (res.data) {
    dataList.week.push(...res.data);
  }
});

fetchServiceStatistics('interface').then((res) => {
  if (res.data) {
    dataList.interface.push(...res.data);
  }
});
</script>
