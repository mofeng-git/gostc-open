<script setup>
import {onBeforeMount, ref} from "vue";
import {
  apiAdminSystemNoticeCreate,
  apiAdminSystemNoticeDelete,
  apiAdminSystemNoticePage,
  apiAdminSystemNoticeUpdate
} from "../../../api/admin/system_notice.js";
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
    },
    list: [],
    total: 0,
    minWidth: 800,
  },
  create: {
    open: false,
    loading: false,
    data: {
      title: '',
      content: '',
      hidden: 1,
      indexValue: 1000,
    },
    dataRules: {
      title: requiredRule('标题不能为空'),
      content: requiredRule('内容不能为空'),
    },
  },
  update: {
    open: false,
    loading: false,
    data: {
      title: '',
      content: '',
      hidden: 1,
      indexValue: 1000,
    },
    dataRules: {
      title: requiredRule('标题不能为空'),
      content: requiredRule('内容不能为空'),
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
    let res = await apiAdminSystemNoticePage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}

const openCreate = () => {
  state.value.create.data = {
    title: '',
    content: '',
    hidden: 1,
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
        await apiAdminSystemNoticeCreate(state.value.create.data)
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
        await apiAdminSystemNoticeUpdate(state.value.update.data)
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
    await apiAdminSystemNoticeDelete({code: row.code})
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
      <SearchItem type="custom">
        <n-space>
          <!--          <n-button type="info" :focusable="false" @click="searchTable">搜索</n-button>-->
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
        title="新增通知公告"
        width="400px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.create.loading"
        @on-confirm="createFunc"
        @on-cancel="closeCreate"
        mask-close
    >
      <n-form ref="createRef" :rules="state.create.dataRules" :model="state.create.data">
        <n-form-item path="title" label="标题">
          <n-input v-model:value.trim="state.create.data.title"></n-input>
        </n-form-item>
        <n-form-item path="content" label="内容">
          <n-input type="textarea" :autosize="{minRows:5,maxRows:20}"
                   v-model:value.trim="state.create.data.content"></n-input>
        </n-form-item>
        <n-form-item path="hidden" label="状态">
          <n-switch
              :value="state.create.data.hidden"
              :checked-value="1"
              :unchecked-value="2"
              :on-update:value="value => state.create.data.hidden = value">
            <template #checked>隐藏</template>
            <template #unchecked>显示</template>
          </n-switch>
        </n-form-item>
        <n-form-item path="indexValue" label="排序(升序)">
          <n-input-number v-model:value="state.create.data.indexValue"></n-input-number>
        </n-form-item>
      </n-form>
    </Modal>

    <Modal
        :show="state.update.open"
        title="修改通知公告"
        width="400px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.update.loading"
        @on-confirm="updateFunc"
        @on-cancel="closeUpdate"
        mask-close
    >
      <n-form ref="updateRef" :rules="state.update.dataRules" :model="state.update.data">
        <n-form-item path="title" label="标题">
          <n-input v-model:value.trim="state.update.data.title"></n-input>
        </n-form-item>
        <n-form-item path="content" label="内容">
          <n-input type="textarea" :autosize="{minRows:5,maxRows:20}"
                   v-model:value.trim="state.update.data.content"></n-input>
        </n-form-item>
        <n-form-item path="hidden" label="状态">
          <n-switch
              :value="state.update.data.hidden"
              :checked-value="1"
              :unchecked-value="2"
              :on-update:value="value => state.update.data.hidden = value">
            <template #checked>隐藏</template>
            <template #unchecked>显示</template>
          </n-switch>
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