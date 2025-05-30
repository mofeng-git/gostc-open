<script setup>
import {onBeforeMount, ref} from "vue";
import {
  apiAdminSystemConfigEmail,
  apiAdminSystemConfigEmailVerify,
  apiAdminSystemConfigQuery
} from "../../../api/admin/system_config.js";
import Modal from "../../../components/Modal.vue";
import {requiredRule} from "../../../utils/formDataRule.js";

const state = ref({
  data: {
    enable: '2',
    nickName: '',
    host: '',
    port: '',
    user: "",
    pwd: "",
    resetPwdTpl: "",
  },
  submitLoading: false,
  verify: {
    data: {
      email: '',
    },
    dataRules: {
      email: requiredRule('邮件不能为空'),
    },
    loading: false,
    open: false,
    openLoading: false,
  },
})


const openVerify = () => {
  state.value.verify.data = {
    email: '',
  }
  state.value.verify.open = true
}

const closeVerify = () => {
  state.value.verify.open = false
}

const verifyRef = ref()
const verifyFunc = () => {
  verifyRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.verify.loading = true
        await apiAdminSystemConfigEmailVerify(state.value.verify.data)
        $message.create('发送成功', {
          type: "success",
          duration: 1500,
          closable: true
        })
      } finally {
        state.value.verify.loading = false
      }
    }
  })
}

const submit = async () => {
  try {
    state.value.submitLoading = true
    await apiAdminSystemConfigEmail(state.value.data)
    $message.success('保存成功', {
      showIcon: true,
      closable: true,
      duration: 1500,
    })
    await querySystemConfigFunc()
  } finally {
    state.value.submitLoading = false
  }
}

const querySystemConfigFunc = async () => {
  try {
    let res = await apiAdminSystemConfigQuery({kind: 'SystemConfigEmail'})
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
    <n-card title="邮件配置">
      <n-form>
        <n-form-item label="启用/禁用">
          <n-switch
              :round="false"
              v-model:value="state.data.enable"
              checked-value="1"
              unchecked-value="2"
          >
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
        <n-form-item label="昵称">
          <n-input v-model:value="state.data.nickName" placeholder="管理员"></n-input>
        </n-form-item>
        <n-form-item label="地址">
          <n-input v-model:value="state.data.host" placeholder="smtp.qq.com"></n-input>
        </n-form-item>
        <n-form-item label="端口">
          <n-input v-model:value="state.data.port" placeholder="465"></n-input>
        </n-form-item>
        <n-form-item label="账号">
          <n-input v-model:value="state.data.user" placeholder="admin@example.com"></n-input>
        </n-form-item>
        <n-form-item label="密码">
          <n-input v-model:value="state.data.pwd"></n-input>
        </n-form-item>
<!--        <n-form-item label="重置密码模板">-->
<!--          <n-space vertical style="width: 100%">-->
<!--            <n-alert :show-icon="false" type="info">-->
<!--              模板变量：<br>-->
<!--              <n-tag size="small"><span v-pre>{{CODE}}</span></n-tag>：验证码<br>-->
<!--              <n-tag size="small"><span v-pre>{{DATETIME}}</span></n-tag>：发件时间<br>-->
<!--            </n-alert>-->
<!--            <n-input type="textarea" :autosize="{minRows:5,maxRows:15}"-->
<!--                     v-model:value="state.data.resetPwdTpl"></n-input>-->
<!--          </n-space>-->
<!--        </n-form-item>-->
        <n-space>
          <n-button type="success" size="small" @click="submit" :loading="state.submitLoading">保存</n-button>
          <n-button type="info" size="small" @click="openVerify">测试(保存后再点击测试)</n-button>
        </n-space>
      </n-form>
    </n-card>

    <Modal title="测试邮箱" :show="state.verify.open"
           @on-confirm="verifyFunc"
           :confirm-loading="state.verify.loading"
           @on-cancel="closeVerify"
           width="400px"
    >
      <n-form
          ref="verifyRef"
          :model="state.verify.data"
          :rules="state.verify.dataRules"
          :show-label="true"
          size="medium"
      >
        <n-form-item path="email" label="邮箱">
          <n-input
              v-model:value="state.verify.data.email"
              placeholder="admin@example.com"
          />
        </n-form-item>
      </n-form>
    </Modal>
  </div>
</template>

<style scoped lang="scss">

</style>