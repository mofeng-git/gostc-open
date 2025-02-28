<script setup>
import {onBeforeMount, ref} from "vue";
import {apiAdminSystemConfigGost, apiAdminSystemConfigQuery} from "../../../api/admin/system_config.js";

const state = ref({
  data: {
    version: '',
    logger: '',
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
  } finally {
    state.value.submitLoading = false
  }
}

onBeforeMount(async () => {
  try {
    let res = await apiAdminSystemConfigQuery({kind: 'SystemConfigGost'})
    state.value.data = res.data
  } finally {

  }
})

</script>

<template>
  <div style="padding: 20px">
    <n-card title="基础配置">
      <n-form>
        <n-form-item label="客户端版本">
          <n-input v-model:value="state.data.version"></n-input>
        </n-form-item>
        <n-form-item label="记录日志">
          <n-switch checked-value="1" unchecked-value="2" v-model:value="state.data.logger" :round="false">
            <template #checked>开启</template>
            <template #unchecked>关闭</template>
          </n-switch>
        </n-form-item>
        <n-button type="success" size="small" @click="submit" :loading="state.submitLoading">保存</n-button>
      </n-form>
    </n-card>
  </div>
</template>

<style scoped lang="scss">

</style>