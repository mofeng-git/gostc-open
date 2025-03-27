<script setup>
import {onBeforeMount, ref} from "vue";
import {
  apiNormalGostClientP2PDelete,
  apiNormalGostClientP2PEnable,
  apiNormalGostClientP2PPage,
  apiNormalGostClientP2PRenew,
  apiNormalGostClientP2PUpdate
} from "../../../api/normal/gost_client_p2p.js";
import AppCard from "../../../layout/components/AppCard.vue";
import SearchCard from "../../../layout/components/SearchCard.vue";
import SearchItem from "../../../layout/components/SearchItem.vue";
import router from "../../../router/index.js";
import {regexpRule, requiredRule} from "../../../utils/formDataRule.js";
import {regexpLocalIp, regexpPort} from "../../../utils/regexp.js";
import Modal from "../../../components/Modal.vue";
import {cLimiterText, configText, limiterText, rLimiterText} from "../gost_client_forward/index.js";
import Empty from "../../../components/Empty.vue";
import Alert from "../../../icon/alert.vue";
import Online from "../../../icon/online.vue";
import {configExpText} from "../gost_client_host/index.js";
import {copyToClipboard} from "../../../utils/copy.js";
import {NButton, NSpace} from "naive-ui";

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
    },
    dataRules: {
      name: requiredRule('请输入名称'),
      targetIp: regexpRule(regexpLocalIp, '内网IP格式错误'),
      targetPort: regexpRule(regexpPort, '内网端口格式错误'),
    },
    open: false,
    loading: false,
  },
  look: {
    open: false,
    key: '',
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
    let res = await apiNormalGostClientP2PPage(state.value.table.search)
    state.value.table.list = res.data.list
    state.value.table.total = res.data.total
  } finally {
    state.value.table.searchLoading = false
  }
}

const openCreate = () => {
  router.push({name: 'NormalGostClientP2PCreate'})
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
        await apiNormalGostClientP2PUpdate(state.value.update.data)
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
    await apiNormalGostClientP2PEnable({code: row.code, enable: enable})
    refreshTable()
  } finally {
    row.enableLoading = false
  }
}

const renewFunc = async (row) => {
  try {
    row.renewLoading = true
    await apiNormalGostClientP2PRenew({code: row.code})
    refreshTable()
  } finally {
    row.renewLoading = false
  }
}

const deleteFunc = async (row) => {
  try {
    await apiNormalGostClientP2PDelete({code: row.code})
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

onBeforeMount(() => {
  pageFunc()
})

const generateCmdString = () => {
  let tls = ' --tls=false'
  if (window.location.protocol.indexOf('https') > 0) {
    tls = ''
  }
  return './gostc' + tls + ' -addr ' + window.location.host + ' -p2p -vts aaaaaa:8080,bbbbbb:8081'
}
</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-alert type="info">
        访客端运行命令：
        <div>{{ generateCmdString() }}</div>
        <div>含义：当前有两条P2P隧道，他们的访问密钥分别为aaaaaa、bbbbbb，访客端运行后，会把aaaaaa密钥的隧道配置的内网服务在访客端设备开启监听8080端口，访问8080端口就相当于访问隧道指向的内网服务，bbbbbb密钥的隧道同理</div>
        <div>其他问题：Linux可能会碰到权限问题，执行以下命令解决：sudo chmod +x gostc</div>
      </n-alert>
    </AppCard>
    <AppCard :show-border="false">
      <n-alert type="info">
        P2P隧道的方案采用的frp的stcp和xtcp组合，直连失败时，会使用节点转发，转发速率受套餐速率限制，直连成功则不受速率限制
      </n-alert>
    </AppCard>
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
<!--          <n-button-->
<!--              :focusable="false"-->
<!--              type="warning"-->
<!--              @click="goToUrl('https://docs.sian.one/gostc/tunnel')">-->
<!--            使用教程-->
<!--          </n-button>-->
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
              <span>速率：{{ limiterText(row.config.limiter) }}</span><br>
              <span>并发数：{{ rLimiterText(row.config.rLimiter) }}</span><br>
              <span>连接数：{{ cLimiterText(row.config.cLimiter) }}</span><br>
              <span>套餐：{{ configText(row.config) }}</span><br>
              <span>到期时间：{{ configExpText(row.config) }}</span><br>
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
              <n-button size="tiny" :focusable="false" quaternary type="info" @click="openLookFunc(row.vKey)">
                访问密钥
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

    <Modal title="访问密钥" :show="state.look.open"
           @on-confirm="copyFunc"
           confirm-text="复制"
           cancel-text="关闭"
           @on-cancel="state.look.open = false"
           :auto-focus="false"
    >
      <n-p>{{ state.look.key }}</n-p>
    </Modal>

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
        <n-form-item path="targetIp" label="内网IP">
          <n-input v-model:value="state.update.data.targetIp" placeholder="127.0.0.1"></n-input>
        </n-form-item>
        <n-form-item path="targetPort" label="内网端口">
          <n-input v-model:value="state.update.data.targetPort" placeholder="80"></n-input>
        </n-form-item>
      </n-form>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>