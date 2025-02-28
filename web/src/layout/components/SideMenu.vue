<script setup>
import {computed, h, nextTick, ref, watch} from "vue";
import MenuIconSvg from "./MenuIconSvg.vue";
import router, {allRouters} from "../../router/index.js";
import {localStore} from "../../store/local.js";
import {appStore} from "../../store/app.js";
import {routerToMenu} from "../../utils/routerToMenu.js";
import {normalRouters} from "../../router/routers/normal.js";


const menuTreeDataComputed = computed(() => {
  if (appStore().userInfo.admin === 1) {
    return routerToMenu(allRouters)
  }
  return routerToMenu(normalRouters)
})

const renderMenuIcon = (option) => {
  return h(MenuIconSvg, {
    svg: option.iconSvg
  }, null);
}

const menuSelectChange = (key, item) => {
  if (item?.link) {
    window.open(item.link, '_blank')
    return
  }
  if (router.currentRoute.value.name === key) {
    return
  }
  localStore().menuKey = key
  router.push({
    name: key
  })
}

const onMaskClick = () => {
  localStore().isCollapsed = true
}
const sideContainer = ref()

const menu = ref()
watch(router.currentRoute, () => {
  nextTick(() => {
    menu.value?.showOption()
  })
})


</script>

<template>
  <div ref="sideContainer">
    <n-drawer :show="!localStore().isCollapsed"
              width="80%"
              placement="left"
              :trap-focus="false"
              :block-scroll="false"
              :on-mask-click="onMaskClick"
              :to="sideContainer"
    >
      <n-drawer-content :title="appStore().siteConfig.title">
        <n-scrollbar style="height: calc(100vh - 60px)">
          <n-menu
              accordion
              ref="menu"
              :options="menuTreeDataComputed"
              :render-icon="renderMenuIcon"
              :on-update:value="menuSelectChange"
              :value="localStore().menuKey"
          />
        </n-scrollbar>
      </n-drawer-content>
    </n-drawer>
  </div>

</template>

<style scoped lang="scss">
:deep(.n-menu-item-content--selected)::before {

}

//:deep(.n-menu-item-content--child-active):not(.n-menu-item-content--collapsed)::before {
:deep(.n-menu-item-content--child-active.n-menu-item-content--collapsed)::before {

}

:deep(.n-drawer .n-drawer-content .n-drawer-body-content-wrapper) {
  box-sizing: border-box;
  padding: 0 !important;
}

:deep(.n-drawer-header) {
  height: 60px !important;
  box-sizing: border-box !important;
}
</style>