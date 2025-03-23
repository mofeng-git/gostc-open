<script setup>
import {onBeforeMount, ref, watch} from "vue";
import {
  apiNormalGostClientForwardAdmission,
  apiNormalGostClientForwardDelete,
  apiNormalGostClientForwardEnable,
  apiNormalGostClientForwardMatcher,
  apiNormalGostClientForwardPage,
  apiNormalGostClientForwardRenew,
  apiNormalGostClientForwardUpdate
} from "../../../api/normal/gost_client_forward.js";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import router from "../../../router/index.js";
import Modal from "../../../components/Modal.vue";
import {regexpRule, requiredRule} from "../../../utils/formDataRule.js";
import {regexpLocalIp, regexpPort} from "../../../utils/regexp.js";
import {NButton, NPopconfirm, NSpace} from "naive-ui";
import Empty from "../../../components/Empty.vue";
import {cLimiterText, configExpText, configText, limiterText, rLimiterText} from "./index.js";
import Alert from "../../../icon/alert.vue";
import Online from "../../../icon/online.vue";
import {flowFormat} from "../../../utils/flow.js";
import {goToUrl} from "../../../utils/browser.js";
import {apiNormalGostObsTunnelMonth} from "../../../api/normal/gost_obs.js";
import moment from "moment";
import Obs from "../../../components/Obs.vue";
import {localStore} from "../../../store/local.js";

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
    minWidth: 1100,
  },
  update: {
    data: {
      name: '',
      targetIp: '',
      targetPort: '',
      proxyProtocol: 0,
    },
    dataRules: {
      name: requiredRule('请输入名称'),
      targetIp: regexpRule(regexpLocalIp, '内网IP格式错误'),
      targetPort: regexpRule(regexpPort, '内网端口格式错误'),
    },
    open: false,
    loading: false,
  },
  matcher: {
    open: false,
    loading: false,
    data: {
      code: '',
      enable: 2,
      matchers: [],
      tcpMatcher: {
        targetIp: '',
        targetPort: '',
      },
      sshMatcher: {
        targetIp: '',
        targetPort: '',
      },
    }
  },
  admission: {
    open: false,
    loading: false,
    data: {
      whiteEnable: 2,
      blackEnable: 2,
      white: [],
      black: [],
    },
    white: '',
    black: '',
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
    let res = await apiNormalGostClientForwardPage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}

const openCreate = () => {
  router.push({name: 'NormalGostClientForwardCreate'})
}

const openUpdate = (row) => {
  state.value.update.data = JSON.parse(JSON.stringify(row))
  state.value.update.open = true
}

const closeUpdate = () => {
  state.value.update.open = false
}

const updateRef = ref()
const updateFunc = () => {
  updateRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.update.loading = true
        await apiNormalGostClientForwardUpdate(state.value.update.data)
        refreshTable()
        closeUpdate()
      } finally {
        state.value.update.loading = false
      }
    }
  })
}

const openMatcher = (row) => {
  state.value.matcher.data = {
    code: row.code,
    enable: row?.matcherEnable,
    matchers: JSON.parse(JSON.stringify(row.matchers || [])),
    tcpMatcher: JSON.parse(JSON.stringify(row.tcpMatcher)),
    sshMatcher: JSON.parse(JSON.stringify(row.sshMatcher)),
  }
  state.value.matcher.open = true
}

const closeMatcher = () => {
  state.value.matcher.open = false
}

const matcherFunc = async () => {
  try {
    state.value.matcher.loading = true
    await apiNormalGostClientForwardMatcher(state.value.matcher.data)
    closeMatcher()
    refreshTable()
  } finally {
    state.value.matcher.loading = false
  }
}

const enableFunc = async (enable, row) => {
  try {
    row.enableLoading = true
    await apiNormalGostClientForwardEnable({code: row.code, enable: enable})
    refreshTable()
  } finally {
    row.enableLoading = false
  }
}

const renewFunc = async (row) => {
  try {
    row.renewLoading = true
    await apiNormalGostClientForwardRenew({code: row.code})
    refreshTable()
  } finally {
    row.renewLoading = false
  }
}

const deleteFunc = async (row) => {
  try {
    await apiNormalGostClientForwardDelete({code: row.code})
    searchTable()
  } finally {
  }
}

const openAdmission = (row) => {
  state.value.admission.data = {
    code: row.code,
    whiteEnable: row.whiteEnable,
    blackEnable: row.blackEnable,
    white: [],
    black: []
  }
  if (row.whiteList) {
    state.value.admission.white = row?.whiteList.join('\n')
  }
  if (row.blackList) {
    state.value.admission.black = row?.blackList.join('\n')
  }
  state.value.admission.open = true
}

const closeAdmission = () => {
  state.value.admission.open = false
}

