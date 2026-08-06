<template>
  <div class="w-full">
    <!-- URL 栏 -->
    <div class="flex items-center gap-8px mb-16px">
      <NSelect
        :value="model.method"
        :options="methodOptions"
        class="w-130px"
        @update:value="onMethodChange"
      />
      <NInput v-model:value="model.url" :placeholder="$t('http.urlPlaceholder')" class="flex-1" @update:value="onChange" />
      <NButton type="primary" :loading="testing" @click="emit('test')">
        {{ $t('page.task.list.testButton') }}
      </NButton>
    </div>

    <NTabs type="line" animated>
      <NTabPane name="query" :tab="$t('http.query')">
        <KvRows :rows="model.queryRows" @change="onChange" />
      </NTabPane>
      <NTabPane name="headers" :tab="$t('http.header')">
        <KvRows :rows="model.headerRows" @change="onChange" />
      </NTabPane>
      <NTabPane name="body" :tab="$t('http.data')">
        <div class="mb-8px w-260px">
          <NSelect v-model:value="model.contentType" :options="contentTypeOptions" @update:value="onContentTypeChange" />
          <NInput
            v-if="model.contentType === 'custom'"
            v-model:value="model.customContentType"
            :placeholder="$t('http.customContentType')"
            class="mt-8px"
            @update:value="onChange"
          />
        </div>
        <KvRows v-if="model.contentType === 'application/x-www-form-urlencoded'" :rows="model.formRows" @change="onChange" />
        <NInput
          v-else-if="model.contentType !== 'none'"
          v-model:value="model.bodyText"
          type="textarea"
          :rows="8"
          :placeholder="$t('http.bodyPlaceholder')"
          @update:value="onChange"
        />
        <div v-else class="text-gray-400">{{ $t('http.bodyNone') }}</div>
      </NTabPane>
    </NTabs>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue';
import { $t } from '@/locales';
import type { SelectOption } from 'naive-ui';
import KvRows from './kv-rows.vue';
import type { KvRow } from './types';

const props = defineProps<{
  command: string;
  params: string;
  testing?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update', command: string, params: string): void;
  (e: 'test'): void;
}>();

const methodOptions: SelectOption[] = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' }
];

const contentTypeOptions: SelectOption[] = [
  { label: $t('http.bodyNone'), value: 'none' },
  { label: 'application/json', value: 'application/json' },
  { label: 'application/x-www-form-urlencoded', value: 'application/x-www-form-urlencoded' },
  { label: 'text/plain', value: 'text/plain' },
  { label: 'application/xml', value: 'application/xml' },
  { label: $t('http.customContentType'), value: 'custom' }
];

const model = reactive({
  method: 'GET',
  url: '',
  contentType: 'none',
  customContentType: '',
  bodyText: '',
  queryRows: [] as KvRow[],
  headerRows: [] as KvRow[],
  formRows: [] as KvRow[]
});

const encodeForm = (rows: KvRow[]) =>
  rows
    .filter((row) => row.enabled && row.key !== '')
    .map((row) => `${encodeURIComponent(row.key)}=${encodeURIComponent(row.value)}`)
    .join('&');

const parseForm = (source: string): KvRow[] => {
  const rows: KvRow[] = [];
  new URLSearchParams(source).forEach((value, key) => {
    rows.push({ key, value, enabled: true, desc: '' });
  });
  return rows;
};

const normalizeRow = (row: Partial<KvRow>): KvRow => ({
  key: row.key ?? '',
  value: row.value ?? '',
  enabled: row.enabled ?? true,
  desc: row.desc ?? ''
});

const effectiveContentType = computed(() => {
  if (model.contentType === 'none') {
    return '';
  }
  return model.contentType === 'custom' ? model.customContentType : model.contentType;
});

const changeData = () => {
  const headers = model.headerRows.map(normalizeRow);
  if (effectiveContentType.value) {
    headers.unshift({
      key: 'Content-Type',
      value: effectiveContentType.value,
      enabled: true,
      desc: ''
    });
  }

  let data = '';
  if (model.contentType === 'application/x-www-form-urlencoded') {
    data = encodeForm(model.formRows);
  } else if (model.contentType !== 'none') {
    data = model.bodyText;
  }

  const command = JSON.stringify({ url: model.url.trim(), method: model.method });
  const params = JSON.stringify({
    headers,
    query: model.queryRows.map(normalizeRow),
    data
  });
  emit('update', command, params);
};

const onChange = () => {
  changeData();
};

const onMethodChange = (value: string) => {
  model.method = value;
  changeData();
};

const onContentTypeChange = (value: string) => {
  model.contentType = value;
  changeData();
};

const parse = () => {
  try {
    const cmd = JSON.parse(props.command || '{}');
    model.url = cmd.url ?? '';
    model.method = cmd.method ?? 'GET';
  } catch (error) {
    console.error('command 解析失败:', error);
  }

  try {
    const params = JSON.parse(props.params || '{}');
    model.queryRows = ((params.query ?? []) as Partial<KvRow>[]).map(normalizeRow);

    const headers = (params.headers ?? []) as Partial<KvRow>[];
    const contentTypeRow = headers.find((row) => (row.key ?? '').toLowerCase() === 'content-type');
    model.headerRows = headers.filter((row) => (row.key ?? '').toLowerCase() !== 'content-type').map(normalizeRow);

    if (contentTypeRow) {
      const value = contentTypeRow.value ?? '';
      if (contentTypeOptions.some((option) => option.value === value)) {
        model.contentType = value;
      } else {
        model.contentType = 'custom';
        model.customContentType = value;
      }
    } else {
      model.contentType = 'none';
    }

    const data = (params.data ?? '') as string;
    if (model.contentType === 'application/x-www-form-urlencoded') {
      model.formRows = parseForm(data);
    } else {
      model.bodyText = data;
    }
  } catch (error) {
    console.error('params 解析失败:', error);
  }
};

watch(
  () => [props.command, props.params],
  () => {
    parse();
  },
  { immediate: true }
);

// 初始挂载时向父组件同步一次默认配置
watch(
  () => props.command,
  () => {
    if (!props.command) {
      changeData();
    }
  },
  { immediate: true }
);
</script>
