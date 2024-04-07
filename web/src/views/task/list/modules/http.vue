<template>
  <NSpace vertical>
    <!-- define component: ParamsContent -->
    <DefineParamsContent v-slot="{ type, items }">
      <div class="flex items-center align-middle" v-for="key in Object.keys(items)" :key="key">
        <NInputGroup>
          <NInput :value="key" />
          <NInput :value="items[key]" />
          <NButton type="danger" ghost @click="handleDelete(type, key)">
            <template #icon>
              <icon-ic-round-delete color="red" />
            </template>
          </NButton>
        </NInputGroup>
      </div>
      <NInputGroup>
        <NInput v-model:value="addData.key" />
        <NInput v-model:value="addData.value" />
        <NButton type="danger" ghost @click="handleAdd(type)">
          <template #icon>
            <icon-ic-round-plus />
          </template>
        </NButton>
      </NInputGroup>
    </DefineParamsContent>
    <!-- define component: ParamsContent -->

    <NFormItem :label="$t('page.task.list.form.request')" path="request">
      <NInputGroup>
        <NSelect :style="{ width: '27%' }" :value="model.method" :options="options" :on-update:value="onSelectChange" />
        <NInputGroupLabel>http://</NInputGroupLabel>
        <NInput v-model:value="model.url" :on-change="onChange" />
      </NInputGroup>
    </NFormItem>
    <NFormItem :show-label="false">
      <NTabs type="line" animated default-value="header" v-model:value="activedTab">
        <NTabPane name="header" :tab="$t('http.header')">
          <ParamsContent :items="form.headers" type="headers" />
        </NTabPane>
        <NTabPane name="query" :tab="$t('http.query')"> <ParamsContent :items="form.query" type="query" /></NTabPane>
        <NTabPane v-if="model.method.toUpperCase() === 'POST'" name="data" :tab="$t('http.data')">
          <ParamsContent :items="form.data" type="data"
        /></NTabPane>
      </NTabs>
    </NFormItem>
  </NSpace>
</template>

<script lang="tsx" setup>
import { $t } from '@/locales';
import { reactive, watch, ref, nextTick } from 'vue';
import type { SelectOption } from 'naive-ui';
import type { HttpPro, HttpParams } from '@/types/task';
import { createReusableTemplate } from '@vueuse/core';

const props = defineProps<{
  command: string;
  params: string;
}>();

const emit = defineEmits<{
  (e: 'update', command: string, params: string): void;
}>();

interface ParamsContentProps {
  type: 'headers' | 'query' | 'data';
  items: Record<string, string>;
}

const [DefineParamsContent, ParamsContent] = createReusableTemplate<ParamsContentProps>();

const model = reactive<HttpPro>({
  url: '',
  method: 'GET'
});

const activedTab = ref<string>('header');
const addData = reactive<{
  key: string;
  value: string;
}>({
  key: '',
  value: ''
});

const form = reactive<HttpParams>({
  headers: {},
  query: {},
  data: {}
});

const options: SelectOption[] = [
  {
    label: 'GET',
    value: 'GET'
  },
  {
    label: 'POST',
    value: 'POST'
  }
];

const handleAdd = (type: string) => {
  if (addData.key.trim() === '') {
    return;
  }

  switch (type) {
    case 'headers':
      form.headers[addData.key.trim()] = addData.value.trim();
      break;
    case 'query':
      form.query[addData.key.trim()] = addData.value.trim();
      break;
    case 'data':
      form.data[addData.key.trim()] = addData.value.trim();
      break;
  }
  addData.key = '';
  addData.value = '';
  changeData();
};

const handleDelete = (type: string, key: string) => {
  switch (type) {
    case 'headers':
      delete form.headers[key.trim()];
      break;
    case 'query':
      delete form.query[key.trim()];
      break;
    case 'data':
      delete form.data[key.trim()];
      break;
  }
  changeData();
};

const onSelectChange = (value: string) => {
  model.method = value;
  changeData();
};

const onChange = () => {
  changeData();
};

const changeData = () => {
  console.log(322323, model.url);
  if (!model.url.trim().startsWith('http://')) {
    model.url = `http://${model.url}`;
  }

  const command = JSON.stringify(model);
  if (model.method.toUpperCase() === 'GET') {
    form.data = {};
  }
  const params = JSON.stringify(form);
  emit('update', command, params);
};

watch(
  () => props.command,
  () => {
    try {
      const res = JSON.parse(props.command);
      model.url = res.url.replace(/^\s*http?:\/\//, '');
      model.method = res.method;
    } catch (err) {
      console.log(err);
    }
  },
  {
    immediate: true
  }
);

watch(
  () => props.params,
  () => {
    try {
      Object.assign(form, JSON.parse(props.params));
    } catch (err) {
      console.log(err);
    }
  },
  {
    immediate: true
  }
);
</script>
