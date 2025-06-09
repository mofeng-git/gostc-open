<script setup>
import router from "../../../router/index.js";
import {onBeforeMount, ref} from "vue";
import {apiNormalFrpClientCfgCreate} from "../../../api/normal/frp_client_cfg.js";
import {requiredRule} from "../../../utils/formDataRule.js";
import {apiNormalGostClientList} from "../../../api/normal/gost_client.js";
import AppCard from "../../../layout/components/AppCard.vue";
import {NSpace} from "naive-ui";
import Online from "../../../icon/online.vue";


const state = ref({
  data: {
    clientCode: '',
    name: '',
    type: 'frpc',
  },
  dataRules: {
    name: requiredRule('请输入名称'),
    type: requiredRule('请选择配置文件类型'),
    clientCode: requiredRule('请选择客户端'),
  },
  loading: true,
  clients: [],
})

const createRef = ref()
const createFunc = () => {
  createRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.loading = true
        await apiNormalFrpClientCfgCreate(state.value.data)
        back()
      } finally {
        state.value.loading = false
      }
    }
  })
}

const clientListFunc = async () => {
  try {
    state.value.loading = true
    let res = await apiNormalGostClientList()
    state.value.clients = res.data || []
    if (state.value.clients.length > 0) {
      state.value.data.clientCode = state.value.clients[0].code
    } else {
      router.back()
      $message.create('暂无客户端信息', {
        type: "warning",
        closable: true,
        duration: 1500,
      })
    }
  } finally {
    state.value.loading = false
  }
}


const back = () => {
  router.back()
}

onBeforeMount(() => {
  clientListFunc()
})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-spin :show="state.loading">
        <n-h4 style="font-weight: bold">选择客户端：</n-h4>
        <n-alert v-if="state.clients.length===0" type="warning">没有客户端，请先新增一个客户端</n-alert>
        <n-radio-group v-else v-model:value="state.data.clientCode"
                       style="width: 100%">
          <n-grid x-gap="12" y-gap="12" cols="550:2 800:3 1400:4 1">
            <n-grid-item v-for="client in state.clients">
              <n-alert
                  type="info"
                  :show-icon="false"
                  :bordered="false"
                  style="height: 100%;cursor: pointer"
                  @click="state.data.clientCode = client.code"
              >
                <n-radio
                    :key="client.code"
                    :value="client.code"
                    style="width: 100%;"
                >
                  <n-space justify="space-between" style="width: 100%">
                    <Online :online="client.online===1"></Online>
                    <span>{{ client.name }}</span>
                  </n-space>
                </n-radio>
              </n-alert>
            </n-grid-item>
          </n-grid>
        </n-radio-group>
      </n-spin>
    </AppCard>

    <AppCard :show-border="false">
      <n-form ref="createRef" :rules="state.dataRules" :model="state.data">
        <n-form-item path="name" label="名称">
          <n-input v-model:value="state.data.name" placeholder="我的服务"></n-input>
        </n-form-item>
        <n-form-item path="type" label="配置文件类型">
          <n-select
              v-model:value="state.data.type"
              :options="[{label:'FRPC',value:'frpc'},{label:'FRPS',value:'frps'}]"
          ></n-select>
        </n-form-item>
        <n-form-item path="content" label="配置内容(兼容INI/TOML)">
          <n-input type="textarea" :autosize="{minRows:10,maxRows:20}" v-model:value="state.data.content"></n-input>
        </n-form-item>
      </n-form>
      <n-space>
        <n-button size="small" @click="back" :focusable="false">取消</n-button>
        <n-button size="small" @click="createFunc" type="primary" :focusable="false" :loading="state.loading">
          保存
        </n-button>
      </n-space>
    </AppCard>
  </div>
</template>

<style scoped lang="scss"></style>