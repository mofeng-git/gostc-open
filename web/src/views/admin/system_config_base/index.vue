<script setup>
import {onBeforeMount, ref} from "vue";
import {apiAdminSystemConfigBase, apiAdminSystemConfigQuery} from "../../../api/admin/system_config.js";
import {randomString} from "../../../utils/random.js";

const state = ref({
  data: {
    title: '',
    favicon: '',
    baseUrl: '',
    checkIn:'2',
    checkInStart:0,
    checkInEnd:0,
    register:'2',
  },
  submitLoading: false,
})

const submit = async () => {
  try {
    state.value.submitLoading = true
    await apiAdminSystemConfigBase(state.value.data)
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

const generateApiKeyFunc = ()=>{
  state.value.data.apiKey = randomString(32)
}

const querySystemConfigFunc = async ()=>{
  try {
    let res = await apiAdminSystemConfigQuery({kind: 'SystemConfigBase'})
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
        <n-form-item label="标题">
          <n-input v-model:value="state.data.title"></n-input>
        </n-form-item>
        <n-form-item label="Favicon">
          <n-input v-model:value="state.data.favicon"></n-input>
        </n-form-item>
        <n-form-item label="基础URL(暂无实际作用，预留)">
          <n-input v-model:value="state.data.baseUrl"></n-input>
        </n-form-item>
        <n-form-item label="注册">
          <n-switch
              :round="false"
              v-model:value="state.data.register"
              checked-value="1"
              unchecked-value="2"
          >
            <template #checked>开启</template>
            <template #unchecked>关闭</template>
          </n-switch>
        </n-form-item>
        <n-form-item label="签到">
          <n-switch
              :round="false"
              v-model:value="state.data.checkIn"
              checked-value="1"
              unchecked-value="2"
          >
            <template #checked>开启</template>
            <template #unchecked>关闭</template>
          </n-switch>
        </n-form-item>
        <n-form-item label="签到积分范围">
          <n-input-group>
            <n-input-number :min="0" placeholder="最小" v-model:value="state.data.checkInStart"></n-input-number>
            <n-input-number :min="0" placeholder="最多" v-model:value="state.data.checkInEnd"></n-input-number>
          </n-input-group>
        </n-form-item>
        <n-button type="success" size="small" @click="submit" :loading="state.submitLoading">保存</n-button>
      </n-form>
    </n-card>
  </div>
</template>

<style scoped lang="scss">

</style>