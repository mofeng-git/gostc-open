<script setup>
import MenuCollapsed from "./MenuCollapsed.vue";
import {NButton, NIcon} from "naive-ui"
import {appStore} from "../../store/app.js";
import {ref} from "vue";
import Drawer from "../../components/Drawer.vue";
import AppCard from "./AppCard.vue";
import {localStore} from "../../store/local.js";


const state = ref({
  pwd: {
    formData: {
      oldPwd: '',
      newPwd: '',
    },
    loading: false,
  },
  notifyOpen: false,
})

const logoutFunc = () => {
  localStore().auth.token = ''
  localStore().auth.tokenExpAt = ''
  window.location.reload()
}

</script>

<template>
  <div class="header-container">
    <MenuCollapsed/>
    <n-divider vertical/>
    <div class="header-title">
      {{ appStore().siteConfig.title }}
    </div>
    <div style="flex: 1"></div>
    <!--  夜间模式  -->
    <NIcon style="margin: 5px;user-select: none;cursor: pointer" :size="20"
           @click="localStore().darkTheme = !localStore().darkTheme">
      <svg v-if="localStore().darkTheme" style="color: #F1F1F1" xmlns="http://www.w3.org/2000/svg"
           xmlns:xlink="http://www.w3.org/1999/xlink"
           viewBox="0 0 24 24">
        <g fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="4"></circle>
          <path d="M3 12h1m8-9v1m8 8h1m-9 8v1M5.6 5.6l.7.7m12.1-.7l-.7.7m0 11.4l.7.7m-12.1-.7l-.7.7"></path>
        </g>
      </svg>
      <svg v-else style="color: #F1F1F1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
           viewBox="0 0 24 24">
        <path
            d="M9.5 4c4.41 0 8 3.59 8 8s-3.59 8-8 8c-.34 0-.68-.02-1.01-.07c1.91-2.16 3.01-4.98 3.01-7.93s-1.1-5.77-3.01-7.93C8.82 4.02 9.16 4 9.5 4m0-2c-1.82 0-3.53.5-5 1.35c2.99 1.73 5 4.95 5 8.65s-2.01 6.92-5 8.65c1.47.85 3.18 1.35 5 1.35c5.52 0 10-4.48 10-10S15.02 2 9.5 2z"
            fill="currentColor"></path>
      </svg>
    </NIcon>
<!--    <n-divider vertical/>-->
    <!--  通知  -->
<!--    <NIcon style="margin: 5px;user-select: none;cursor: pointer" :size="20" @click="state.notifyOpen = true">-->
<!--      <svg t="1731032347646" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"-->
<!--           p-id="4502" width="200" height="200">-->
<!--        <path-->
<!--            d="M512.64 0a42.666667 42.666667 0 0 1 42.56 39.466667l0.128 3.2-0.021333 49.92c173.226667 15.573333 307.584 176.576 307.584 370.773333v126.72c0 54.954667 18.176 107.733333 50.602666 147.733333l17.408 21.482667c37.738667 46.549333 8.277333 121.706667-53.610666 121.706667h-188.202667a180.074667 180.074667 0 0 1-345.109333 0H155.733333c-61.888 0-91.349333-75.178667-53.589333-121.706667l17.408-21.504c32.426667-39.978667 50.581333-92.757333 50.581333-147.712v-126.72c0-191.296 130.389333-350.357333 299.84-369.962667V42.666667a42.666667 42.666667 0 0 1 42.666667-42.666667z m83.413333 881.002667h-159.082666a94.656 94.656 0 0 0 159.104 0zM527.765333 176.704H512.64l-0.469333-0.021333h-6.890667c-134.741333 0-246.485333 123.605333-249.728 279.68l-0.085333 7.018666v126.698667c0 74.261333-24.704 146.048-69.610667 201.429333l-3.392 4.138667h668.074667l-3.328-4.117333c-43.2-53.269333-67.712-121.685333-69.546667-192.896l-0.106667-8.554667v-126.72c0-159.36-113.045333-286.656-249.813333-286.656z"-->
<!--            fill="#F1F1F1" p-id="4503"></path>-->
<!--      </svg>-->
<!--    </NIcon>-->
<!--    <n-divider vertical/>-->
    <!--  头像  -->
