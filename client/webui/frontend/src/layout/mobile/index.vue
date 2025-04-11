<script setup>
import Header from "../components/Header.vue";
</script>

<template>
  <n-layout class="layout-container">
    <n-layout-header>
      <Header/>
    </n-layout-header>
    <n-layout-content class="layout-body">
      <n-scrollbar style="max-height: calc(100vh - 60px)">
        <slot/>
      </n-scrollbar>
    </n-layout-content>
  </n-layout>
</template>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;

  .layout-body {
    height: calc(100vh - 60px);
  }
}
</style>