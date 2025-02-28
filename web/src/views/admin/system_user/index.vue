<script setup>
import {onBeforeMount, ref} from "vue";
import {
  apiAdminSystemUserCreate,
  apiAdminSystemUserDelete,
  apiAdminSystemUserPage,
  apiAdminSystemUserUpdate
} from "../../../api/admin/system_user.js";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import {generateTableColumns} from "./tableColumns.js";
import Modal from "../../../components/Modal.vue";
import {requiredRule} from "../../../utils/formDataRule.js";

const state = ref({
  table: {
    searchLoading: false,
    search: {
      page: 1,
      size: 10,
      account: '',
      admin: 0,
    },
    list: [],
    total: 0,
    minWidth: 1100,
  },
  create: {
    open: false,
    loading: false,
    data: {
      account: '',
      password: '',
      amount: '0',
    },
    dataRules: {
      account: requiredRule('账号不能为空'),
      password: requiredRule('密码不能为空'),
    },
  },
  update: {
    open: false,
    loading: false,
    data: {
      code: '',
      account: '',
      password: '',
      amount: '0',
    },
    dataRules: {
      account: requiredRule('账号不能为空'),
    },
  }
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
    let res = await apiAdminSystemUserPage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}

const openCreate = () => {
  state.value.create.data = {
    account: '',
    password: '',
    amount: '0',
  }
  state.value.create.open = true
}

const closeCreate = () => {
  state.value.create.open = false
}

const openUpdate = (row) => {
  state.value.update.open = true
  state.value.update.data = {
    code: row.code,
    account: row.account,
    password: '',
    amount: row.amount,
  }
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
        await apiAdminSystemUserCreate(state.value.create.data)
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
        await apiAdminSystemUserUpdate(state.value.update.data)
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
    await apiAdminSystemUserDelete({code: row.code})
    searchTable()
  } finally {
  }
}

const tableColumns = generateTableColumns(openUpdate, deleteFunc)

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
          label="角色"
          :default="0"
          :options="[
              {value:'全部',key:0},
              {value:'管理员',key:1},
              {value:'普通用户',key:2}
          ]"
          @onChange="value => state.table.search.admin=value"
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
        title="新增用户"
        width="400px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.create.loading"
        @on-confirm="createFunc"
        @on-cancel="closeCreate"
        mask-close
    >
      <n-form ref="createRef" :rules="state.create.dataRules" :model="state.create.data">
        <n-form-item path="account" label="账号">
          <n-input v-model:value.trim="state.create.data.account"></n-input>
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input v-model:value.trim="state.create.data.password"></n-input>
        </n-form-item>
        <n-form-item path="amount" label="积分">
          <n-input v-model:value.trim="state.create.data.amount"></n-input>
        </n-form-item>
      </n-form>
    </Modal>

    <Modal
        :show="state.update.open"
        title="修改用户"
        width="400px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.update.loading"
        @on-confirm="updateFunc"
        @on-cancel="closeUpdate"
        mask-close
    >
      <n-form ref="updateRef" :rules="state.update.dataRules" :model="state.update.data">
        <n-form-item path="account" label="账号">
          <n-input v-model:value.trim="state.update.data.account"></n-input>
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input v-model:value.trim="state.update.data.password"></n-input>
        </n-form-item>
        <n-form-item path="amount" label="积分">
          <n-input v-model:value.trim="state.update.data.amount"></n-input>
        </n-form-item>
      </n-form>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>