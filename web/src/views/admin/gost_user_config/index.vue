<script setup>
import {computed, onBeforeMount, ref} from "vue";
import {
  apiAdminGostUserConfigCreate,
  apiAdminGostUserConfigDelete,
  apiAdminGostUserConfigPage,
  apiAdminGostUserConfigUpdate
} from "../../../api/admin/gost_user_config.js";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import {generateTableColumns} from "./tableColumns.js";
import Modal from "../../../components/Modal.vue";
import {requiredRule} from "../../../utils/formDataRule.js";
import {apiAdminGostNodeList} from "../../../api/admin/gost_node.js";
import {apiAdminGostConfigList} from "../../../api/admin/gost_config.js";
import {apiAdminSystemUserList} from "../../../api/admin/system_user.js";

const state = ref({
  table: {
    searchLoading: false,
    search: {
      page: 1,
      size: 10,
      name: '',
    },
    list: [],
    total: 0,
    minWidth: 1100,
  },
  create: {
    open: false,
    loading: false,
    configCode: '',
    data: {
      userCode: '',
      name: '',
      chargingType: 1,
      cycle: 10,
      amount: '0',
      limiter: 1,
      rLimiter: 100,
      cLimiter: 100,
      onlyChina: 1,
      nodes: [],
      expAt: new Date().getTime(),
    },
    dataRules: {
      userCode: requiredRule('请选择用户'),
      name: requiredRule('请输入名称'),
    },
    users: [],
    usersLoading: false,
  },
  update: {
    open: false,
    loading: false,
    data: {
      code: '',
      name: '',
      userCode: '',
      userAccount: '',
      chargingType: 1,
      cycle: 10,
      amount: '0',
      limiter: 1,
      rLimiter: 100,
      cLimiter: 100,
      onlyChina: 1,
      nodes: [],
      expAt: new Date().getTime(),
    },
    dataRules: {
      name: requiredRule('请输入名称'),
    },
  },
  nodes: [],
  configs: [],
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
    let res = await apiAdminGostUserConfigPage(state.value.table.search)
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
    chargingType: 1,
    cycle: 10,
    amount: '0',
    limiter: 1,
    rLimiter: 100,
    cLimiter: 100,
    onlyChina: 1,
    nodes: [],
    expAt: new Date().getTime(),
  }
  state.value.create.configCode = ''
  state.value.create.open = true
}

const closeCreate = () => {
  state.value.create.open = false
}

const selectConfigChange = (code) => {
  for (let i = 0; i < state.value.configs.length; i++) {
    let cfg = state.value.configs[i]
    if (cfg.code === code) {
      state.value.create.data.name = cfg.name
      state.value.create.data.chargingType = cfg.chargingType
      state.value.create.data.cycle = cfg.cycle
      state.value.create.data.amount = cfg.amount
      state.value.create.data.limiter = cfg.limiter
      state.value.create.data.rLimiter = cfg.rLimiter
      state.value.create.data.cLimiter = cfg.cLimiter
      state.value.create.data.onlyChina = cfg.onlyChina
      state.value.create.data.nodes = cfg.nodes
      return
    }
  }
}

const openUpdate = (row) => {
  state.value.update.data = JSON.parse(JSON.stringify(row))
  state.value.update.data.expAt = state.value.update.data.expAt * 1000
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
        await apiAdminGostUserConfigCreate(Object.assign({}, state.value.create.data, {expAt: state.value.create.data.expAt / 1000}))
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
        await apiAdminGostUserConfigUpdate(Object.assign({}, state.value.update.data, {expAt: state.value.update.data.expAt / 1000}))
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
    await apiAdminGostUserConfigDelete({code: row.code})
    searchTable()
  } finally {
  }
}

const getNodes = async () => {
  try {
    let res = await apiAdminGostNodeList()
    state.value.nodes = res.data?.map(item => {
      return {
        label: item.name,
        value: item.code,
      }
    }) || []
  } finally {

  }
}

const getConfigs = async () => {
  try {
    let res = await apiAdminGostConfigList({})
    state.value.configs = res.data || []
  } finally {

  }
}

