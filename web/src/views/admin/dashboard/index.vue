<script setup>
import {onBeforeMount, ref} from "vue";
import {
  apiAdminDashboardClientForwardObsDate,
  apiAdminDashboardClientHostObsDate,
  apiAdminDashboardClientObsDate,
  apiAdminDashboardClientTunnelObsDate,
  apiAdminDashboardCount,
  apiAdminDashboardNodeObs,
  apiAdminDashboardNodeObsDate,
  apiAdminDashboardUserObs,
  apiAdminDashboardUserObsDate
} from "../../../api/admin/dashboard.js";
import {flowFormat} from "../../../utils/flow.js";
import AppCard from "../../../layout/components/AppCard.vue";
import Online from "../../../icon/online.vue";
import DatePicker from "../../../components/DatePicker.vue";
import moment from "moment";

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
    checkInTotal: 0,
  },
  userObs: [],
  nodeObs: [],
  userObsDate: [],
  nodeObsDate: [],
  clientObsDate: [],
  clientHostObsDate: [],
  clientForwardObsDate: [],
  clientTunnelObsDate: [],
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

const nodeObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiAdminDashboardNodeObsDate(date)
    state.value.nodeObsDate = res.data || []
  } finally {

  }
}
const userObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiAdminDashboardUserObsDate(date)
    state.value.userObsDate = res.data || []
  } finally {

  }
}

const clientObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiAdminDashboardClientObsDate(date)
    state.value.clientObsDate = res.data || []
  } finally {

  }
}

const clientHostObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiAdminDashboardClientHostObsDate(date)
    state.value.clientHostObsDate = res.data || []
  } finally {

  }
}

const clientForwardObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiAdminDashboardClientForwardObsDate(date)
    state.value.clientForwardObsDate = res.data || []
  } finally {

  }
}

const clientTunnelObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiAdminDashboardClientTunnelObsDate(date)
    state.value.clientTunnelObsDate = res.data || []
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
  let date = new Date()
  nodeObsDateFunc(date)
  userObsDateFunc(date)
  clientObsDateFunc(date)
  clientHostObsDateFunc(date)
  clientForwardObsDateFunc(date)
  clientTunnelObsDateFunc(date)
  countFunc()
})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-grid cols="1 300:2 600:3 800:4 1000:6" x-gap="12" y-gap="12">
        <n-grid-item>
          <n-card size="small" title="今日签到(人)" :content-style="{padding:'6px 16px !important'}">
            <div style="font-size: 1.5em">
              <span style="font-weight: bold">{{ state.count.checkInTotal }}</span>
            </div>
          </n-card>
        </n-grid-item>
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
          <n-card size="small" :content-style="{padding:'6px 16px !important'}">
            <template #header>
              <n-space justify="space-between" align="center">
                <span>节点流量排行 ( IN | OUT )</span>
                <DatePicker :max-value="new Date()" @on-before="args => nodeObsDateFunc(args)"
                            @on-after="args => nodeObsDateFunc(args)"></DatePicker>
              </n-space>
            </template>
            <n-scrollbar style="height: 400px">
              <n-list v-if="state.nodeObsDate?.length">
                <n-list-item v-for="(obs,index) in state.nodeObsDate">
                  <div style="display: flex;justify-content: space-between;align-items: center">
                    <div style="width: 50px">{{ index + 1 }}</div>
                    <div style="flex: 1">
                      <Online :online="obs.online === 1"></Online>
                      {{ obs.name }}
                    </div>
                    <div>
                      <span>{{ flowFormat(obs.inputBytes) }}</span>
                      <n-divider vertical></n-divider>
                      <span>{{ flowFormat(obs.outputBytes) }}</span>
                    </div>
                  </div>
                </n-list-item>
              </n-list>
              <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
            </n-scrollbar>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" :content-style="{padding:'6px 16px !important'}">
            <template #header>
              <n-space justify="space-between" align="center">
                <span>用户流量排行 ( IN | OUT )</span>
                <DatePicker :max-value="new Date()" @on-before="args => userObsDateFunc(args)"
                            @on-after="args => userObsDateFunc(args)"></DatePicker>
              </n-space>
            </template>
            <n-scrollbar style="height: 400px">
              <n-list v-if="state.userObsDate?.length">
                <n-list-item v-for="(obs,index) in state.userObsDate">
                  <div style="display: flex;justify-content: space-between;align-items: center">
                    <div style="width: 50px">{{ index + 1 }}</div>
                    <div style="flex: 1">{{ obs.account }}</div>
                    <div>
                      <span>{{ flowFormat(obs.inputBytes) }}</span>
                      <n-divider vertical></n-divider>
                      <span>{{ flowFormat(obs.outputBytes) }}</span>
                    </div>
                  </div>
                </n-list-item>
              </n-list>
              <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
            </n-scrollbar>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" :content-style="{padding:'6px 16px !important'}">
            <template #header>
              <n-space justify="space-between" align="center">
                <span>客户端流量排行 ( IN | OUT )</span>
                <DatePicker :max-value="new Date()" @on-before="args => clientObsDateFunc(args)"
                            @on-after="args => clientObsDateFunc(args)"></DatePicker>
              </n-space>
            </template>
            <n-scrollbar style="height: 400px">
              <n-list v-if="state.clientObsDate?.length">
                <n-list-item v-for="(obs,index) in state.clientObsDate">
                  <div style="display: flex;justify-content: space-between;align-items: center">
                    <div style="width: 50px">{{ index + 1 }}</div>
                    <div style="flex: 1">
                      <Online :online="obs.online === 1"></Online>
                      {{ obs.name }}
                    </div>
                    <div>
                      <span>{{ flowFormat(obs.inputBytes) }}</span>
                      <n-divider vertical></n-divider>
                      <span>{{ flowFormat(obs.outputBytes) }}</span>
                    </div>
                  </div>
                </n-list-item>
              </n-list>
              <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
            </n-scrollbar>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" :content-style="{padding:'6px 16px !important'}">
            <template #header>
              <n-space justify="space-between" align="center">
                <span>域名解析流量排行 ( IN | OUT )</span>
                <DatePicker :max-value="new Date()" @on-before="args => clientHostObsDateFunc(args)"
                            @on-after="args => clientHostObsDateFunc(args)"></DatePicker>
              </n-space>
            </template>
            <n-scrollbar style="height: 400px">
              <n-list v-if="state.clientHostObsDate?.length">
                <n-list-item v-for="(obs,index) in state.clientHostObsDate">
                  <div style="display: flex;justify-content: space-between;align-items: center">
                    <div style="width: 50px">{{ index + 1 }}</div>
                    <div style="flex: 1">
                      {{ obs.name }}
                    </div>
                    <div>
                      <span>{{ flowFormat(obs.inputBytes) }}</span>
                      <n-divider vertical></n-divider>
                      <span>{{ flowFormat(obs.outputBytes) }}</span>
                    </div>
                  </div>
                </n-list-item>
              </n-list>
              <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
            </n-scrollbar>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card size="small" :content-style="{padding:'6px 16px !important'}">
            <template #header>
              <n-space justify="space-between" align="center">
                <span>端口转发流量排行 ( IN | OUT )</span>
                <DatePicker :max-value="new Date()" @on-before="args => clientForwardObsDateFunc(args)"
                            @on-after="args => clientForwardObsDateFunc(args)"></DatePicker>
              </n-space>
            </template>
            <n-scrollbar style="height: 400px">
              <n-list v-if="state.clientForwardObsDate?.length">
                <n-list-item v-for="(obs,index) in state.clientForwardObsDate">
                  <div style="display: flex;justify-content: space-between;align-items: center">
                    <div style="width: 50px">{{ index + 1 }}</div>
                    <div style="flex: 1">
                      {{ obs.name }}
                    </div>
                    <div>
                      <span>{{ flowFormat(obs.inputBytes) }}</span>
                      <n-divider vertical></n-divider>
                      <span>{{ flowFormat(obs.outputBytes) }}</span>
                    </div>
                  </div>
                </n-list-item>
              </n-list>
              <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
            </n-scrollbar>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card size="small" :content-style="{padding:'6px 16px !important'}">
            <template #header>
              <n-space justify="space-between" align="center">
                <span>私有隧道流量排行 ( IN | OUT )</span>
                <DatePicker :max-value="new Date()" @on-before="args => clientTunnelObsDateFunc(args)"
                            @on-after="args => clientTunnelObsDateFunc(args)"></DatePicker>
              </n-space>
            </template>
            <n-scrollbar style="height: 400px">
              <n-list v-if="state.clientTunnelObsDate?.length">
                <n-list-item v-for="(obs,index) in state.clientTunnelObsDate">
                  <div style="display: flex;justify-content: space-between;align-items: center">
                    <div style="width: 50px">{{ index + 1 }}</div>
                    <div style="flex: 1">
                      {{ obs.name }}
                    </div>
                    <div>
                      <span>{{ flowFormat(obs.inputBytes) }}</span>
                      <n-divider vertical></n-divider>
                      <span>{{ flowFormat(obs.outputBytes) }}</span>
                    </div>
                  </div>
                </n-list-item>
              </n-list>
              <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
            </n-scrollbar>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card size="small" title="30天节点流量排行 ( IN | OUT )" :content-style="{padding:'6px 16px !important'}">
            <n-scrollbar style="height: 400px">
              <n-list v-if="state.nodeObs?.length">
                <n-list-item v-for="(obs,index) in state.nodeObs">
                  <div style="display: flex;justify-content: space-between;align-items: center">
                    <div style="width: 50px">{{ index + 1 }}</div>
                    <div style="flex: 1">
                      <Online :online="obs.online === 1"></Online>
                      {{ obs.name }}
                    </div>
                    <div>
                      <span>{{ flowFormat(obs.inputBytes) }}</span>
                      <n-divider vertical></n-divider>
                      <span>{{ flowFormat(obs.outputBytes) }}</span>
                    </div>
                  </div>
                </n-list-item>
              </n-list>
              <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
            </n-scrollbar>
          </n-card>
        </n-grid-item>
        <n-grid-item>
          <n-card size="small" title="30天用户流量排行 ( IN | OUT )" :content-style="{padding:'6px 16px !important'}">
            <n-scrollbar style="height: 400px">
              <n-list v-if="state?.userObs">
                <n-list-item v-for="(obs,index) in state.userObs">
                  <div style="display: flex;justify-content: space-between;align-items: center">
                    <div style="width: 50px">{{ index + 1 }}</div>
                    <div style="flex: 1">{{ obs.account }}</div>
                    <div>
                      <span>{{ flowFormat(obs.inputBytes) }}</span>
                      <n-divider vertical></n-divider>
                      <span>{{ flowFormat(obs.outputBytes) }}</span>
                    </div>
                  </div>
                </n-list-item>
              </n-list>
              <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
            </n-scrollbar>
          </n-card>
        </n-grid-item>
      </n-grid>
    </AppCard>
  </div>
</template>

<style scoped lang="scss">

</style>