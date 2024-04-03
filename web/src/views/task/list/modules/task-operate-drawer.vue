<template>
  <NDrawer v-model:show="visible" :title="title" display-directive="show" :width="360">
    <NDrawerContent :title="title" :native-scrollbar="false" closable>
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.service.list.serviceName')" path="serviceName">
          <NInput v-model:value="model.serviceName" :placeholder="$t('page.service.list.form.serviceName')" />
        </NFormItem>
        <NFormItem :label="$t('page.service.list.serviceDesc')" path="serviceDesc">
          <NInput
            v-model:value="model.serviceDesc"
            :placeholder="$t('page.service.list.form.serviceDesc')"
            type="textarea"
          />
        </NFormItem>
        <NFormItem :label="$t('page.service.list.rule')" path="rule">
          <NInput v-model:value="model.rule" :placeholder="$t('page.service.list.form.rule')" />
        </NFormItem>
        <NFormItem :label="$t('page.service.list.urlRewrite')" path="urlRewrite">
          <NInput v-model:value="model.urlRewrite" :placeholder="$t('page.service.list.form.urlRewrite')" />
        </NFormItem>
        <NFormItem :label="$t('page.service.list.needHttps')" path="needHttps">
          <NSwitch v-model:value="model.needHttps" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace :size="16">
          <NButton @click="closeDrawer">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<script setup lang="ts">
import { $t } from '@/locales';
import type { FormRules } from 'naive-ui';
import { message } from '@/utils/message';
import { computed, reactive, watch } from 'vue';
import type { ServiceUpdateHTTPInput } from '@/api/service/types';
import { fetchServiceGet, fetchTaskEdit } from '@/api/task';
import { useForm } from '@/hooks';

defineOptions({
  name: 'ServiceOperateDrawer'
});

interface Props {
  /** the type of operation */
  operateType: OperateType;
  /** the edit row data */
  dataId?: number | null;
}
const props = defineProps<Props>();
interface Emits {
  (e: 'submitted'): void;
}
const emit = defineEmits<Emits>();
export type OperateType = 'add' | 'edit';

const visible = defineModel<boolean>('visible', {
  default: false
});

const { formRef, validate, restoreValidation } = useForm();

const rules: FormRules = {};
const title = computed(() => {
  const titles: Record<OperateType, string> = {
    add: $t('page.service.list.addService'),
    edit: $t('page.service.list.editService')
  };
  return titles[props.operateType];
});

async function handleSubmit() {
  await validate();
  // request
  if (props.operateType === 'add') {
    const res = await fetchServiceAdd(model);
    if (res.code !== 200) {
      res.message && message.error(res.message);
      return;
    }
  } else {
    const res = await fetchServiceUpdate(model);
    if (res.code !== 200) {
      res.message && message.error(res.message);
      return;
    }
  }
  message.success($t('common.updateSuccess'));
  closeDrawer();
  emit('submitted');
}

const model: ServiceUpdateHTTPInput = reactive(createDefaultModel());
function createDefaultModel(): ServiceUpdateHTTPInput {
  return {
    id: null,
    serviceName: '',
    serviceDesc: '',
    rule: '',
    urlRewrite: '',
    needHttps: 0
  };
}

async function handleUpdateModelWhenEdit() {
  if (props.operateType === 'add') {
    Object.assign(model, createDefaultModel());
    return;
  }

  if (props.operateType === 'edit' && props.dataId) {
    fetchServiceGet(props.dataId).then((res) => {
      if (res.code === 200) {
        model.id = res.data.info.id;
        model.serviceName = res.data.info.serviceName;
        model.serviceDesc = res.data.info.serviceDesc;
        model.needHttps = res.data.httpRule.needHttps;
        model.rule = res.data.httpRule.rule;
        model.urlRewrite = res.data.httpRule.urlRewrite;
      }
    });
  }
}

function closeDrawer() {
  visible.value = false;
}

watch(visible, () => {
  if (visible.value) {
    handleUpdateModelWhenEdit();
    restoreValidation();
  }
});
</script>

<style scoped></style>
