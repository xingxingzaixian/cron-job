<template>
  <NDropdown :value="lang" :options="langOptions" trigger="hover" @select="changeLang">
    <div class="h-36px text-icon">
      <SvgIcon icon="heroicons:language" class="h-full" />
    </div>
  </NDropdown>
</template>

<script setup lang="ts">
defineOptions({
  name: 'LangSwitch'
});

interface Props {
  /** Current language */
  lang: I18n.LangType;
  /** Language options */
  langOptions: I18n.LangOption[];
}

defineProps<Props>();

type Emits = {
  (e: 'changeLang', lang: I18n.LangType): void;
};

const emit = defineEmits<Emits>();

function changeLang(lang: I18n.LangType) {
  emit('changeLang', lang);
}
</script>

<style scoped></style>
