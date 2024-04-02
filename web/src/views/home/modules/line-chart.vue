<template>
  <NCard :bordered="false" class="rd-8 shadow-sm">
    <div ref="domRef" class="h-360px overflow-hidden"></div>
  </NCard>
</template>

<script lang="ts" setup>
import { onMounted, watch } from 'vue';
import { useEchart } from '@/hooks';
import { $t } from '@/locales';
import type { ServiceStatisticsRes } from '@/api/service/types';

defineOptions({
  name: 'LineChart'
});

const props = defineProps<{
  title: string;
  color: string;
  xAxis: string;
  data: ServiceStatisticsRes[];
}>();

const { updateOptions } = useEchart(() => {
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      containLabel: true
    },
    xAxis: {
      name: props.xAxis,
      type: 'category',
      nameLocation: 'middle',
      nameGap: 25,
      boundaryGap: false,
      data: [] as string[]
    },
    yAxis: {
      name: $t('page.home.visits'),
      nameGap: 30,
      min: 30,
      boundaryGap: false,
      type: 'value'
    },
    series: [
      {
        name: props.title,
        color: props.color,
        type: 'line',
        stack: 'Total',
        smooth: true,
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0.25,
                color: props.color
              },
              {
                offset: 1,
                color: '#fff'
              }
            ]
          }
        },
        emphasis: {
          focus: 'series'
        },
        data: [] as number[]
      }
    ]
  };
});

watch(
  props.data,
  () => {
    if (!props.data) return;
    updateOptions((opts) => {
      opts.xAxis.data = props.data.map((item) => item.key.toString());
      opts.series[0].data = props.data.map((item) => item.value);
      return opts;
    });
  },
  { deep: true }
);
</script>
