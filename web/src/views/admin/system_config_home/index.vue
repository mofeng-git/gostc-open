<script setup>
import {onBeforeMount, ref} from "vue";
import {apiAdminSystemConfigHome, apiAdminSystemConfigQuery} from "../../../api/admin/system_config.js";

const state = ref({
  data: {
    homeEnable: '1',
    homeTpl: '',
  },
  submitLoading: false,
})

const submit = async () => {
  try {
    state.value.submitLoading = true
    await apiAdminSystemConfigHome(state.value.data)
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
    let res = await apiAdminSystemConfigQuery({kind: 'SystemConfigHome'})
    state.value.data = res.data || {homeEnable:'1', homeTpl:''}
  } finally {
  }
}

onBeforeMount(() => {
  querySystemConfigFunc()
})
</script>

<template>
  <div style="padding: 20px">
    <n-card title="首页配置">
      <n-form>
        <n-form-item label="动态首页">
          <n-switch
              :round="false"
              v-model:value="state.data.homeEnable"
              checked-value="1"
              unchecked-value="2"
          >
            <template #checked>开启</template>
            <template #unchecked>关闭</template>
          </n-switch>
        </n-form-item>

        <n-form-item label="自定义HTML模板">
          <n-space vertical style="width: 100%">
            <n-alert :show-icon="false" type="info">
              模板变量（Go模板语法）：<br>
              <n-space wrap>
                <n-tag size="small"><span v-pre>{{ .config.title }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .config.favicon }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .stats.today_flow }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .stats.user_total }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .stats.checkin_today }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .stats.tunnel_total }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .stats.node_online }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .stats.client_online }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .stats.updated_at }}</span></n-tag>
                <n-tag size="small"><span v-pre>{{ .login_url }}</span></n-tag>
              </n-space>
            </n-alert>
            <n-input
                type="textarea"
                :autosize="{minRows:10,maxRows:28}"
                v-model:value="state.data.homeTpl"
                placeholder="留空将使用内置默认模板"
            ></n-input>
          </n-space>
        </n-form-item>

        <n-button type="success" size="small" @click="submit" :loading="state.submitLoading">保存</n-button>
      </n-form>
    </n-card>
  </div>
  
</template>

<style scoped lang="scss">

</style>

