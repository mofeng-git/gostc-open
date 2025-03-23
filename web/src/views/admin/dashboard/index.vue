<script setup>
import {onBeforeMount, ref} from "vue";
import {
  apiAdminDashboardCount,
  apiAdminDashboardNodeObs,
  apiAdminDashboardUserObs
} from "../../../api/admin/dashboard.js";
import {flowFormat} from "../../../utils/flow.js";
import AppCard from "../../../layout/components/AppCard.vue";
import Online from "../../../icon/online.vue";

const state = ref({
  count: {
    client: 0,
    clientOnline: 0,
    node: 0,
    nodeOnline: 0,
    host: 0,
    forward: 0,
    tunnel: 0,
    user: 0,
    inputBytes: 0,
    outputBytes: 0,
  },
  userObs: [],
  nodeObs: [],
})

const nodeObsFunc = async () => {
  try {
    let res = await apiAdminDashboardNodeObs()
    state.value.nodeObs = res.data || []
  } finally {

  }
}
const userObsFunc = async () => {
  try {
    let res = await apiAdminDashboardUserObs()
    state.value.userObs = res.data || []
  } finally {

  }
}
const countFunc = async () => {
  try {
    let res = await apiAdminDashboardCount()
    state.value.count = res.data
  } finally {

  }
}

onBeforeMount(() => {
  nodeObsFunc()
  userObsFunc()
  countFunc()
})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-grid cols="1 350:2 600:3 800:4 1000:6" x-gap="12" y-gap="12">
        <n-grid-item>
          <n-card size="small" title="节点" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.nodeOnline }}</span>
              <n-divider vertical></n-divider>
              <span style="opacity: 0.8">{{ state.count.node }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="客户端" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.clientOnline }}</span>
              <n-divider vertical></n-divider>
              <span style="opacity: 0.8">{{ state.count.client }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="用户" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.user }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="域名解析" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.host }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="端口转发" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.forward }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="私有隧道" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.tunnel }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="代理隧道" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.proxy }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="P2P隧道" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.p2p }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="下行流量(今日)" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ flowFormat(state.count.inputBytes) }}</span>
            </div>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="上行流量(今日)" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ flowFormat(state.count.outputBytes) }}</span>
            </div>
          </n-card>
        </n-grid-item>
      </n-grid>

      <p></p>
      <n-grid cols="1 500:2" x-gap="12" y-gap="12">
        <n-grid-item>
          <n-card size="small" title="30天节点流量排行 ( IN | OUT )" :content-style="{padding:'6px 16px !important'}">
            <n-list>
              <n-list-item v-for="obs in state.nodeObs">
                <div style="display: flex;justify-content: space-between;align-items: center">
                  <span><Online :online="obs.online === 1"></Online> {{ obs.name }}</span>
                  <div>
                    <span>{{ flowFormat(obs.inputBytes) }}</span>
                    <n-divider vertical></n-divider>
                    <span>{{ flowFormat(obs.outputBytes) }}</span>
                  </div>
                </div>
              </n-list-item>
            </n-list>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="30天用户流量排行 ( IN | OUT )" :content-style="{padding:'6px 16px !important'}">
            <n-list>
              <n-list-item v-for="obs in state.userObs">
                <div style="display: flex;justify-content: space-between;align-items: center">
                  <span>{{ obs.account }}</span>
                  <div>
                    <span>{{ flowFormat(obs.inputBytes) }}</span>
                    <n-divider vertical></n-divider>
                    <span>{{ flowFormat(obs.outputBytes) }}</span>
                  </div>
                </div>
              </n-list-item>
            </n-list>
          </n-card>
        </n-grid-item>
      </n-grid>
    </AppCard>
  </div>
</template>

<style scoped lang="scss">

</style>