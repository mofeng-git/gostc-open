<script setup>
import {onBeforeMount, ref, watch} from "vue";
import {
  apiAdminGostNodeCleanPort,
  apiAdminGostNodeCreate,
  apiAdminGostNodeDelete,
  apiAdminGostNodePage,
  apiAdminGostNodeUpdate
} from "../../../api/admin/gost_node.js";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import Modal from "../../../components/Modal.vue";
import {requiredRule} from "../../../utils/formDataRule.js";
import {protocols} from "./protocols.js"
import {apiAdminGostNodeRuleList} from "../../../api/admin/gost_node_rule.js";
import {apiAdminSystemUserList} from "../../../api/admin/system_user.js";
import {apiNormalGostNodeBindUpdate} from "../../../api/admin/gost_node_bind.js";
import {flowFormat} from "../../../utils/flow.js";
import {NButton, NPopconfirm, NSpace} from "naive-ui";
import Empty from "../../../components/Empty.vue";
import moment from "moment/moment.js";
import {apiNormalGostObsNodeMonth} from "../../../api/normal/gost_obs.js";
import Obs from "../../../components/Obs.vue";
import {localStore} from "../../../store/local.js";

const state = ref({
  table: {
    searchLoading: false,
    search: {
      page: 1,
      size: 12,
      name: '',
      bind: 2,
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
      proxy: 2,
      p2p: 1,
      domain: '',
      denyDomainPrefix: '',
      address: '',
      protocol: protocols[0].value,
      urlTpl: '',
      tunnelConnPort: '',
      tunnelInPort: '',
      tunnelMetadata: '',
      forwardConnPort: '',
      forwardPorts: '',
      forwardMetadata: '',
      tunnelReplaceAddress: '',
      forwardReplaceAddress: '',
      p2pPort: '',
      p2pDisableForward: 0,
      rules: [],
      tags: [],
      indexValue: 1000,
      limitResetIndex: 0,
      limitTotal: 1,
      limitKind: 0,
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
      proxy: 2,
      p2p: 1,
      domain: '',
      denyDomainPrefix: '',
      address: '',
      protocol: protocols[0].value,
      urlTpl: '',
      tunnelConnPort: '',
      tunnelInPort: '',
      tunnelMetadata: '',
      forwardConnPort: '',
      forwardPorts: '',
      forwardMetadata: '',
      tunnelReplaceAddress: '',
      forwardReplaceAddress: '',
      p2pPort: '',
      p2pDisableForward: 0,
      rules: [],
      tags: [],
      indexValue: 1000,
      limitResetIndex: 0,
      limitTotal: 1,
      limitKind: 0,
    },
    dataRules: {
      name: requiredRule('请输入名称'),
      address: requiredRule('请输入服务器地址'),
    },
  },
  rules: [],
  bind: {
    loading: false,
    open: false,
    users: [],
    usersLoading: false,
    data: {
      nodeCode: '',
      userCode: '',
    },
  },
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
    let res = await apiAdminGostNodePage(state.value.table.search)
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
    proxy: 2,
    p2p: 1,
    domain: '',
    denyDomainPrefix: '',
    address: '',
    protocol: protocols[0].value,
    urlTpl: '',
    tunnelConnPort: '',
    tunnelInPort: '',
    tunnelMetadata: '',
    forwardConnPort: '',
    forwardPorts: '',
    forwardMetadata: '',
    tunnelReplaceAddress: '',
    forwardReplaceAddress: '',
    p2pPort: '',
    p2pDisableForward: 0,
    rules: [],
    tags: [],
    indexValue: 1000,
    limitResetIndex: 0,
    limitTotal: 1,
    limitKind: 0,
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
        await apiAdminGostNodeCreate(state.value.create.data)
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
        await apiAdminGostNodeUpdate(state.value.update.data)
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
    await apiAdminGostNodeDelete({code: row.code})
    searchTable()
  } finally {
  }
}

const openBind = (row) => {
  state.value.bind.data.nodeCode = row.code
  state.value.bind.data.userCode = null
  state.value.bind.open = true
}

const closeBind = () => {
  state.value.bind.open = false
}

const bindRef = ref()
const bindFunc = () => {
  bindRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.bind.loading = true
        await apiNormalGostNodeBindUpdate(state.value.bind.data)
        refreshTable()
        closeBind()
      } finally {
        state.value.bind.loading = false
      }
    }
  })
}

