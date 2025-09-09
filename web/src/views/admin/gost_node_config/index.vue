<script setup>
import {onBeforeMount, ref, watch} from "vue";
import {
  apiAdminGostNodeConfigCreate,
  apiAdminGostNodeConfigDelete,
  apiAdminGostNodeConfigList,
  apiAdminGostNodeConfigUpdate
} from "../../../api/admin/gost_node_config.js";
import AppCard from "../../../layout/components/AppCard.vue";
import {requiredRule} from "../../../utils/formDataRule.js";
import {apiAdminGostNodeList, apiAdminGostNodeQuery} from "../../../api/admin/gost_node.js";
import Modal from "../../../components/Modal.vue";
import Empty from "../../../components/Empty.vue";
import {NButton, NPopconfirm, NSpace} from "naive-ui";
import {cLimiterText, configText, limiterText, rLimiterText} from "./index.js";
import router from "../../../router/index.js";
import Online from "../../../icon/online.vue";

const state = ref({
  create: {
    open: false,
    loading: false,
    data: {
      name: '',
      chargingType: 1,
      cycle: 10,
      amount: '0',
      limiter: 1,
      rLimiter: 100,
      cLimiter: 100,
      nodeCode: '',
      indexValue: 1000,
    },
    dataRules: {
      name: requiredRule('请输入名称'),
    },
  },
  update: {
    open: false,
    loading: false,
    data: {
      code: '',
      name: '',
      chargingType: 1,
      cycle: 10,
      amount: '0',
      limiter: 1,
      rLimiter: 100,
      cLimiter: 100,
      nodeCode: '',
      indexValue: 1000,
    },
    dataRules: {
      name: requiredRule('请输入名称'),
    },
  },
  search: {
    bind: 2
  },
  nodes: [],
  node: {},
  nodeLoading: false,
  configs: [],
  selectNodeCode: '',
})

const openCreate = () => {
  state.value.create.data = {
    name: '',
    chargingType: 1,
    cycle: 10,
    amount: '0',
    limiter: 1,
    rLimiter: 100,
    cLimiter: 100,
    nodeCode: state.value.node.code,
    indexValue: 1000,
  }
  state.value.create.open = true
}

const closeCreate = () => {
  state.value.create.open = false
}

const openUpdate = (row) => {
  state.value.update.data = JSON.parse(JSON.stringify(row))
  state.value.update.open = true
}

const closeUpdate = () => {
  state.value.update.open = false
}

const createRef = ref()
const createFunc = () => {
  createRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.create.loading = true
        await apiAdminGostNodeConfigCreate(state.value.create.data)
        await listFunc()
        closeCreate()
      } finally {
        state.value.create.loading = false
      }
    }
  })
}

const updateRef = ref()
const updateFunc = () => {
  updateRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.update.loading = true
        await apiAdminGostNodeConfigUpdate(state.value.update.data)
        await listFunc()
        closeUpdate()
      } finally {
        state.value.update.loading = false
      }
    }
  })
}

const deleteFunc = async (row) => {
  try {
    await apiAdminGostNodeConfigDelete({code: row.code})
    await listFunc()
  } finally {
  }
}

const getNodes = async () => {
  try {
    let res = await apiAdminGostNodeList(state.value.search)
    state.value.nodes = res.data || []
    if (state.value.nodes.length > 0) {
      state.value.selectNodeCode = state.value.nodes[0].code
    } else {
      $message.create('请先新增节点', {
        type: "warning",
        closable: true,
        duration: 1500,
      })
      await router.push({name: 'AdminGostNode'})
    }
  } finally {

  }
}


watch(() => ({
  code: state.value.selectNodeCode,
}), (value) => {
  queryNode()
})

const queryNode = async () => {
  try {
    state.value.nodeLoading = true
    let res = await apiAdminGostNodeQuery({code: state.value.selectNodeCode})
    state.value.node = res.data
    await listFunc()
  } finally {
    state.value.nodeLoading = false
  }
}

const listFunc = async () => {
  try {
    let res = await apiAdminGostNodeConfigList({nodeCode: state.value.node.code})
    state.value.configs = res.data || []
  } finally {
  }
}

watch(() => ({bind: state.value.search.bind}), () => {
  getNodes()
})

