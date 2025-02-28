<script setup>

import AppCard from "../../layout/components/AppCard.vue";
</script>

<template>
  <AppCard :show-border="false">
    <n-result status="403" title="403 无权限"></n-result>
  </AppCard>
</template>

<style scoped lang="scss">

</style>