const admissionFunc = async () => {
  try {
    state.value.admission.loading = true
    state.value.admission.data.white = state.value.admission.white.split('\n')
    state.value.admission.data.black = state.value.admission.black.split('\n')
    await apiNormalGostClientForwardAdmission(state.value.admission.data)
    closeAdmission()
    refreshTable()
  } finally {
    state.value.admission.loading = false
  }
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
    let res = await apiNormalGostObsTunnelMonth(data)
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
        请注意，删除隧道不会退还积分
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
          :label-width="70"
          label="状态"
          :default="0"
          :options="[
              {value:'全部',key:0},
              {value:'启用',key:1},
              {value:'停用',key:2}
          ]"
          @onChange="value => state.table.search.enable=value"
      ></SearchItem>
      <SearchItem type="custom">
        <n-space>
          <n-button type="info" :focusable="false" @click="searchTable">搜索</n-button>
          <n-button type="info" :focusable="false" @click="refreshTable">刷新</n-button>
          <n-button type="info" :focusable="false" @click="openCreate">新增</n-button>
          <n-button
              :focusable="false"
              type="warning"
              @click="goToUrl('https://docs.sian.one/gostc/forward')">
            使用教程
          </n-button>
        </n-space>
      </SearchItem>
    </SearchCard>

    <AppCard :show-border="false" :loading="state.table.searchLoading">
      <Empty v-if="state.table.list.length===0" border description="暂无数据"></Empty>
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
                <div style="display: flex;justify-content: center">
                  <n-tooltip v-if="row?.warnMsg">
                    <template #trigger>
                      <Alert :size="20"></Alert>
                    </template>
                    {{ row.warnMsg }}
                  </n-tooltip>
                  <n-switch
                      style="margin-left: 8px"
                      size="small"
                      :loading="row.enableLoading"
                      v-model:value="row.enable"
                      :checked-value="1"
                      :unchecked-value="2"
                      :on-update:value="value => {enableFunc(value,row)}"
                  ></n-switch>
                </div>
              </n-space>
            </n-h4>
            <div>
              <span>节点：
                <Online :size="10" :online="row.node.online===1"></Online>
                &nbsp&nbsp{{ row.node.name }}</span><br>
              <span>客户端：<Online :size="10" :online="row.client.online===1"></Online>&nbsp&nbsp{{
                  row.client.name
                }}</span><br>
              <span>内网目标：{{ row.targetIp + ':' + row.targetPort }}</span><br>
              <span>访问地址：{{ row.node.address + ':' + row.port }}</span><br>
              <span>速率：{{ limiterText(row.config.limiter) }}</span><br>
              <span>并发数：{{ rLimiterText(row.config.rLimiter) }}</span><br>
              <span>连接数：{{ cLimiterText(row.config.cLimiter) }}</span><br>
              <span>套餐：{{ configText(row.config) }}</span><br>
              <span>到期时间：{{ configExpText(row.config) }}</span><br>
              <span>流量( IN | OUT )：{{ flowFormat(row.inputBytes) + ' | ' + flowFormat(row.outputBytes) }}</span><br>
            </div>
            <n-space justify="end" style="width: 100%">
              <n-popconfirm
                  v-if="row.config.chargingType===2"
                  @positive-click="renewFunc(row)"
                  :positive-button-props="{loading:row.renewLoading}"
              >
                <template #trigger>
                  <n-button
                      size="tiny"
                      :focusable="false"
                      quaternary
                      type="info"
                  >续费
                  </n-button>
                </template>
                确认续费吗？
              </n-popconfirm>
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openObsModal(row)">
                流量
              </n-button>
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openAdmission(row)">
                白/黑名单
              </n-button>
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openMatcher(row)">
                规则匹配
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
        title="修改"
        :show="state.update.open"
        @on-confirm="updateFunc"
        @on-cancel="closeUpdate"
        :confirm-loading="state.update.loading"
    >
      <n-form
          ref="updateRef"
          :model="state.update.data"
          :rules="state.update.dataRules"
          :show-label="true"
          size="medium"
      >
        <n-form-item path="name" label="名称">
          <n-input v-model:value="state.update.data.name" placeholder="我的服务"></n-input>
        </n-form-item>

        <n-alert type="info" :bordered="false" v-if="state.update.data?.matcherEnable===1">已开启规则匹配</n-alert>

        <n-form-item path="targetIp" label="内网IP">
          <n-input
              v-model:value="state.update.data.targetIp"
              :disabled="state.update.data?.matcherEnable===1"
              placeholder="127.0.0.1"
          ></n-input>
        </n-form-item>
        <n-form-item path="targetPort" label="内网端口">
          <n-input
              v-model:value="state.update.data.targetPort"
              :disabled="state.update.data?.matcherEnable===1"
              placeholder="80"
          ></n-input>
        </n-form-item>
        <n-form-item path="proxyProtocol" label="PROXY Protocol">
          <n-select
              :options="[{label:'不启用',value:0},{label:'V1',value:1},{label:'V2',value:2}]"
              v-model:value="state.update.data.proxyProtocol"
          ></n-select>
        </n-form-item>
      </n-form>
    </Modal>

    <Modal
        title="规则匹配"
        :show="state.matcher.open"
        @on-confirm="matcherFunc"
        @on-cancel="closeMatcher"
        :confirm-loading="state.matcher.loading"
    >
      <n-alert type="info" :bordered="false">开启规则匹配，普通配置的内网设置将失效</n-alert>
      <br>
      <n-switch v-model:value="state.matcher.data.enable" :checked-value="1" :unchecked-value="2" :round="false">
        <template #checked>开启规则匹配</template>
        <template #unchecked>关闭规则匹配</template>
      </n-switch>
      <n-divider></n-divider>
      <n-tabs>
        <n-tab-pane name="http" tab="HTTP协议">
          <div v-for="(matcher,index) in state.matcher.data.matchers" :key="index">
            <n-h5 style="font-weight: bold">
              <n-space justify="space-between">
                <span>规则{{ index + 1 }}</span>
                <n-button text type="error" @click="state.matcher.data.matchers.splice(index,1)">删除</n-button>
              </n-space>
            </n-h5>
            <n-space vertical>
              <n-input-group>
                <n-input-group-label style="width: 50%;max-width: 100px;">内网目标</n-input-group-label>
                <n-input
                    placeholder="127.0.0.1"
                    v-model:value="state.matcher.data.matchers[index].targetIp"
                ></n-input>
                <n-input
                    placeholder="8080"
                    v-model:value="state.matcher.data.matchers[index].targetPort"
                ></n-input>
              </n-input-group>
              <n-input-group>
                <n-input-group-label style="width: 50%;max-width: 100px;">匹配域名</n-input-group-label>
                <n-input
                    placeholder="www.example.com"
                    v-model:value="state.matcher.data.matchers[index].host"
                ></n-input>
              </n-input-group>
            </n-space>
            <br>
          </div>
          <n-button
              type="info"
              size="small"
              @click="()=>state.matcher.data.matchers.push({host:'',targetIp:'',targetPort:''})"
              :focusable="false"
          >
            添加
          </n-button>
        </n-tab-pane>
        <n-tab-pane name="tcp" tab="TCP协议">
          <n-space vertical>
            <n-input-group>
              <n-input-group-label style="width: 50%;max-width: 100px;">内网目标</n-input-group-label>
              <n-input
                  placeholder="127.0.0.1"
                  v-model:value="state.matcher.data.tcpMatcher.targetIp"
              ></n-input>
              <n-input
                  placeholder="8080"
                  v-model:value="state.matcher.data.tcpMatcher.targetPort"
              ></n-input>
            </n-input-group>
          </n-space>
        </n-tab-pane>
        <n-tab-pane name="ssh" tab="SSH协议">
          <n-space vertical>
            <n-input-group>
              <n-input-group-label style="width: 50%;max-width: 100px;">内网目标</n-input-group-label>
              <n-input
                  placeholder="127.0.0.1"
                  v-model:value="state.matcher.data.sshMatcher.targetIp"
              ></n-input>
              <n-input
                  placeholder="8080"
                  v-model:value="state.matcher.data.sshMatcher.targetPort"
              ></n-input>
            </n-input-group>
          </n-space>
        </n-tab-pane>
      </n-tabs>
    </Modal>

    <Modal
        title="白/黑名单"
        :show="state.admission.open"
        @on-confirm="admissionFunc"
        @on-cancel="closeAdmission"
        :confirm-loading="state.admission.loading"
    >
      <n-tabs>
        <n-tab-pane name="white" tab="白名单">
          <n-alert type="info">
            白名单：只允许配置的地址访问服务
          </n-alert>
          <p></p>
          <n-switch v-model:value="state.admission.data.whiteEnable" :checked-value="1" :unchecked-value="2"
                    :round="false">
            <template #checked>开启白名单</template>
            <template #unchecked>关闭白名单</template>
          </n-switch>
          <p></p>
          <n-input type="textarea" :autosize="{minRows:5,maxRows:20}" v-model:value="state.admission.white"
                   :placeholder="`127.0.0.1\n192.168.0.0/16`"></n-input>
        </n-tab-pane>
        <n-tab-pane name="black" tab="黑名单">
          <n-alert type="info">
            黑名单：不允许配置的地址访问服务
          </n-alert>
          <p></p>
          <n-switch v-model:value="state.admission.data.blackEnable" :checked-value="1" :unchecked-value="2"
                    :round="false">
            <template #checked>开启黑名单</template>
            <template #unchecked>关闭黑名单</template>
          </n-switch>
          <p></p>
          <n-input type="textarea" :autosize="{minRows:5,maxRows:20}" v-model:value="state.admission.black"
                   :placeholder="`127.0.0.1\n192.168.0.0/16`"></n-input>
        </n-tab-pane>
      </n-tabs>
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