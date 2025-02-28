<script setup>
import {onBeforeMount, ref} from "vue";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import {generateTableColumns} from "./tableColumns.js";
import router from "../../../router/index.js";
import {apiAdminGostNodeLoggerPage} from "../../../api/admin/gost_node_logger.js";

const state = ref({
  table: {
    searchLoading: false,
    search: {
      page: 1,
      size: 10,
      level: '',
      nodeCode: '',
    },
    list: [],
    total: 0,
    minWidth: 1100,
  },
  levels: [
    {value: 'ALL', key: ''},
    {value: 'INFO', key: 'info'},
    {value: 'WARNING', key: 'warning'},
    {value: 'ERROR', key: 'error'},
  ]
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
    if (state.value.table.search.nodeCode === '') {
      return
    }
    state.value.table.searchLoading = true
    let res = await apiAdminGostNodeLoggerPage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}

const tableColumns = generateTableColumns()

onBeforeMount(() => {
  state.value.table.search.nodeCode = router.currentRoute.value.query.nodeCode
  pageFunc()
})

</script>

<template>
  <div>
    <SearchCard :show-border="false" space>
      <SearchItem
          type="select"
          :label-width="120"
          label="日志级别"
          :options="state.levels"
          @onChange="value => state.table.search.level=value"
      ></SearchItem>
      <SearchItem type="custom">
        <n-space>
          <n-button type="info" :focusable="false" @click="searchTable">搜索</n-button>
          <n-button type="info" :focusable="false" @click="refreshTable">刷新</n-button>
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
  </div>
</template>

<style scoped lang="scss">

</style>