<!--    <div class="header-avatar">-->
<!--      <n-popover trigger="click">-->
<!--        <template #trigger>-->
<!--          <n-icon :size="20">-->
<!--            <svg t="1717400206462" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"-->
<!--                 p-id="4269" width="200" height="200">-->
<!--              <path-->
<!--                  d="M908.6432 802.848q0 68.576-41.728 108.288t-110.848 39.712l-499.424 0q-69.152 0-110.848-39.712t-41.728-108.288q0-30.272 2.016-59.136t8-62.272 15.136-62.016 24.576-55.712 35.424-46.272 48.864-30.56 63.712-11.424q5.152 0 24 12.288t42.56 27.424 61.728 27.424 76.288 12.288 76.288-12.288 61.728-27.424 42.56-27.424 24-12.288q34.848 0 63.712 11.424t48.864 30.56 35.424 46.272 24.576 55.712 15.136 62.016 8 62.272 2.016 59.136zM725.7952 292.576q0 90.848-64.288 155.136t-155.136 64.288-155.136-64.288-64.288-155.136 64.288-155.136 155.136-64.288 155.136 64.288 64.288 155.136z"-->
<!--                  p-id="4270" fill="#F1F1F1"></path>-->
<!--            </svg>-->
<!--          </n-icon>-->
<!--        </template>-->
<!--        <div style="margin: 20px 0">-->
<!--          <div style="display: flex;justify-content: center;align-items: center">-->
<!--            <n-avatar-->
<!--                :size="72"-->
<!--                round-->
<!--                style="color:#F1F1F1;background-color: #bfbfbf;font-size:2rem;margin-bottom: 20px"-->
<!--            >G-->
<!--            </n-avatar>-->
<!--          </div>-->
<!--          <div-->
<!--              style="padding:0 8px;height:64px;flex: 1;display: flex;flex-direction: column;justify-content: space-between">-->
<!--            <n-ellipsis :style="{maxWidth: appStore().drawerWidthAdapter - 120 + 'px',fontSize: '1.5em'}">-->
<!--              {{ appStore().userInfo.email }}-->
<!--            </n-ellipsis>-->
<!--            <n-space>-->
<!--              <n-button size="small" type="error" :focusable="false" tertiary @click="logoutFunc">退出登录</n-button>-->
<!--            </n-space>-->
<!--          </div>-->
<!--        </div>-->
<!--      </n-popover>-->
<!--    </div>-->

    <Drawer :show="state.notifyOpen" :width="appStore().drawerWidthAdapter">
      <template #header>
        <div style="display: flex;justify-content: space-between;align-items: center">
          <span>消息</span>
          <n-button type="error" tertiary :focusable="false" size="small" @click="state.notifyOpen = false">关闭
          </n-button>
        </div>
      </template>
<!--      <AppCard :show-border="false" :empty="useStore.app().siteConfig.notifyList.length===0">-->
<!--        <n-alert-->
<!--            v-for="item in appStore().siteConfig.notifyList"-->
<!--            style="margin-bottom: 8px"-->
<!--            :show-icon="false">-->
<!--          <template #header>-->
<!--            <div style="display: flex;justify-content: space-between;align-items: center">-->
<!--              <span>{{ item.title }}</span>-->
<!--              <span>{{ item.date }}</span>-->
<!--            </div>-->
<!--          </template>-->
<!--          {{ item.content }}-->
<!--        </n-alert>-->
<!--      </AppCard>-->
    </Drawer>
  </div>
</template>

<style lang="scss" scoped>
.header-container {
  padding: 0 10px;
  background-color: #303133;
  height: 60px;
  display: flex;
  align-items: center;

  .header-title {
    color: #F1F1F1;
    font-size: 1.8em;
  }

  .header-avatar {
    cursor: pointer;
    height: 30px;
    padding: 5px;
    box-sizing: border-box;
  }
}

@media screen and (max-width: 350px) {
  .header-title {
    display: none;
  }
}
</style>