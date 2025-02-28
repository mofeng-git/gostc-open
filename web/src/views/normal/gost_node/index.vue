<script setup>
import {onBeforeMount, ref, watch} from "vue";
import {
  apiNormalGostNodeCreate,
  apiNormalGostNodeDelete,
  apiNormalGostNodePage,
  apiNormalGostNodeUpdate
} from "../../../api/normal/gost_node.js";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import Modal from "../../../components/Modal.vue";
import {requiredRule} from "../../../utils/formDataRule.js";
import {protocols} from "./protocols.js"
import {flowFormat} from "../../../utils/flow.js";
import {NButton, NPopconfirm, NSpace} from "naive-ui";
import Empty from "../../../components/Empty.vue";
import {goToUrl} from "../../../utils/browser.js";
import {apiNormalGostObsNodeMonth} from "../../../api/normal/gost_obs.js";
import Obs from "../../../components/Obs.vue";
import moment from "moment";

const state = ref({
  table: {
    searchLoading: false,
    search: {
      page: 1,
      size: 12,
      name: '',
    },
    list: [],
    total: 0,
    minWidth: 2600,
  },
  create: {
    open: false,
    loading: false,
    data: {
      name: '',
      remark: '',
      web: 1,
      tunnel: 1,
      forward: 1,
      domain: '',
      denyDomainPrefix: '',
      address: '',
      protocol: protocols[0].value,
      tunnelConnPort: '',
      tunnelInPort: '',
      tunnelMetadata: '',
      forwardConnPort: '',
      forwardPorts: '',
      forwardMetadata: '',
      tunnelReplaceAddress: '',
      forwardReplaceAddress: '',
      rules: [],
      tags: [],
      indexValue: 1000,
    },
    dataRules: {
      name: requiredRule('请输入名称'),
      address: requiredRule('请输入服务器地址'),
    },
  },
  update: {
    open: false,
    loading: false,
    data: {
      code: '',
      name: '',
      remark: '',
      web: 2,
      tunnel: 2,
      forward: 2,
      domain: '',
      denyDomainPrefix: '',
      address: '',
      protocol: protocols[0].value,
      tunnelConnPort: '',
      tunnelInPort: '',
      tunnelMetadata: '',
      forwardConnPort: '',
      forwardPorts: '',
      forwardMetadata: '',
      tunnelReplaceAddress: '',
      forwardReplaceAddress: '',
      rules: [],
      tags: [],
      indexValue: 1000,
    },
    dataRules: {
      name: requiredRule('请输入名称'),
      address: requiredRule('请输入服务器地址'),
    },
  },
  rules: [],
  obs: {
    open: false,
    code: '',
    loading: false,
    data: [],
    dataRange: 1,
  },
})

const refreshTable = () => {
  pageFunc()
}

const searchTable = () => {
  state.value.table.search.page = 1
  pageFunc()
}

const pageFunc = async () => {
  try {
    state.value.table.searchLoading = true
    let res = await apiNormalGostNodePage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}

const openCreate = () => {
  state.value.create.data = {
    name: '',
    remark: '',
    web: 1,
    tunnel: 1,
    forward: 1,
    domain: '',
    denyDomainPrefix: '',
    address: '',
    protocol: protocols[0].value,
    tunnelConnPort: '',
    tunnelInPort: '',
    tunnelMetadata: '',
    forwardConnPort: '',
    forwardPorts: '',
    forwardMetadata: '',
    tunnelReplaceAddress: '',
    forwardReplaceAddress: '',
    rules: [],
    tags: [],
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
        await apiNormalGostNodeCreate(state.value.create.data)
        searchTable()
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
        await apiNormalGostNodeUpdate(state.value.update.data)
        searchTable()
        closeUpdate()
      } finally {
        state.value.update.loading = false
      }
    }
  })
}

const deleteFunc = async (row) => {
  try {
    await apiNormalGostNodeDelete({code: row.code})
    searchTable()
  } finally {
  }
}

const statusStyle = (row) => {
  return {
    display: 'inline-block',
    fontSize: '0.8em',
    color: row?.online === 1 ? '#0F0' : '#F00',
  }
}

const statusText = (row) => {
  return row?.online === 1 ? '在线' : '离线'
}

const openObsModal = (row) => {
  state.value.obs.code = row.code
  obsFunc()
  state.value.obs.open = true
}

const closeObsModal = () => {
  state.value.obs.open = false
}

