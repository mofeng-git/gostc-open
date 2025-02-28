<script setup>
import Header from "../components/Header.vue";
import SideMenu from "../components/SideMenu.vue";
import Menu from "../components/Menu.vue"
import {computed} from "vue";
import {localStore} from "../../store/local.js";
import {appStore} from "../../store/app.js";


// 是否为侧滑菜单
const layoutDirectionIsSideComputed = computed(() => {
  return appStore().width >= 500
})

</script>

<template>
  <n-layout class="layout-container">
    <n-layout-header>
      <Header/>
    </n-layout-header>
    <SideMenu v-if="!layoutDirectionIsSideComputed"/>

    <n-layout-content has-sider class="layout-body">
      <n-layout-sider
          bordered
          v-if="layoutDirectionIsSideComputed"
          collapse-mode="width"
          :collapsed-width="60"
          :width="240"
          :collapsed="localStore().isCollapsed"
      >
        <Menu/>
      </n-layout-sider>
      <n-scrollbar style="max-height: calc(100vh - 60px)">
        <slot></slot>
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