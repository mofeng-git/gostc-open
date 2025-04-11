<script setup>
import {NIcon} from "naive-ui"
import {computed, h, nextTick, onBeforeMount, ref, watch} from "vue";
import {localStore} from "../../store/local.js";
import {apiConfigQuery} from "../../api/config.js";


const state = ref({
  data:{
    version:'v1.0.0'
  }
})

const configQueryFunc =async ()=>{
  try {
    let res = await apiConfigQuery()
    state.value.data = res.data
  }finally {

  }
}

onBeforeMount(() => {
  configQueryFunc()
})

</script>

<template>
  <div class="header-container">
    <div class="header-title">
      GOSTC WEBUI
    </div>
    <div style="height: 100%;display: flex;align-items: end;padding: 14px 10px;box-sizing: border-box;color: #F1F1F1;">
      <span>{{` ${state.data.version} `}}</span>
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
    white-space: nowrap;
  }
}

@media screen and (max-width: 350px) {
  .header-title {
    display: none;
  }
}
</style>