const obsFunc = async () => {
  try {
    state.value.obs.loading = false
    state.value.obsLoading = true
    let data = {
      start: moment().add(-29, 'days').format('yyyy-MM-DD'),
      end: moment().format('yyyy-MM-DD'),
      code: state.value.obs.code,
    }
    if (state.value.obs.dataRange === 1) {
      data = {
        start: moment().add(-6, 'days').format('yyyy-MM-DD'),
        end: moment().format('yyyy-MM-DD'),
        code: state.value.obs.code,
      }
    }
    let res = await apiNormalGostObsNodeMonth(data)
    state.value.obs.data = res.data || []
  } finally {
    state.value.obs.loading = false
  }
}

watch(() => ({type: state.value.obs.dataRange}), () => {
  obsFunc()
})

onBeforeMount(() => {
  pageFunc()
})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-alert type="info">
        节点运行命令：
        <div>./gostc -s -key xxxxxx</div>
        <div>将xxxxxx修改为你的节点连接密钥</div>

        <div>其他问题：Linux可能会碰到权限问题，执行以下命令解决：sudo chmod +x gostc</div>
      </n-alert>
    </AppCard>

    <SearchCard :show-border="false" space>
      <SearchItem
          type="input"
          :label-width="70"
          clearable
          label="名称"
          @onChange="value => state.table.search.name=value"
      ></SearchItem>
      <SearchItem type="custom">
        <n-space>
          <n-button type="info" :focusable="false" @click="searchTable">搜索</n-button>
          <n-button type="info" :focusable="false" @click="refreshTable">刷新</n-button>
          <n-button type="info" :focusable="false" @click="openCreate">新增</n-button>
          <n-button
              :focusable="false"
              type="warning"
              @click="goToUrl('https://docs.sian.one/gostc/node')">
            自建教程
          </n-button>
        </n-space>
      </SearchItem>
    </SearchCard>

    <AppCard :show-border="false" :loading="state.table.searchLoading">
      <Empty v-if="state.table.list.length===0" border description="暂无节点"></Empty>
      <n-grid v-else x-gap="12" y-gap="12" cols="1 520:2 900:3 1400:4">
        <n-grid-item v-for="row in state.table.list">
          <n-el class="client-item" tag="div" :style="{
                border: '1px solid var(--border-color)',
                borderRadius:'var(--border-radius)',
                padding: '12px',
                cursor: 'pointer'}">
            <n-h4 style="margin-bottom: 8px !important;">
              <n-space justify="space-between">
                <span style="font-weight: bold">{{ row.name }}</span>
                <span :style="statusStyle(row)">
                {{ statusText(row) }}
              </span>
              </n-space>
            </n-h4>
            <div>
              <span>版本：{{ row.version || '未知' }}</span><br>
              <span>地址：{{ row.address }}</span><br>
              <span>协议：{{ row.protocol }}</span><br>
              <span>密钥：{{ row.key }}</span><br>
              <span>介绍：{{ row.remark }}</span><br>
              <span>标签：{{ row.tags.join('、') }}</span><br>
              <span>规则：{{ row.ruleNames.join('、') }}</span><br>
              <span>流量( IN | OUT )：{{ flowFormat(row.inputBytes) + ' | ' + flowFormat(row.outputBytes) }}</span><br>
              <n-tabs animated size="small">
                <n-tab-pane name="web" tab="域名解析">
                  <span>基础域名：{{ row.domain }}</span><br>
                  <span>不允许的域名前缀：{{ row.denyDomainPrefix || '暂无' }}</span><br>
                </n-tab-pane>
                <n-tab-pane name="forward" tab="端口转发">
                  <span>连接端口：{{ row.forwardConnPort }}</span><br>
                  <span>开发端口：{{ row.forwardPorts }}</span><br>
                </n-tab-pane>
                <n-tab-pane name="tunnel" tab="私有隧道">
                  <span>连接端口：{{ row.tunnelConnPort }}</span><br>
                  <span>访问端口：{{ row.tunnelInPort }}</span><br>
                </n-tab-pane>
              </n-tabs>
            </div>
            <n-space justify="end" style="width: 100%">
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openObsModal(row)">
                流量
              </n-button>
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
    </AppCard>

    <AppCard :show-border="false">
      <n-pagination
          :page-size="state.table.search.size"
          :page="state.table.search.page"
          :item-count="state.table.total"
          :simple="true"
          :on-update-page="(val) => {state.table.search.page = val;refreshTable()}"
          :on-update-page-size="(val) => {state.table.search.size = val;refreshTable()}"
      />
    </AppCard>

    <Modal
        :show="state.create.open"
        title="新增节点"
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
        <n-form-item path="remark" label="介绍">
          <n-input type="textarea" v-model:value.trim="state.create.data.remark"></n-input>
        </n-form-item>
        <n-form-item path="indexValue" label="排序(升序)">
          <n-input-number v-model:value="state.create.data.indexValue"></n-input-number>
        </n-form-item>
        <n-form-item path="address" label="地址(IP|HOST)">
          <n-input v-model:value.trim="state.create.data.address"></n-input>
        </n-form-item>
        <n-form-item path="tunnelProtocol" label="连接协议">
          <n-select
              :default-value="state.create.data.protocol"
              v-model:value="state.create.data.protocol"
              :options="protocols"
              label-field="label"
              value-field="value"
          ></n-select>
        </n-form-item>
        <n-form-item label="功能">
          <n-space>
            <n-checkbox
                v-model:value="state.create.data.web"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.create.data.web = value}"
            >域名解析
            </n-checkbox>
            <n-checkbox
                v-model:value="state.create.data.forward"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.create.data.forward = value}"
            >端口转发
            </n-checkbox>
            <n-checkbox
                v-model:value="state.create.data.tunnel"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.create.data.tunnel = value}"
            >私有隧道
            </n-checkbox>
          </n-space>
        </n-form-item>
        <n-tabs type="line" animated>
          <n-tab-pane name="web" tab="域名解析">
            <div v-show="state.create.data.tunnel!==1 || state.create.data.web !==1">
              <n-alert type="warning" :show-icon="false" v-show="state.create.data.web!==1">
                未启用域名解析
              </n-alert>
              <n-alert type="warning" :show-icon="false"
                       v-show="state.create.data.web===1 && state.create.data.tunnel!==1">
                域名解析，必须启用私有隧道
              </n-alert>
              <br>
            </div>
            <n-form-item path="tunnelInPort" label="访问端口">
              <n-input v-model:value.trim="state.create.data.tunnelInPort" placeholder="3333"></n-input>
            </n-form-item>
            <n-form-item path="domain" label="基础域名">
              <n-input v-model:value.trim="state.create.data.domain" placeholder="example.com"></n-input>
            </n-form-item>
            <n-form-item path="denyDomainPrefix" label="不允许的域名前缀">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.create.data.denyDomainPrefix"
                  placeholder="一行一个"
              ></n-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="forward" tab="端口转发">
            <div v-show="state.create.data.forward!==1">
              <n-alert type="warning" :show-icon="false">
                未启用端口转发
              </n-alert>
              <br>
            </div>
            <n-form-item path="forwardConnPort" label="连接端口">
              <n-input v-model:value="state.create.data.forwardConnPort" placeholder="2097"></n-input>
            </n-form-item>
            <n-form-item path="forwardPorts" label="端口配额">
              <n-input v-model:value="state.create.data.forwardPorts" placeholder="10001-11000,20000,30000"></n-input>
            </n-form-item>
            <n-form-item path="forwardReplaceAddress" label="替换地址(一般留空)">
              <n-input v-model:value="state.create.data.forwardReplaceAddress" placeholder="grpc://1.1.1.1:8080"></n-input>
            </n-form-item>
            <n-form-item path="forwardMetadata" label="METADATA">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.create.data.forwardMetadata"
                  placeholder="JSON格式"
              ></n-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="tunnel" tab="私有隧道">
            <div v-show="state.create.data.tunnel!==1">
              <n-alert type="warning" :show-icon="false">
                未启用私有隧道
              </n-alert>
              <br>
            </div>
            <n-form-item path="tunnelConnPort" label="连接端口">
              <n-input v-model:value="state.create.data.tunnelConnPort" placeholder="2096"></n-input>
            </n-form-item>
            <n-form-item path="tunnelReplaceAddress" label="替换地址(一般留空)">
              <n-input v-model:value="state.create.data.tunnelReplaceAddress" placeholder="grpc://1.1.1.1:8080"></n-input>
            </n-form-item>
            <n-form-item path="tunnelMetadata" label="METADATA">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.create.data.tunnelMetadata"
                  placeholder="JSON格式"
              ></n-input>
            </n-form-item>
          </n-tab-pane>
        </n-tabs>
      </n-form>
    </Modal>

    <Modal
        :show="state.update.open"
        title="修改节点"
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
        <n-form-item path="remark" label="介绍">
          <n-input type="textarea" v-model:value.trim="state.update.data.remark"></n-input>
        </n-form-item>
        <n-form-item path="indexValue" label="排序(升序)">
          <n-input-number v-model:value="state.update.data.indexValue"></n-input-number>
        </n-form-item>
        <n-form-item path="address" label="地址(IP|HOST)">
          <n-input v-model:value.trim="state.update.data.address"></n-input>
        </n-form-item>
        <n-form-item path="tunnelProtocol" label="连接协议">
          <n-select
              :default-value="state.update.data.protocol"
              v-model:value="state.update.data.protocol"
              :options="protocols"
              label-field="label"
              value-field="value"
          ></n-select>
        </n-form-item>
        <n-form-item label="功能">
          <n-space>
            <n-checkbox
                v-model:checked="state.update.data.web"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.update.data.web = value}"
            >域名解析
            </n-checkbox>
            <n-checkbox
                v-model:checked="state.update.data.forward"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.update.data.forward = value}"
            >端口转发
            </n-checkbox>
            <n-checkbox
                v-model:checked="state.update.data.tunnel"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.update.data.tunnel = value}"
            >私有隧道
            </n-checkbox>
          </n-space>
        </n-form-item>
        <n-tabs type="line" animated>
          <n-tab-pane name="web" tab="域名解析">
            <div v-show="state.update.data.tunnel!==1 || state.update.data.web !==1">
              <n-alert type="warning" :show-icon="false" v-show="state.update.data.web!==1">
                未启用域名解析
              </n-alert>
              <n-alert type="warning" :show-icon="false"
                       v-show="state.update.data.web===1 && state.update.data.tunnel!==1">
                域名解析，必须启用私有隧道
              </n-alert>
              <br>
            </div>
            <n-form-item path="tunnelInPort" label="访问端口">
              <n-input v-model:value.trim="state.update.data.tunnelInPort" placeholder="3333"></n-input>
            </n-form-item>
            <n-form-item path="domain" label="基础域名">
              <n-input v-model:value.trim="state.update.data.domain" placeholder="example.com"></n-input>
            </n-form-item>
            <n-form-item path="denyDomainPrefix" label="不允许的域名前缀">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.update.data.denyDomainPrefix"
                  placeholder="一行一个"
              ></n-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="forward" tab="端口转发">
            <div v-show="state.update.data.forward!==1">
              <n-alert type="warning" :show-icon="false">
                未启用端口转发
              </n-alert>
              <br>
            </div>
            <n-form-item path="forwardConnPort" label="连接端口">
              <n-input v-model:value="state.update.data.forwardConnPort" placeholder="2097"></n-input>
            </n-form-item>
            <n-form-item path="forwardPorts" label="端口配额">
              <n-input v-model:value="state.update.data.forwardPorts" placeholder="10001-11000,20000,30000"></n-input>
            </n-form-item>
            <n-form-item path="forwardReplaceAddress" label="替换地址(一般留空)">
              <n-input v-model:value="state.update.data.forwardReplaceAddress" placeholder="grpc://1.1.1.1:8080"></n-input>
            </n-form-item>
            <n-form-item path="forwardMetadata" label="METADATA">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.update.data.forwardMetadata"
                  placeholder="JSON格式"
              ></n-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="tunnel" tab="私有隧道">
            <div v-show="state.update.data.tunnel!==1">
              <n-alert type="warning" :show-icon="false">
                未启用私有隧道
              </n-alert>
              <br>
            </div>
            <n-form-item path="tunnelConnPort" label="连接端口">
              <n-input v-model:value="state.update.data.tunnelConnPort" placeholder="2096"></n-input>
            </n-form-item>
            <n-form-item path="tunnelReplaceAddress" label="替换地址(一般留空)">
              <n-input v-model:value="state.update.data.tunnelReplaceAddress" placeholder="grpc://1.1.1.1:8080"></n-input>
            </n-form-item>
            <n-form-item path="tunnelMetadata" label="METADATA">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.update.data.tunnelMetadata"
                  placeholder="JSON格式"
              ></n-input>
            </n-form-item>
          </n-tab-pane>
        </n-tabs>
      </n-form>
    </Modal>

    <Modal
        title="流量情况"
        :show="state.obs.open"
        confirm-text=""
        cancel-text="关闭"
        @on-cancel="closeObsModal"
        mask-close
    >
      <n-space justify="space-between">
        <n-h4 style="font-weight: bold">最近{{ state.obs.dataRange === 1 ? '7' : '30' }}天流量使用趋势</n-h4>
        <n-radio-group size="small" v-model:value="state.obs.dataRange">
          <n-radio-button :value="1">最近7天</n-radio-button>
          <n-radio-button :value="2">最近30天</n-radio-button>
        </n-radio-group>
      </n-space>
      <Obs :data="state.obs.data" :loading="state.obs.loading"></Obs>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>