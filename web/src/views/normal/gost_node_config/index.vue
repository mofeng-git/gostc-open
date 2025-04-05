<script setup>
import {onBeforeMount, ref, watch} from "vue";
import AppCard from "../../../layout/components/AppCard.vue";
import Empty from "../../../components/Empty.vue";
import {NSpace} from "naive-ui";
import {cLimiterText, configText, limiterText, rLimiterText} from "./index.js";
import {apiNormalGostNodeList} from "../../../api/normal/gost_node.js";
import Online from "../../../icon/online.vue";

const state = ref({
  search: {
    bind: 0,
  },
  nodes: [],
  node: {},
  nodeLoading: false,
  configs: [],
  selectNodeCode: '',
})

const nodeListFunc = async () => {
  try {
    let res = await apiNormalGostNodeList(state.value.search)
    state.value.nodes = res.data || []
    if (state.value.nodes.length > 0) {
      state.value.selectNodeCode = state.value.nodes[0].code
      state.value.node = state.value.nodes[0]
      state.value.configs = state.value.node?.configs
    } else {
      $message.create('暂无节点套餐信息', {
        type: "warning",
        closable: true,
        duration: 1500,
      })
    }
  } finally {

  }
}

watch(() => ({
  bind: state.value.search.bind
}), () => {
  nodeListFunc()
})


watch(() => ({
  code: state.value.selectNodeCode,
}), (value) => {
  for (let i = 0; i < state.value.nodes.length; i++) {
    if (state.value.selectNodeCode === state.value.nodes[i].code) {
      state.value.node = state.value.nodes[i]
      state.value.configs = state.value.node?.configs
      break
    }
  }
})

onBeforeMount(() => {
  nodeListFunc()
})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-h4 style="font-weight: bold">节点类型：</n-h4>
      <n-radio-group v-model:value="state.search.bind" size="medium" :default-value="state.search.bind">
        <n-radio-button
            v-for="tp in [{label:'全部',value:0},{label:'我的节点',value:1},{label:'系统节点',value:2}]"
            :key="tp.value"
            :value="tp.value">
          {{ tp.label }}
        </n-radio-button>
      </n-radio-group>
      <n-h4 style="font-weight: bold">选择节点：</n-h4>
      <n-spin :show="state.nodeLoading">
        <n-radio-group v-model:value="state.selectNodeCode" name="nodeRadioGroup"
                       style="width: 100%">
          <n-grid x-gap="12" y-gap="12" cols="550:2 800:3 1400:4 1">
            <n-grid-item v-for="nodeItem in state.nodes">
              <n-alert
                  type="info"
                  :show-icon="false"
                  :bordered="false"
                  style="height: 100%;cursor: pointer"
                  @click="state.selectNodeCode = nodeItem.code"
              >
                <n-radio
                    :key="nodeItem.code"
                    :value="nodeItem.code"
                    style="width: 100%;"
                >
                  <n-space justify="space-between" style="width: 100%">
                    <Online :online="nodeItem.online===1"></Online>
                    <span>{{ nodeItem.name }}</span>
                  </n-space>
                </n-radio>
              </n-alert>
            </n-grid-item>
          </n-grid>
        </n-radio-group>

        <n-h4 style="font-weight: bold">节点信息：</n-h4>
        <n-space vertical>
          <n-space>
            <span>功能：</span>
            <n-tag type="info" size="small" bordered v-if="state.node.web===1">域名解析</n-tag>
            <n-tag type="info" size="small" bordered v-if="state.node.forward===1">端口转发</n-tag>
            <n-tag type="info" size="small" bordered v-if="state.node.tunnel===1">私有隧道</n-tag>
            <n-tag type="info" size="small" bordered v-if="state.node.proxy===1">代理隧道</n-tag>
            <n-tag type="info" size="small" bordered v-if="state.node.p2p===1">P2P隧道</n-tag>
          </n-space>
          <n-space>
            <span>标签：</span>
            <n-tag type="info" size="small" bordered v-for="tag in state.node.tags">{{ tag }}</n-tag>
          </n-space>
          <n-space>
            <span>规则：</span>
            <n-tag type="info" size="small" bordered v-for="rule in state.node.ruleNames">{{ rule }}</n-tag>
          </n-space>
          <n-space>
            <span>自定义域名：</span>
            <span>{{ state.node.customDomain === 1 ? '支持' : '不支持' }}</span>
          </n-space>
          <n-space>
            <span>介绍：</span>
            <span>{{ state.node.remark }}</span>
          </n-space>
        </n-space>

        <n-h4 style="font-weight: bold">套餐信息：</n-h4>
        <p></p>
        <Empty v-if="state.configs.length===0" border description="该节点暂无套餐"></Empty>
        <n-grid v-else x-gap="12" y-gap="12" cols="1 520:2 900:3 1400:4">
          <n-grid-item v-for="row in state.configs">
            <n-el class="client-item" tag="div" :style="{
                border: '1px solid var(--border-color)',
                borderRadius:'var(--border-radius)',
                padding: '12px',
                cursor: 'pointer'}">
              <n-h4 style="margin-bottom: 8px !important;">
                <n-space justify="space-between">
                  <span style="font-weight: bold">{{ row.name }}</span>
                </n-space>
              </n-h4>
              <div style="display: flex;justify-content: space-between">
                <span style="font-weight: bold;">说明：</span>
                <span>{{ configText(row) }}</span>
              </div>
              <div style="display: flex;justify-content: space-between">
                <span style="font-weight: bold;">速率：</span>
                <span>{{ limiterText(row.limiter) }}</span>
              </div>
              <div style="display: flex;justify-content: space-between">
                <span style="font-weight: bold;">并发数：</span>
                <span>{{ rLimiterText(row.rLimiter) }}</span>
              </div>
              <div style="display: flex;justify-content: space-between">
                <span style="font-weight: bold;">连接数：</span>
                <span>{{ cLimiterText(row.cLimiter) }}</span>
              </div>
            </n-el>
          </n-grid-item>
        </n-grid>
      </n-spin>
    </AppCard>
  </div>
</template>

<style scoped lang="scss">

</style>