const getRules = async () => {
  try {
    let res = await apiAdminGostNodeRuleList()
    state.value.rules = res.data || []
  } finally {

  }
}

const searchUsers = async (account) => {
  if (!account) {
    return
  }
  try {
    state.value.bind.usersLoading = true
    let res = await apiAdminSystemUserList({account: account})
    state.value.bind.users = res.data || []
  } finally {
    state.value.bind.usersLoading = false
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

const nodeCleanPortFunc = async (row) => {
  try {
    await apiAdminGostNodeCleanPort({code: row.code})
    $message.create('清除成功', {
      duration: 1500,
      closable: true,
      type: "success"
    })
  } finally {

  }
}

const nodeObsLimitFunc = (row) => {
  if (row.limitResetIndex === 0 || row.limitUseTotal === -1) {
    return '未配置'
  }
  return flowFormat(row.limitUseTotal) + " | " + row.limitTotal + 'GB'
}

onBeforeMount(() => {
  pageFunc()
  getRules()
})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-alert type="warning">
        不推荐启用节点的代理隧道功能，非正常用途，可能导致节点IP被墙
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
      <SearchItem
          type="select"
          :label-width="120"
          label="节点类型"
          :default="state.table.search.bind"
          :options="[
              {value:'全部',key:0},
              {value:'用户节点',key:1},
              {value:'系统节点',key:2}
          ]"
          @onChange="value => state.table.search.bind=value"
      ></SearchItem>
      <SearchItem type="custom">
        <n-space>
          <n-button type="info" :focusable="false" @click="searchTable">搜索</n-button>
          <n-button type="info" :focusable="false" @click="refreshTable">刷新</n-button>
          <n-button type="info" :focusable="false" @click="openCreate">新增</n-button>
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
              <span>归属用户：{{ row.userAccount || '系统' }}</span><br>
              <span>版本：{{ row.version || '未知' }}</span><br>
              <span>地址：{{ row.address }}</span><br>
              <span>协议：{{ row.protocol }}</span><br>
              <span>密钥：{{ row.key }}</span><br>
              <span>介绍：{{ row.remark }}</span><br>
              <span>标签：{{ row.tags.join('、') }}</span><br>
              <span>规则：{{ row.ruleNames.join('、') }}</span><br>
              <span>流量循环：{{ nodeObsLimitFunc(row) }}</span><br>
              <span>30天流量( IN | OUT )：{{
                  flowFormat(row.inputBytes) + ' | ' + flowFormat(row.outputBytes)
                }}</span><br>
              <span>自定义域名：{{ row.customDomain === 1 ? '支持' : '不支持' }}</span><br>
              <span>代理隧道：{{ row.proxy === 1 ? '启用' : '禁用' }}</span><br>
              <n-tabs animated size="small">
                <n-tab-pane name="web" tab="域名解析">
                  <span>访问端口：{{ row.tunnelInPort }}</span><br>
                  <span>基础域名：{{ row.domain }}</span><br>
                  <span>不允许的域名前缀：{{ row.denyDomainPrefix || '暂无' }}</span><br>
                </n-tab-pane>
                <n-tab-pane name="tunnel" tab="私有隧道">
                  <span>连接端口：{{ row.tunnelConnPort }}</span><br>
                </n-tab-pane>
                <n-tab-pane name="forward" tab="端口转发">
                  <span>连接端口：{{ row.forwardConnPort }}</span><br>
                  <span>开放端口：{{ row.forwardPorts }}</span><br>
                </n-tab-pane>
                <n-tab-pane name="p2p" tab="P2P隧道">
                  <span>连接端口：{{ row.p2pPort }}</span><br>
                  <span>中继转发：{{ row.p2pDisableForward === 0 ? '启用' : '禁用' }}</span><br>
                </n-tab-pane>
              </n-tabs>
            </div>
            <n-space justify="end" style="width: 100%">
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openObsModal(row)">
                流量
              </n-button>
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openBind(row)">
                绑定
              </n-button>
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="nodeCleanPortFunc(row)">
                清除端口占用缓存
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
        <n-form-item path="tags" label="标签">
          <n-select
              v-model:value="state.create.data.tags"
              filterable multiple tag
              placeholder="请输入"
              :show-arrow="false"
              :show="false"
          />
        </n-form-item>
        <n-form-item path="rules" label="规则">
          <n-select
              v-model:value="state.create.data.rules"
              multiple
              label-field="name"
              value-field="code"
              placeholder="请选择"
              :options="state.rules"
          />
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
                v-model:checked="state.create.data.web"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.create.data.web = value}"
            >域名解析
            </n-checkbox>
            <n-checkbox
                v-model:checked="state.create.data.tunnel"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.create.data.tunnel = value}"
            >私有隧道
            </n-checkbox>
            <n-checkbox
                v-model:checked="state.create.data.forward"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.create.data.forward = value}"
            >端口转发
            </n-checkbox>
            <n-checkbox
                v-model:checked="state.create.data.proxy"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.create.data.proxy = value}"
            >代理隧道
            </n-checkbox>
            <n-checkbox
                v-model:checked="state.create.data.p2p"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.create.data.p2p = value}"
            >P2P隧道
            </n-checkbox>
          </n-space>
        </n-form-item>
        <n-alert type="warning" :show-icon="false"
                 v-show="state.create.data.proxy===1 && state.create.data.forward!==1">
          代理隧道，必须启用端口转发
        </n-alert>
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
            <n-form-item path="urlTpl" label="域名模板(展示域名解析访问地址)">
              <n-input v-model:value.trim="state.create.data.urlTpl" placeholder="https://{{DOMAIN}}:443"></n-input>
            </n-form-item>
            <n-form-item path="denyDomainPrefix" label="不允许的域名前缀">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.create.data.denyDomainPrefix"
                  placeholder="一行一个"
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
              <n-input v-model:value="state.create.data.tunnelReplaceAddress"
                       placeholder="grpc://1.1.1.1:8080"></n-input>
            </n-form-item>
            <n-form-item path="tunnelMetadata" label="METADATA">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.create.data.tunnelMetadata"
                  placeholder="JSON格式"
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
              <n-input v-model:value="state.create.data.forwardReplaceAddress"
                       placeholder="grpc://1.1.1.1:8080"></n-input>
            </n-form-item>
            <n-form-item path="forwardMetadata" label="METADATA">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.create.data.forwardMetadata"
                  placeholder="JSON格式"
              ></n-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="p2p" tab="P2P隧道">
            <div v-show="state.create.data.p2p!==1">
              <n-alert type="warning" :show-icon="false">
                未启用P2P隧道
              </n-alert>
              <br>
            </div>
            <n-form-item path="p2pPort" label="连接端口">
              <n-input v-model:value="state.create.data.p2pPort" placeholder="7000"></n-input>
            </n-form-item>
            <n-form-item path="p2pDisableForward" label="中继转发">
              <n-switch
                  :round="false"
                  :default-value="state.create.data.p2pDisableForward"
                  :checked-value="0"
                  :unchecked-value="1"
                  :on-update-value="value => {state.create.data.p2pDisableForward = value}"
              >
                <template #checked>启用转发</template>
                <template #unchecked>禁用转发</template>
              </n-switch>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="limit" tab="流量循环">
            <n-form-item path="limitResetIndex" label="重置日期(0表示流量不循环计算)">
              <n-input-number
                  :min="0" :max="31"
                  v-model:value="state.create.data.limitResetIndex"
              ></n-input-number>
            </n-form-item>
            <n-form-item path="limitTotal" label="预警流量(GB)">
              <n-input-number
                  :min="1"
                  v-model:value="state.create.data.limitTotal"
              ></n-input-number>
            </n-form-item>
            <n-form-item path="limitKind" label="计算方式">
              <n-select
                  v-model:value="state.create.data.limitKind"
                  :default-value="state.create.data.limitKind"
                  :options="[{label:'全部',value:0},{label:'上行',value:1},{label:'下行',value:2}]"
              ></n-select>
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
        <n-form-item path="tags" label="标签">
          <n-select
              v-model:value="state.update.data.tags"
              filterable multiple tag
              placeholder="请输入"
              :show-arrow="false"
              :show="false"
          />
        </n-form-item>
        <n-form-item path="rules" label="规则">
          <n-select
              v-model:value="state.update.data.rules"
              multiple
              label-field="name"
              value-field="code"
              placeholder="请选择"
              :options="state.rules"
          />
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
                v-model:checked="state.update.data.tunnel"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.update.data.tunnel = value}"
            >私有隧道
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
                v-model:checked="state.update.data.proxy"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.update.data.proxy = value}"
            >代理隧道
            </n-checkbox>
            <n-checkbox
                v-model:checked="state.update.data.p2p"
                :focusable="false"
                :checked-value="1"
                :unchecked-value="2"
                :on-update-checked="value=>{state.update.data.p2p = value}"
            >P2P隧道
            </n-checkbox>
          </n-space>
        </n-form-item>
        <n-alert type="warning" :show-icon="false"
                 v-show="state.update.data.proxy===1 && state.update.data.forward!==1">
          代理隧道，必须启用端口转发
        </n-alert>
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
            <n-form-item path="urlTpl" label="域名模板(展示域名解析访问地址)">
              <n-input v-model:value.trim="state.update.data.urlTpl" placeholder="https://{{DOMAIN}}:443"></n-input>
            </n-form-item>
            <n-form-item path="denyDomainPrefix" label="不允许的域名前缀">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.update.data.denyDomainPrefix"
                  placeholder="一行一个"
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
              <n-input v-model:value="state.update.data.tunnelReplaceAddress"
                       placeholder="grpc://1.1.1.1:8080"></n-input>
            </n-form-item>
            <n-form-item path="tunnelMetadata" label="METADATA">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.update.data.tunnelMetadata"
                  placeholder="JSON格式"
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
              <n-input v-model:value="state.update.data.forwardReplaceAddress"
                       placeholder="grpc://1.1.1.1:8080"></n-input>
            </n-form-item>
            <n-form-item path="forwardMetadata" label="METADATA">
              <n-input
                  type="textarea"
                  v-model:value.trim="state.update.data.forwardMetadata"
                  placeholder="JSON格式"
              ></n-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="p2p" tab="P2P隧道">
            <div v-show="state.update.data.p2p!==1">
              <n-alert type="warning" :show-icon="false">
                未启用P2P隧道
              </n-alert>
              <br>
            </div>
            <n-form-item path="p2pPort" label="连接端口">
              <n-input v-model:value="state.update.data.p2pPort" placeholder="7000"></n-input>
            </n-form-item>
            <n-form-item path="p2pDisableForward" label="中继转发">
              <n-switch
                  :round="false"
                  :default-value="state.update.data.p2pDisableForward"
                  :checked-value="0"
                  :unchecked-value="1"
                  :on-update-value="value => {state.update.data.p2pDisableForward = value}"
              >
                <template #checked>启用转发</template>
                <template #unchecked>禁用转发</template>
              </n-switch>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="limit" tab="流量循环">
            <n-form-item path="limitResetIndex" label="重置日期(0表示流量不循环计算)">
              <n-input-number
                  :min="0" :max="31"
                  v-model:value="state.update.data.limitResetIndex"
              ></n-input-number>
            </n-form-item>
            <n-form-item path="limitTotal" label="预警流量(GB)">
              <n-input-number
                  :min="1"
                  v-model:value="state.update.data.limitTotal"
              ></n-input-number>
            </n-form-item>
            <n-form-item path="limitKind" label="计算方式">
              <n-select
                  v-model:value="state.update.data.limitKind"
                  :default-value="state.update.data.limitKind"
                  :options="[{label:'全部',value:0},{label:'上行',value:1},{label:'下行',value:2}]"
              ></n-select>
            </n-form-item>
          </n-tab-pane>
        </n-tabs>
      </n-form>
    </Modal>

    <Modal
        :show="state.bind.open"
        title="绑定用户"
        width="400px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.bind.loading"
        @on-confirm="bindFunc"
        @on-cancel="closeBind"
        mask-close
    >
      <n-form ref="bindRef" :model="state.bind.data">
        <n-form-item path="userCode" label="用户">
          <n-select
              v-model:value="state.bind.data.userCode"
              filterable
              placeholder="搜索用户"
              :options="state.bind.users"
              :loading="state.bind.usersLoading"
              label-field="account"
              value-field="code"
              remote
              clearable
              @search="searchUsers"
          />
        </n-form-item>
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
      <Obs :data="state.obs.data" style="width:100%" :loading="state.obs.loading" :dark="localStore().darkTheme"></Obs>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>