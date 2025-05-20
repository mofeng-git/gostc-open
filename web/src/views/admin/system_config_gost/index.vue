<script setup>
import {onBeforeMount, ref} from "vue";
import {apiAdminSystemConfigGost, apiAdminSystemConfigQuery} from "../../../api/admin/system_config.js";

const state = ref({
  data: {
    version: '',
    logger: '2',
    funcWeb:"2",
    funcForward:"2",
    funcTunnel:"2",
    funcP2P:"2",
    funcProxy:"2",
    funcTun:"2",
    funcNode:"2",
  },
  submitLoading: false,
})

const submit = async () => {
  try {
    state.value.submitLoading = true
    await apiAdminSystemConfigGost(state.value.data)
    $message.success('保存成功，刷新生效', {
      showIcon: true,
      closable: true,
      duration: 1500,
    })
    await querySystemConfigFunc()
  } finally {
    state.value.submitLoading = false
  }
}

const querySystemConfigFunc = async ()=>{
  try {
    let res = await apiAdminSystemConfigQuery({kind: 'SystemConfigGost'})
    state.value.data = res.data
  } finally {

  }
}

onBeforeMount(() => {
  querySystemConfigFunc()
})

</script>

<template>
  <div style="padding: 20px">
    <n-card title="基础配置">
      <n-form>
        <n-form-item label="客户端版本(仅用于提示用户客户端版本)">
          <n-input v-model:value="state.data.version"></n-input>
        </n-form-item>
        <n-form-item label="功能菜单(显示/隐藏)">
          <n-space>
            <n-checkbox checked-value="1" unchecked-value="2" v-model:checked="state.data.funcWeb">域名解析</n-checkbox>
            <n-checkbox checked-value="1" unchecked-value="2" v-model:checked="state.data.funcForward">端口转发</n-checkbox>
            <n-checkbox checked-value="1" unchecked-value="2" v-model:checked="state.data.funcTunnel">私有隧道</n-checkbox>
            <n-checkbox checked-value="1" unchecked-value="2" v-model:checked="state.data.funcP2P">P2P隧道</n-checkbox>
            <n-checkbox checked-value="1" unchecked-value="2" v-model:checked="state.data.funcProxy">代理隧道</n-checkbox>
<!--            <n-checkbox checked-value="1" unchecked-value="2" v-model:checked="state.data.funcTun">虚拟组网</n-checkbox>-->
            <n-checkbox checked-value="1" unchecked-value="2" v-model:checked="state.data.funcNode">自建节点</n-checkbox>
          </n-space>
        </n-form-item>
        <n-button type="success" size="small" @click="submit" :loading="state.submitLoading">保存</n-button>
      </n-form>
    </n-card>
  </div>
</template>

<style scoped lang="scss">

</style>