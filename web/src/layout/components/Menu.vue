<script setup>
import {computed, h, nextTick, ref, watch} from "vue";
import router, {allRouters} from "../../router/index.js";
import MenuIconSvg from "./MenuIconSvg.vue";
import {localStore} from "../../store/local.js";
import {routerToMenu} from "../../utils/routerToMenu.js";
import {appStore} from "../../store/app.js";
import {normalRouters} from "../../router/routers/normal.js";

const menuTreeDataComputed = computed(() => {
  var funcMap = new Map()
  funcMap.set("funcWeb",appStore().siteConfig.funcWeb)
  funcMap.set("funcForward",appStore().siteConfig.funcForward)
  funcMap.set("funcTunnel",appStore().siteConfig.funcTunnel)
  funcMap.set("funcP2P",appStore().siteConfig.funcP2P)
  funcMap.set("funcProxy",appStore().siteConfig.funcProxy)
  funcMap.set("funcTun",appStore().siteConfig.funcTun)
  funcMap.set("funcNode",appStore().siteConfig.funcNode)
  if (appStore().userInfo.admin === 1) {
    return routerToMenu(allRouters,funcMap)
  }
  return routerToMenu(normalRouters,funcMap)
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

const menu = ref()
watch(router.currentRoute, () => {
  nextTick(() => {
    menu.value?.showOption()
  })
})

</script>

<template>
  <n-scrollbar style="height: calc(100vh - 60px)">
    <n-menu
        accordion
        ref="menu"
        :collapsed="localStore().isCollapsed"
        :collapsed-width="60"
        :collapsed-icon-size="22"
        :options="menuTreeDataComputed"
        :render-icon="renderMenuIcon"
        :on-update:value="menuSelectChange"
        :value="localStore().menuKey"
    />
  </n-scrollbar>
</template>

<style scoped lang="scss">
:deep(.n-menu-item-content--selected)::before {
  //background-color: var(--n-item-color-active);
  //border-left: 4px solid #2299dd;
  //background-color: #2299dd !important;
}

//:deep(.n-menu-item-content--child-active):not(.n-menu-item-content--collapsed)::before {
:deep(.n-menu-item-content--child-active.n-menu-item-content--collapsed)::before {
  //background-color: var(--n-item-color-active) !important;
}
</style>