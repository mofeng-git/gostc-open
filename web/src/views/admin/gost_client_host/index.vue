<script setup>
import {onBeforeMount, ref, watch} from "vue";
import {
  apiAdminGostClientHostConfig,
  apiAdminGostClientHostDelete,
  apiAdminGostClientHostPage
} from "../../../api/admin/gost_client_host.js";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import Empty from "../../../components/Empty.vue";
import Alert from "../../../icon/alert.vue";
import Online from "../../../icon/online.vue";
import {cLimiterText, configExpText, configText, limiterText, rLimiterText} from "./index.js";
import moment from "moment/moment.js";
import {NButton, NSpace} from "naive-ui";
import Modal from "../../../components/Modal.vue";
import {flowFormat} from "../../../utils/flow.js";
import {apiNormalGostObsTunnelMonth} from "../../../api/normal/gost_obs.js";
import Obs from "../../../components/Obs.vue";

const state = ref({
  table: {
    searchLoading: false,
    search: {
      page: 1,
      size: 12,
      account: '',
    },
    list: [],
    total: 0,
  },
  config: {
    data: {
      code: '',
      chargingType: 1,
      cycle: 10,
      amount: '0',
      limiter: 1,
      rLimiter: 100,
      cLimiter: 100,
      expAt: '',
    },
    expAt: 0,
    dataRules: {},
    loading: false,
    open: false,
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
    let res = await apiAdminGostClientHostPage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}


const configRef = ref()
const configFunc = () => {
  configRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.config.loading = true
        state.value.config.data.expAt = moment(state.value.config.expAt).format('yyyy-MM-DD HH:mm:ss')
        await apiAdminGostClientHostConfig(state.value.config.data)
        closeConfig()
        await pageFunc()
      } finally {
        state.value.config.loading = false
      }
    }
  })
}

const openConfig = (row) => {
  state.value.config.data = {
    code: row.code,
    chargingType: row.config.chargingType,
    cycle: row.config.cycle,
    amount: row.config.amount,
    limiter: row.config.limiter,
    rLimiter: row.config.rLimiter,
    cLimiter: row.config.cLimiter,
  }
  state.value.config.expAt = moment(row.config.expAt).unix() * 1000
  state.value.config.open = true
}

const closeConfig = () => {
  state.value.config.open = false
}

const deleteFunc = async (row) => {
  try {
    await apiAdminGostClientHostDelete({code: row.code})
    searchTable()
  } finally {
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
    <SearchCard :show-border="false" space>
      <SearchItem
          type="input"
          :label-width="70"
          clearable
          label="账号"
          @onChange="value => state.table.search.account=value"
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
                      disabled
                  ></n-switch>
                </div>
              </n-space>
            </n-h4>
            <div>
              <span>所属用户：{{ row.userAccount }}</span><br>
              <span>节点：
                <Online :size="10" :online="row.node.online===1"></Online>
                &nbsp&nbsp{{ row.node.name }}</span><br>
              <span>客户端：<Online :size="10" :online="row.client.online===1"></Online>&nbsp&nbsp{{
                  row.client.name
                }}</span><br>
              <span>内网目标：{{ row.targetIp + ':' + row.targetPort }}</span><br>
              <span>访问地址：{{ row.domainPrefix + '.' + row.node.domain }}</span><br>
              <span>速率：{{ limiterText(row.config.limiter) }}</span><br>
              <span>并发数：{{ rLimiterText(row.config.rLimiter) }}</span><br>
              <span>连接数：{{ cLimiterText(row.config.cLimiter) }}</span><br>
              <span>套餐：{{ configText(row.config) }}</span><br>
              <span>到期时间：{{ configExpText(row.config) }}</span><br>
              <span>流量( IN | OUT )：{{ flowFormat(row.inputBytes) + ' | ' + flowFormat(row.outputBytes) }}</span><br>
            </div>
            <n-space justify="end" style="width: 100%">
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openObsModal(row)">
                流量
              </n-button>
              <n-button size="tiny" :focusable="false" quaternary type="success" @click="openConfig(row)">
                套餐
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
        :show="state.config.open"
        title="修改套餐"
        width="600px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.config.loading"
        @on-confirm="configFunc"
        @on-cancel="closeConfig"
        mask-close
    >
      <n-form ref="configRef" :rules="state.config.dataRules" :model="state.config.data">
        <n-form-item path="chargingType" label="计费方式">
          <n-radio-group v-model:value="state.config.data.chargingType">
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
        <n-form-item path="cycle" label="续费周期(天)" v-show="state.config.data.chargingType===2">
          <n-input-number v-model:value="state.config.data.cycle" :min="1"></n-input-number>
        </n-form-item>
        <n-form-item path="amount" label="积分" v-show="state.config.data.chargingType!==3">
          <n-input v-model:value.trim="state.config.data.amount"></n-input>
        </n-form-item>
        <n-form-item path="limiter" label="速率(mbps)">
          <n-input-number v-model:value="state.config.data.limiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="rLimiter" label="并发数">
          <n-input-number v-model:value="state.config.data.rLimiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="cLimiter" label="连接数">
          <n-input-number v-model:value="state.config.data.cLimiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="expAt" label="到期时间" v-if="state.config.data.chargingType===2">
          <n-date-picker v-model:value="state.config.expAt" :default-value="state.config.expAt" type="datetime"/>
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
      <Obs :data="state.obs.data" :loading="state.obs.loading"></Obs>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>