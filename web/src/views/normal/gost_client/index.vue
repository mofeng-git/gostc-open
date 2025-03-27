<script setup>

import {onBeforeMount, ref, watch} from "vue";
import {NButton, NPopconfirm, NSpace} from 'naive-ui'
import AppCard from "../../../layout/components/AppCard.vue";
import Modal from "../../../components/Modal.vue";
import Empty from "../../../components/Empty.vue";
import {requiredRule} from "../../../utils/formDataRule.js";
import {
  apiNormalGostClientCreate,
  apiNormalGostClientDelete,
  apiNormalGostClientPage,
  apiNormalGostClientUpdate
} from "../../../api/normal/gost_client.js";
import {copyToClipboard} from "../../../utils/copy.js";
import router from "../../../router/index.js";
import {appStore} from "../../../store/app.js";
import {flowFormat} from "../../../utils/flow.js";
import {goToUrl} from "../../../utils/browser.js";
import moment from "moment/moment.js";
import {apiNormalGostObsClientMonth} from "../../../api/normal/gost_obs.js";
import Obs from "../../../components/Obs.vue";
import {localStore} from "../../../store/local.js";

const state = ref({
  table: {
    searchLoading: false,
    search: {
      page: 1,
      size: 12,
      account: '',
      admin: 0,
    },
    list: [],
    total: 0,
    minWidth: 1100,
  },
  create: {
    data: {
      name: '',
    },
    dataRules: {
      name: requiredRule('名称不能为空'),
    },
    loading: false,
    open: false,
    openLoading: false,
  },
  update: {
    data: {
      name: '',
      code: '',
    },
    dataRules: {
      name: requiredRule('名称不能为空'),
    },
    loading: false,
    open: false,
  },
  look: {
    open: false,
    key: '',
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
    let res = await apiNormalGostClientPage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}

const openCreate = () => {
  state.value.create.data = {
    userCode: '',
    name: '',
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
        await apiNormalGostClientCreate(state.value.create.data)
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
        await apiNormalGostClientUpdate(state.value.update.data)
        refreshTable()
        closeUpdate()
      } finally {
        state.value.update.loading = false
      }
    }
  })
}


const deleteFunc = async (row) => {
  try {
    await apiNormalGostClientDelete({code: row.code})
    searchTable()
  } finally {
  }
}


const openLookFunc = (key) => {
  state.value.look.open = true
  state.value.look.key = key
}

const copyFunc = () => {
  copyToClipboard(state.value.look.key).then(() => {
    $message.create('已复制到剪切板', {
      type: "success",
      closable: true,
      duration: 1500,
    })
  }).catch(err => {
    $message.create('复制失败' + err, {
      type: "error",
      closable: true,
      duration: 1500,
    })
  })
  state.value.look.open = false
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
    let res = await apiNormalGostObsClientMonth(data)
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

const generateCmdString = () => {
  let tls = ' --tls=false'
  if (window.location.protocol.indexOf('https') > -1) {
    tls = ''
  }
  return './gostc' + tls + ' -addr ' + window.location.host + ' -key xxxxxx'
}
</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-alert type="info">
        客户端运行命令：
        <div>{{ generateCmdString() }}</div>
        <div>将xxxxxx修改为你的客户端密钥</div>

        <div>其他问题：Linux可能会碰到权限问题，执行以下命令解决：sudo chmod +x gostc</div>
      </n-alert>
    </AppCard>

    <AppCard :show-border="false">
      <n-alert type="info">
        客户端最新版本：{{ appStore().siteConfig.version }}
      </n-alert>
    </AppCard>

    <AppCard :show-border="false">
      <n-space>
        <n-button type="info" :focusable="false" @click="refreshTable">
          刷新
        </n-button>
        <n-button type="info" :focusable="false" @click="searchTable">
          搜索
        </n-button>
        <n-button type="info" :focusable="false" @click="openCreate">
          新增
        </n-button>
        <n-button
            :focusable="false"
            type="warning"
            @click="goToUrl('https://docs.sian.one/gostc/client')">
          使用教程
        </n-button>
      </n-space>
    </AppCard>

    <AppCard :show-border="false" :loading="state.table.searchLoading">
      <Empty v-if="state.table.list.length===0" border description="暂无客户端"></Empty>
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
              <span>上线IP：{{ row.ip || '未知' }}</span><br>
              <span>上线时间：{{ row.lastTime || '未知' }}</span><br>
              <span>流量( IN | OUT )：{{ flowFormat(row.inputBytes) + ' | ' + flowFormat(row.outputBytes) }}</span><br>
            </div>
            <n-space justify="end" style="width: 100%">
              <n-button size="tiny"
                        :focusable="false" quaternary
                        type="info"
                        @click="router.push({name:'NormalGostClientLogger',query:{clientCode:row.code}})"
              >
                日志
              </n-button>
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openObsModal(row)">
                流量
              </n-button>
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openLookFunc(row.key)">
                连接密钥
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


    <Modal title="客户端密钥" :show="state.look.open"
           @on-confirm="copyFunc"
           confirm-text="复制"
           cancel-text="关闭"
           @on-cancel="state.look.open = false"
           :auto-focus="false"
    >
      <n-p>{{ state.look.key }}</n-p>
    </Modal>

    <Modal title="新增客户端" :show="state.create.open"
           @on-confirm="createFunc"
           :confirm-loading="state.create.loading"
           @on-cancel="closeCreate"
           width="400px"
    >
      <n-form
          ref="createRef"
          :model="state.create.data"
          :rules="state.create.dataRules"
          :show-label="true"
          size="medium"
      >
        <n-form-item path="name" label="名称">
          <n-input
              v-model:value="state.create.data.name"
              placeholder="我的工作电脑"
          />
        </n-form-item>
      </n-form>
    </Modal>

    <Modal title="编辑客户端" :show="state.update.open"
           @on-confirm="updateFunc"
           :confirm-loading="state.update.loading"
           @on-cancel="closeUpdate"
           width="400px"
    >
      <n-form
          ref="updateRef"
          :model="state.update.data"
          :rules="state.update.dataRules"
          :show-label="true"
          size="medium"
      >
        <n-form-item path="name" label="名称">
          <n-input
              v-model:value="state.update.data.name"
              placeholder="我的工作电脑"
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

<style lang="scss">
.client-item {
  //transition: transform 0.2s ease-in-out;
  //transform: scale(1);
  //&:hover {
  //  transform: scale(1.01);
  //}
}
</style>