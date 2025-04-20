<script setup xmlns="http://www.w3.org/1999/html">
import {onBeforeMount, ref} from "vue";
import {apiTunnelCreate, apiTunnelDelete, apiTunnelList, apiTunnelStatus, apiTunnelUpdate} from "../../api/tunnel.js";
import Modal from "../../components/Modal.vue";
import {requiredRule} from "../../utils/formDataRule.js";

const state = ref({
  data: [],
  create: {
    data: {
      key: '',
      name: '',
      bind: '',
      port: '',
      tls: 1,
      address: 'gost.sian.one',
      autoStart:0,
    },
    dataRules: {
      key: requiredRule('请输入密钥'),
      name: requiredRule('请输入备注名称'),
      address: requiredRule('请输入服务器地址'),
      port: requiredRule('请输入本地端口'),
    },
    open: false,
  },
  update: {
    data: {
      key: '',
      name: '',
      bind: '',
      port: '',
      tls: 1,
      address: 'gost.sian.one',
      autoStart:0,
    },
    dataRules: {
      key: requiredRule('请输入密钥'),
      name: requiredRule('请输入备注名称'),
      address: requiredRule('请输入服务器地址'),
      port: requiredRule('请输入本地端口'),
    },
    open: false,
  }
})

const listFunc = async () => {
  try {
    let res = await apiTunnelList()
    state.value.data = res.data || []
  } finally {

  }
}

const openCreateModal = () => {
  state.value.create.data = {
    key: '',
    name: '',
    bind: '',
    port: '',
    tls: 1,
    address: 'gost.sian.one',
    autoStart:0,
  }
  state.value.create.open = true
}

const closeCreateModal = () => {
  state.value.create.open = false
}

const createRef = ref()
const confirmCreateFunc = () => {
  createRef.value.validate(async valid => {
    if (!valid) {
      try {
        await apiTunnelCreate(state.value.create.data)
        closeCreateModal()
        await listFunc()
      } finally {

      }
    }
  })
}

const openUpdateModal = (row) => {
  state.value.update.data = JSON.parse(JSON.stringify(row))
  state.value.update.open = true
}

const closeUpdateModal = () => {
  state.value.update.open = false
}

const updateRef = ref()
const confirmUpdateFunc = () => {
  updateRef.value.validate(async valid => {
    if (!valid) {
      try {
        await apiTunnelUpdate(state.value.update.data)
        closeUpdateModal()
        await listFunc()
      } finally {

      }
    }
  })
}

const statusChangeFunc = async (row, value) => {
  try {
    row.statsuLoading = true
    await apiTunnelStatus({
      key: row.key,
      status: value,
    })
    await listFunc()
  } finally {
    row.statsuLoading = false
  }
}

const deleteFunc = async (row) => {
  try {
    row.deleteLoading = true
    await apiTunnelDelete({key: row.key})
    await listFunc()
  } finally {
    row.deleteLoading = false
  }
}

onBeforeMount(() => {
  listFunc()
})

const generateBindAddressString = (bind,port) => {
  let address = ''
  if (bind === '' || bind === '0.0.0.0'){
    address = window.location.hostname
  }else {
    address = bind
  }
  return address+':'+port
}
const generateServerString = (address, tls) => {
  return tls === 1 ? `wss://${address}` : `ws://${address}`
}
</script>

