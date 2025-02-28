<script setup>
</script>

<template>
  <n-scrollbar style="max-height: calc(100vh)">
    <n-card class="empty-container">
      <slot/>
    </n-card>
  </n-scrollbar>
</template>

<style lang="scss" scoped>
//.empty-container {
//  width: 100%;
//  height: 100vh;
//  overflow: hidden;
//  padding: 0;
//  margin: 0;
//}

.n-card {
  //margin: 0 !important;
  width: 100%;
  border-radius: 0 !important;
  min-height: 100vh;

  :deep(.n-card__content) {
    padding: 0 !important;
  }

}

.n-card.n-card--bordered {
  border: 0 !important;
}
</style>