onBeforeMount(() => {
  getNodes()
})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-h4 style="font-weight: bold">节点类型：</n-h4>
      <n-radio-group v-model:value="state.search.bind" size="medium" :default-value="state.search.bind">
        <n-radio-button
            v-for="tp in [{label:'全部',value:0},{label:'用户节点',value:1},{label:'系统节点',value:2}]"
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
        <!--        <n-radio-group v-model:value="state.selectNodeCode">-->
        <!--          <n-space>-->
        <!--            <n-radio-->
        <!--                v-for="node in state.nodes"-->
        <!--                :key="node.code"-->
        <!--                :value="node.code">-->
        <!--              {{ node.name }}-->
        <!--            </n-radio>-->
        <!--          </n-space>-->
        <!--        </n-radio-group>-->
        <n-h4 style="font-weight: bold">节点信息：</n-h4>
        <n-space vertical>
          <n-space>
            <span>地址：</span>
            <span>{{ state.node.address }}</span>
          </n-space>
          <n-space>
            <span>协议：</span>
            <span>{{ state.node.protocol }}</span>
          </n-space>
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
            <span>绑定域名：</span>
            <span>{{ state.node.customDomain === 1?'支持':'不支持' }}</span>
          </n-space>
          <n-space>
            <span>介绍：</span>
            <span>{{ state.node.remark }}</span>
          </n-space>
        </n-space>

        <n-h4 style="font-weight: bold">套餐信息：</n-h4>
        <n-button type="info" size="small" :focusable="false" @click="openCreate">新增套餐</n-button>
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
              <n-space justify="end" style="width: 100%">
                <n-button size="tiny" :focusable="false" quaternary type="success" @click="openUpdate(row)">
                  编辑
                </n-button>
                <n-popconfirm
                    @positive-click="deleteFunc(row)"
                    :positive-button-props="{loading:row.deleteLoading}"
                >
                  <template #trigger>
                    <n-button
                        size="tiny"
                        :focusable="false"
                        type="error"
                        quaternary
                    >删除
                    </n-button>
                  </template>
                  确认删除吗？
                </n-popconfirm>
              </n-space>
            </n-el>
          </n-grid-item>
        </n-grid>
      </n-spin>
    </AppCard>

    <Modal
        :show="state.create.open"
        title="新增套餐"
        width="600px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.create.loading"
        @on-confirm="createFunc"
        @on-cancel="closeCreate"
    >
      <n-form ref="createRef" :rules="state.create.dataRules" :model="state.create.data">
        <n-form-item path="name" label="名称">
          <n-input v-model:value.trim="state.create.data.name"></n-input>
        </n-form-item>
        <n-form-item path="chargingType" label="计费方式">
          <n-radio-group v-model:value="state.create.data.chargingType">
            <n-radio :checked="state.create.data.chargingType===1" :value="1">
              一次性
            </n-radio>
            <n-radio :checked="state.create.data.chargingType===2" :value="2">
              循环
            </n-radio>
            <n-radio :checked="state.create.data.chargingType===3" :value="3">
              免费
            </n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item path="cycle" label="续费周期(天)" v-show="state.create.data.chargingType===2">
          <n-input-number v-model:value="state.create.data.cycle" :min="1"></n-input-number>
        </n-form-item>
        <n-form-item path="amount" label="积分" v-show="state.create.data.chargingType!==3">
          <n-input v-model:value.trim="state.create.data.amount"></n-input>
        </n-form-item>
        <n-form-item path="limiter" label="速率(mbps)">
          <n-input-number v-model:value="state.create.data.limiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="indexValue" label="排序(升序)">
          <n-input-number v-model:value="state.create.data.indexValue"></n-input-number>
        </n-form-item>
      </n-form>
    </Modal>

    <Modal
        :show="state.update.open"
        title="修改套餐"
        width="600px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.update.loading"
        @on-confirm="updateFunc"
        @on-cancel="closeUpdate"
        mask-close
    >
      <n-form ref="updateRef" :rules="state.update.dataRules" :model="state.update.data">
        <n-form-item path="name" label="名称">
          <n-input v-model:value.trim="state.update.data.name"></n-input>
        </n-form-item>
        <n-form-item path="chargingType" label="计费方式">
          <n-radio-group v-model:value="state.update.data.chargingType">
            <n-radio :value="1">
              一次性
            </n-radio>
            <n-radio :value="2">
              循环
            </n-radio>
            <n-radio :value="3">
              免费
            </n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item path="cycle" label="续费周期(天)" v-show="state.update.data.chargingType===2">
          <n-input-number v-model:value="state.update.data.cycle" :min="1"></n-input-number>
        </n-form-item>
        <n-form-item path="amount" label="积分" v-show="state.update.data.chargingType!==3">
          <n-input v-model:value.trim="state.update.data.amount"></n-input>
        </n-form-item>
        <n-form-item path="limiter" label="速率(mbps)">
          <n-input-number v-model:value="state.update.data.limiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="indexValue" label="排序(升序)">
          <n-input-number v-model:value="state.update.data.indexValue"></n-input-number>
        </n-form-item>
      </n-form>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>