<template>
  <div>
    <n-space>
      <n-button size="small" type="info" @click="listFunc">刷新</n-button>
      <n-button size="small" type="info" @click="openCreateModal">新增</n-button>
    </n-space>

    <p/>

    <n-grid x-gap="12" y-gap="12" cols="1 520:2 900:3 1400:4">
      <n-grid-item v-for="row in state.data">
        <n-el class="client-item" tag="div" :style="{
                border: '1px solid var(--border-color)',
                borderRadius:'var(--border-radius)',
                padding: '12px',
                cursor: 'pointer'}">
          <n-h4 style="margin-bottom: 8px !important;">
            <n-space justify="space-between">
              <span style="font-weight: bold">{{ row.name }}</span>
              <div style="display: flex;justify-content: center">
                <n-switch
                    style="margin-left: 8px"
                    size="small"
                    v-model:value="row.status"
                    :loading="row.statsuLoading"
                    :checked-value="1"
                    :unchecked-value="2"
                    :on-update:value="value => {statusChangeFunc(row,value)}"
                ></n-switch>
              </div>
            </n-space>
          </n-h4>
          <div>
            <span>服务器：{{ generateServerString(row.address, row.tls) }}</span><br>
            <span>密钥：{{ row.key }}</span><br>
            <span>地址1：{{ generateBindAddressString(row.bind,row.port) }}</span><br>
            <span>地址2：{{ generateBindAddressString('',row.port) }}</span><br>
            <span>自启动：{{ row.autoStart===1?'是':'否' }}</span><br>
          </div>
          <n-space justify="end" style="width: 100%">
            <n-button
                size="tiny"
                :focusable="false"
                quaternary type="success"
                @click="openUpdateModal(row)"
                :disabled="row.status===1"
            >编辑
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
                    :disabled="row.status===1"
                >删除
                </n-button>
              </template>
              确认删除吗？
            </n-popconfirm>
          </n-space>
        </n-el>
      </n-grid-item>
    </n-grid>

    <Modal
        :show="state.create.open"
        title="新增私有隧道"
        width="600px"
        confirm-text="保存"
        cancel-text="取消"
        @on-confirm="confirmCreateFunc"
        @on-cancel="closeCreateModal"
        mask-close
    >
      <n-form ref="createRef" :rules="state.create.dataRules" :model="state.create.data">
        <n-form-item path="name" label="名称">
          <n-input v-model:value.trim="state.create.data.name"></n-input>
        </n-form-item>
        <n-form-item path="key" label="密钥">
          <n-input v-model:value.trim="state.create.data.key"></n-input>
        </n-form-item>
        <n-form-item path="bind" label="本地地址">
          <n-input v-model:value.trim="state.create.data.bind" placeholder="0.0.0.0"></n-input>
        </n-form-item>
        <n-form-item path="port" label="本地端口">
          <n-input v-model:value.trim="state.create.data.port" placeholder="8080"></n-input>
        </n-form-item>
        <n-form-item path="address" label="服务器">
          <n-input-group>
            <n-select
                :style="{ width: '33%' }"
                :options="[{label:'启用TLS',value:1},{label:'禁用TLS',value:2}]"
                v-model:value="state.create.data.tls"/>
            <n-input v-model:value.trim="state.create.data.address"></n-input>
          </n-input-group>
        </n-form-item>
        <n-form-item path="autoStart" label="自启动">
          <n-switch :round="false" v-model:value="state.create.data.autoStart" :unchecked-value="0" :checked-value="1">
            <template #checked>开启</template>
            <template #unchecked>关闭</template>
          </n-switch>
        </n-form-item>
      </n-form>
    </Modal>

    <Modal
        :show="state.update.open"
        title="修改私有隧道"
        width="600px"
        confirm-text="保存"
        cancel-text="取消"
        @on-confirm="confirmUpdateFunc"
        @on-cancel="closeUpdateModal"
        mask-close
    >
      <n-form ref="updateRef" :rules="state.update.dataRules" :model="state.update.data">
        <n-form-item path="name" label="名称">
          <n-input v-model:value.trim="state.update.data.name"></n-input>
        </n-form-item>
        <n-form-item path="key" label="密钥">
          <n-input v-model:value.trim="state.update.data.key" disabled></n-input>
        </n-form-item>
        <n-form-item path="bind" label="本地地址">
          <n-input v-model:value.trim="state.update.data.bind" placeholder="0.0.0.0"></n-input>
        </n-form-item>
        <n-form-item path="port" label="本地端口">
          <n-input v-model:value.trim="state.update.data.port" placeholder="8080"></n-input>
        </n-form-item>
        <n-form-item path="address" label="服务器">
          <n-input-group>
            <n-select
                :style="{ width: '33%' }"
                :options="[{label:'启用TLS',value:1},{label:'禁用TLS',value:2}]"
                v-model:value="state.update.data.tls"/>
            <n-input v-model:value.trim="state.update.data.address"></n-input>
          </n-input-group>
        </n-form-item>
        <n-form-item path="autoStart" label="自启动">
          <n-switch :round="false" v-model:value="state.update.data.autoStart" :unchecked-value="0" :checked-value="1">
            <template #checked>开启</template>
            <template #unchecked>关闭</template>
          </n-switch>
        </n-form-item>
      </n-form>
    </Modal>
  </div>
</template>
