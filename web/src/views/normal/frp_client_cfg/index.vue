<script setup>
import {h, onBeforeMount, ref} from "vue";
import {
  apiNormalFrpClientCfgDelete,
  apiNormalFrpClientCfgEnable,
  apiNormalFrpClientCfgMigrate,
  apiNormalFrpClientCfgPage,
  apiNormalFrpClientCfgUpdate
} from "../../../api/normal/frp_client_cfg.js";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import router from "../../../router/index.js";
import Modal from "../../../components/Modal.vue";
import {requiredRule} from "../../../utils/formDataRule.js";
import {NButton, NPopconfirm, NSpace} from "naive-ui";
import Empty from "../../../components/Empty.vue";
import Alert from "../../../icon/alert.vue";
import Online from "../../../icon/online.vue";
import {apiNormalGostClientList} from "../../../api/normal/gost_client.js";
import {appStore} from "../../../store/app.js";
import {goToUrl} from "../../../utils/browser.js";

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
      clientCode: '',
      name: '',
      type: 'frpc',
      address: '',
      platform: '',
    },
    dataRules: {
      name: requiredRule('请输入名称'),
      type: requiredRule('请选择配置文件类型'),
      clientCode: requiredRule('请选择客户端'),
    },
    open: false,
    loading: false,
  },
  clients: [],
  migrate: {
    open: false,
    data: {
      code: '',
      clientCode: '',
    },
    loading: false,
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
    let res = await apiNormalFrpClientCfgPage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}

const openCreate = () => {
  router.push({name: 'NormalFrpClientCfgCreate'})
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
        await apiNormalFrpClientCfgUpdate(state.value.update.data)
        refreshTable()
        closeUpdate()
      } finally {
        state.value.update.loading = false
      }
    }
  })
}

const enableFunc = async (enable, row) => {
  try {
    row.enableLoading = true
    await apiNormalFrpClientCfgEnable({code: row.code, enable: enable})
    refreshTable()
  } finally {
    row.enableLoading = false
  }
}

const deleteFunc = async (row) => {
  try {
    await apiNormalFrpClientCfgDelete({code: row.code})
    searchTable()
  } finally {
  }
}

const clientListFunc = async () => {
  try {
    let res = await apiNormalGostClientList()
    state.value.clients = res.data || []
  } finally {

  }
}

const openMigrateModal = (row) => {
  state.value.migrate.data.code = row.code
  state.value.migrate.data.clientCode = row.client.code
  state.value.migrate.open = true
}

const closeMigrateModal = () => {
  state.value.migrate.open = false
}

const migrateFunc = async () => {
  try {
    state.value.migrate.loading = true
    await apiNormalFrpClientCfgMigrate(state.value.migrate.data)
    closeMigrateModal()
    refreshTable()
  } finally {
    state.value.migrate.loading = false
  }
}

onBeforeMount(() => {
  pageFunc()
  clientListFunc()
})

const operatorOptions = [
  {
    label: '转移隧道',
    key: 'migrate',
    disabled: false,
    func: openMigrateModal,
  },
]
const operatorSelect = (key, row) => {
  for (let i = 0; i < operatorOptions.length; i++) {
    if (operatorOptions[i].key === key) {
      operatorOptions[i].func(row)
      return
    }
  }
}

const operatorRenderLabel = (option) => {
  return h(NButton, {
    text: true,
    size: "tiny",
    focusable: false,
    type: "info",
  }, {
    default: () => option.label,
  })
}

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-alert type="info">
        可使用官方的FRP配置下发到客户端，兼容INI和TOML，第三方的魔改的FRP自测兼容性
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
          label="客户端"
          :default="null"
          :options="state.clients"
          :optionsKeyValue="{key: 'code',value: 'name'}"
          @onChange="value => state.table.search.clientCode=value"
          clearable
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
              v-if="appStore().siteConfig.guideConfig.cfgURL"
              :focusable="false"
              type="warning"
              @click="goToUrl(appStore().siteConfig.guideConfig.cfgURL)">
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
              <span>客户端：<Online :size="10" :online="row.client.online===1"></Online>&nbsp&nbsp{{
                  row.client.name
                }}</span><br>
              <span>FRP平台：{{ row.platform }}</span><br>
              <span>访问地址：{{ row.address }}</span><br>
              <span>配置类型：{{ row.type }}</span><br>
            </div>
            <n-space justify="end" style="width: 100%">
              <n-dropdown trigger="hover" size="small" :options="operatorOptions"
                          @select="value => operatorSelect(value,row)" :render-label="operatorRenderLabel">
                <n-button size="tiny" :focusable="false" quaternary type="info">更多操作</n-button>
              </n-dropdown>

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
        <n-form-item path="platform" label="FRP平台(仅展示)">
          <n-input v-model:value="state.update.data.platform" placeholder="FRP平台(仅展示)"></n-input>
        </n-form-item>
        <n-form-item path="address" label="访问地址(仅展示)">
          <n-input v-model:value="state.update.data.address" placeholder="访问地址(仅展示)"></n-input>
        </n-form-item>
        <n-form-item path="type" label="配置文件类型">
          <n-select
              v-model:value="state.update.data.type"
              :options="[{label:'FRPC',value:'frpc'},{label:'FRPS',value:'frps'}]"
          ></n-select>
        </n-form-item>
        <n-form-item path="content" label="配置内容(兼容INI/TOML)">
          <n-input type="textarea" :autosize="{minRows:10,maxRows:20}"
                   v-model:value="state.update.data.content"></n-input>
        </n-form-item>
      </n-form>
    </Modal>

    <Modal
        title="转移隧道"
        :show="state.migrate.open"
        @on-confirm="migrateFunc"
        @on-cancel="closeMigrateModal"
        :confirm-loading="state.migrate.loading"
        width="400px"
    >
      <n-alert type="info">
        请注意，迁移到新的客户端后，请确认新的客户端依然能正常访问到内网目标地址
      </n-alert>
      <br>
      <n-form>
        <n-form-item label="新客户端" path="clientCode">
          <n-select
              :options="state.clients"
              label-field="name"
              value-field="code"
              v-model:value="state.migrate.data.clientCode"
              :default-value="state.migrate.data.clientCode"
              placeholder="请选择目标客户端"
          ></n-select>
        </n-form-item>
      </n-form>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>