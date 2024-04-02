import * as echarts from 'echarts';
import { onUnmounted, ref, computed, onMounted, watch, nextTick } from 'vue';
import useThemeStore from '@/store/modules/theme';
import { useResizeObserver } from '@vueuse/core';
import { templateRef } from '@vueuse/core';

const useEchart = <T extends echarts.EChartsOption>(optionsFactory: () => T) => {
  const domRef = templateRef('domRef', null);
  let chart: echarts.ECharts | null = null;
  const themeStore = useThemeStore();
  const darkMode = computed(() => (themeStore.darkMode ? 'light' : 'dark'));
  const chartOptions: T = optionsFactory();

  useResizeObserver(domRef, (entries) => {
    const entry = entries[0];
    const { width, height } = entry.contentRect;
    chart?.resize({ width, height });
  });

  onMounted(() => {
    render();
  });

  onUnmounted(() => {
    if (chart) {
      chart.clear();
      chart.dispose();
    }
  });

  const updateOptions = (callback: (opts: T, optsFactory: () => T) => T = () => chartOptions) => {
    if (!isRendered()) return;

    const updatedOpts = callback(chartOptions, optionsFactory);

    Object.assign(chartOptions, updatedOpts);

    if (isRendered()) {
      chart?.clear();
    }

    chart?.setOption({ ...chartOptions, backgroundColor: 'transparent' });
  };

  const isRendered = () => {
    return Boolean(domRef.value && chart);
  };

  const render = () => {
    if (!isRendered()) {
      nextTick(() => {
        chart = echarts.init(domRef.value!, darkMode.value);
        chart.clear();
        if (chart) {
          chart.setOption({ ...chartOptions, backgroundColor: 'transparent' });
        }
      });
    }
  };

  const destroy = () => {
    if (chart) {
      chart.dispose();
      chart = null;
    }
  };

  const changeTheme = () => {
    destroy();
    render();
  };

  watch(darkMode, () => {
    changeTheme();
  });

  return { domRef, updateOptions };
};

export default useEchart;
