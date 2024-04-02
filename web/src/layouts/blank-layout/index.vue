<script setup lang="ts">
import LayoutContent from '../modules/content/index.vue';

defineOptions({
  name: 'BlankLayout'
});
</script>

<template>
  <h1>空白</h1>
  <LayoutContent :show-padding="true" />
</template>

<style scoped></style>
