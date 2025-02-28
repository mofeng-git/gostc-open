<script setup>

import AppCard from "../../layout/components/AppCard.vue";
</script>

<template>
  <AppCard :show-border="false">
    <n-result status="404" title="404 资源不存在"></n-result>
  </AppCard>
</template>

<style scoped lang="scss">

</style>