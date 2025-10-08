<script setup>
import router from "../../../router/index.js";
import {onBeforeMount, ref, watch} from "vue";
import {apiNormalGostClientP2PCreate} from "../../../api/normal/gost_client_p2p.js";
import {regexpRule, requiredRule} from "../../../utils/formDataRule.js";
import {regexpLocalIp, regexpPort} from "../../../utils/regexp.js";
import {apiNormalGostClientList} from "../../../api/normal/gost_client.js";
import {cLimiterText, configText, limiterText, rLimiterText} from "../gost_client_host/index.js";
import AppCard from "../../../layout/components/AppCard.vue";
import {NSpace} from "naive-ui";
import {apiNormalGostNodeList} from "../../../api/normal/gost_node.js";
import Online from "../../../icon/online.vue";


const state = ref({
  search: {
    bind: 0,
  },
  data: {
    configCode: '',
    nodeCode: '',
    clientCode: '',
    name: '',
    targetIp: '',
    targetPort: '',
    useEncryption: 1,
    useCompression: 1,
  },
  dataRules: {
    name: requiredRule('请输入名称'),
    configCode: requiredRule('请选择节点套餐'),
    nodeCode: requiredRule('请选择节点套餐'),
    clientCode: requiredRule('请选择节点套餐'),
    targetIp: regexpRule(regexpLocalIp, '内网IP格式错误'),
    targetPort: regexpRule(regexpPort, '内网端口格式错误'),
  },
  nodes: [],
  node: {},
  loading: true,
  config: {},
  clients: [],
})

const createRef = ref()
const createFunc = () => {
  createRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.loading = true
        await apiNormalGostClientP2PCreate(state.value.data)
        back()
      } finally {
        state.value.loading = false
      }
    }
  })
}

const clientListFunc = async () => {
  try {
    let res = await apiNormalGostClientList()
    state.value.clients = res.data || []
    if (state.value.clients.length > 0) {
      state.value.data.clientCode = state.value.clients[0].code
    } else {
      router.back()
      $message.create('暂无客户端信息', {
        type: "warning",
        closable: true,
        duration: 1500,
      })
    }
  } finally {

  }
}

const nodeListFunc = async () => {
  try {
    state.value.loading = true
    let res = await apiNormalGostNodeList(state.value.search)
    state.value.nodes = res.data || []
    state.value.nodes = state.value.nodes.filter(item=>{
      return item.p2p === 1
    })

    state.value.data.nodeCode = ''
    if (state.value.nodes.length > 0) {
      state.value.data.nodeCode = state.value.nodes[0].code
    }
    if (state.value.data.nodeCode === '') {
      $message.create('暂无节点信息', {
        type: "warning",
        closable: true,
        duration: 1500,
      })
    }
  } finally {
    state.value.loading = false
  }
}

watch(() => ({bind: state.value.search.bind}), () => {
  nodeListFunc()
})

watch(() => ({
  code: state.value.data.nodeCode
}), () => {
  for (let i = 0; i < state.value.nodes.length; i++) {
    if (state.value.nodes[i].code === state.value.data.nodeCode) {
      state.value.node = state.value.nodes[i]
      state.value.data.configCode = state.value.node.configs[0].code
      state.value.config = state.value.node.configs[0]
      break
    }
  }
})
watch(() => ({
  code: state.value.data.configCode
}), () => {
  for (let i = 0; i < state.value.node.configs.length; i++) {
    if (state.value.node.configs[i].code === state.value.data.configCode) {
      state.value.config = state.value.node.configs[i]
      break
    }
  }
})

const back = () => {
  router.back()
}

onBeforeMount(() => {
  nodeListFunc()
  clientListFunc()
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
      <n-spin :show="state.loading">
        <n-radio-group v-model:value="state.data.nodeCode" name="nodeRadioGroup"
                       style="width: 100%">
          <n-grid x-gap="12" y-gap="12" cols="550:2 800:3 1400:4 1">
            <n-grid-item v-for="nodeItem in state.nodes">
              <n-alert
                  type="info"
                  :show-icon="false"
                  :bordered="false"
                  style="height: 100%;cursor: pointer"
                  @click="state.data.nodeCode = nodeItem.code"
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
            <span>介绍：</span>
            <span>{{ state.node.remark }}</span>
          </n-space>
        </n-space>

        <n-h4 style="font-weight: bold">选择套餐：</n-h4>
        <n-radio-group v-model:value="state.data.configCode">
          <n-space>
            <n-radio
                v-for="config in state.node.configs"
                :key="config.code"
                :value="config.code">
              {{ config.name }}
            </n-radio>
          </n-space>
        </n-radio-group>
        <n-h4 style="font-weight: bold">套餐信息：</n-h4>
        <n-space vertical>
          <n-space>
            <span>说明：</span>
            <span>{{ configText(state.config) }}</span>
          </n-space>
          <n-space>
            <span>速率：</span>
            <span>{{ limiterText(state.config.limiter) }}</span>
          </n-space>
        </n-space>

        <n-h4 style="font-weight: bold">选择客户端：</n-h4>
        <n-alert v-if="state.clients.length===0" type="warning">没有客户端，请先新增一个客户端</n-alert>
        <n-radio-group v-else v-model:value="state.data.clientCode"
                       style="width: 100%">
          <n-grid x-gap="12" y-gap="12" cols="550:2 800:3 1400:4 1">
            <n-grid-item v-for="client in state.clients">
              <n-alert
                  type="info"
                  :show-icon="false"
                  :bordered="false"
                  style="height: 100%;cursor: pointer"
                  @click="state.data.clientCode = client.code"
              >
                <n-radio
                    :key="client.code"
                    :value="client.code"
                    style="width: 100%;"
                >
                  <n-space justify="space-between" style="width: 100%">
                    <Online :online="client.online===1"></Online>
                    <span>{{ client.name }}</span>
                  </n-space>
                </n-radio>
              </n-alert>
            </n-grid-item>
          </n-grid>
        </n-radio-group>
      </n-spin>
    </AppCard>

    <AppCard :show-border="false">
      <n-form ref="createRef" :rules="state.dataRules" :model="state.data">
        <n-form-item path="name" label="名称">
          <n-input v-model:value="state.data.name" placeholder="我的服务"></n-input>
        </n-form-item>
        <n-form-item path="targetIp" label="内网IP">
          <n-input v-model:value="state.data.targetIp" placeholder="127.0.0.1"></n-input>
        </n-form-item>
        <n-form-item path="targetPort" label="内网端口">
          <n-input v-model:value="state.data.targetPort" placeholder="80"></n-input>
        </n-form-item>
        <n-form-item label="加密(开启后，会增加一些延迟)">
          <n-select
              :options="[{label:'停用',value:2},{label:'启用',value:1}]"
              v-model:value="state.data.useEncryption"
          ></n-select>
        </n-form-item>
        <n-form-item label="压缩(开启后，会增加一些延迟)">
          <n-select
              :options="[{label:'停用',value:2},{label:'启用',value:1}]"
              v-model:value="state.data.useCompression"
          ></n-select>
        </n-form-item>
      </n-form>
      <n-space>
        <n-button size="small" @click="back" :focusable="false">取消</n-button>
        <n-button size="small" @click="createFunc" type="primary" :focusable="false" :loading="state.loading">
          保存
        </n-button>
      </n-space>
    </AppCard>
  </div>
</template>

<style scoped lang="scss"></style>