const cfgs = computed(() => {
  return state.value.configs?.map(item => {
    return {
      label: item.name,
      value: item.code,
    }
  })
})

const searchUsers = async (account) => {
  if (!account) {
    return
  }
  try {
    state.value.create.usersLoading = true
    let res = await apiAdminSystemUserList({account: account})
    state.value.create.users = res.data || []
  } finally {
    state.value.create.usersLoading = false
  }
}

const tableColumns = generateTableColumns(openUpdate, deleteFunc)

onBeforeMount(() => {
  pageFunc()
  getNodes()
  getConfigs()
})

</script>

<template>
  <div>
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
        </n-space>
      </SearchItem>
    </SearchCard>

    <AppCard>
      <n-data-table
          remote
          :columns="tableColumns"
          :data="state.table.list"
          :bordered="false"
          :loading="state.table.searchLoading"
          :scroll-x="state.table.minWidth"
      />
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
        title="新增套餐"
        width="600px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.create.loading"
        @on-confirm="createFunc"
        @on-cancel="closeCreate"
    >
      <n-form ref="createRef" :rules="state.create.dataRules" :model="state.create.data">
        <n-form-item path="userCode" label="用户">
          <n-select
              v-model:value="state.create.data.userCode"
              filterable
              placeholder="搜索用户"
              :options="state.create.users"
              :loading="state.create.usersLoading"
              label-field="account"
              value-field="code"
              remote
              @search="searchUsers"
          />
        </n-form-item>
        <n-popselect
            v-model:value="state.create.configCode"
            :options="cfgs"
            trigger="click"
            :on-update:value="selectConfigChange"
        >
          <n-button type="info" size="small">快速应用套餐配置</n-button>
        </n-popselect>
        <p/>
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
        <n-form-item path="onlyChina" label="仅中国大陆客户端可用">
          <n-switch :checked-value="1" :unchecked-value="2" :default-value="state.create.data.onlyChina">
            <template #checked>开启</template>
            <template #unchecked>关闭</template>
          </n-switch>
        </n-form-item>
        <n-form-item path="limiter" label="速率(mbps)">
          <n-input-number v-model:value="state.create.data.limiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="rLimiter" label="并发数">
          <n-input-number v-model:value="state.create.data.rLimiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="cLimiter" label="连接数">
          <n-input-number v-model:value="state.create.data.cLimiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="expAt" label="到期时间" v-if="state.create.data.chargingType===2">
          <n-date-picker v-model:value="state.create.data.expAt" type="date"/>
        </n-form-item>
        <n-transfer v-model:value="state.create.data.nodes" :options="state.nodes"/>
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
        <n-form-item path="userAccount" label="用户">
          <n-input :default-value="state.update.data.userAccount" disabled></n-input>
        </n-form-item>
        <n-form-item path="chargingType" label="计费方式">
          <n-radio-group v-model:value="state.update.data.chargingType">
            <n-radio :checked="state.update.data.chargingType===1" :value="1">
              一次性
            </n-radio>
            <n-radio :checked="state.update.data.chargingType===2" :value="2">
              循环
            </n-radio>
            <n-radio :checked="state.update.data.chargingType===3" :value="3">
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
        <n-form-item path="onlyChina" label="仅中国大陆客户端可用">
          <n-switch :checked-value="1" :unchecked-value="2" v-model:value="state.update.data.onlyChina">
            <template #checked>开启</template>
            <template #unchecked>关闭</template>
          </n-switch>
        </n-form-item>
        <n-form-item path="limiter" label="速率(mbps)">
          <n-input-number v-model:value="state.update.data.limiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="rLimiter" label="并发数">
          <n-input-number v-model:value="state.update.data.rLimiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="cLimiter" label="连接数">
          <n-input-number v-model:value="state.update.data.cLimiter" :min="0"></n-input-number>
        </n-form-item>
        <n-form-item path="expAt" label="到期时间" v-if="state.update.data.chargingType===2">
          <n-date-picker v-model:value="state.update.data.expAt" type="date"/>
        </n-form-item>
        <n-transfer v-model:value="state.update.data.nodes" :options="state.nodes"/>
      </n-form>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>