<script setup>
import {onBeforeMount, ref} from "vue";
import {apiAuthCaptcha, apiAuthGenResetPwdEmailCode, apiAuthResetPwd} from "../../api/auth/index.js";
import {appStore} from "../../store/app.js";
import {requiredRule} from "../../utils/formDataRule.js";
import router from "../../router/index.js";

const state = ref({
  captcha: {
    src: '',
    loading: false,
    security: true,
  },
  resetPwd: {
    loading: false,
    genLoading: false,
    data: {
      account: '',
      captchaKey: '',
      captchaValue: '',
      key: '',
      code: ''
    },
    dataRules: {
      account: requiredRule('请输入账号'),
      captchaValue: requiredRule('请输入图形验证码'),
      code: requiredRule('请输入邮件验证码'),
    }
  }
})

const refreshCaptchaFunc = async () => {
  try {
    state.value.captcha.loading = true
    let res = await apiAuthCaptcha()
    state.value.resetPwd.data.captchaKey = res.data.key
    state.value.captcha.src = res.data?.bs64
    state.value.captcha.security = res.data?.security
  } finally {
    state.value.captcha.loading = false
  }
}

const resetPwdRef = ref()
const genResetPwdEmailCodeFunc = () => {
  resetPwdRef.value.validate(async (errors) => {
    if (!errors) {
      try {
        state.value.resetPwd.genLoading = true
        let res = await apiAuthGenResetPwdEmailCode(state.value.resetPwd.data)
        state.value.resetPwd.data.key = res.data
      } catch (err) {
        await refreshCaptchaFunc()
      } finally {
        state.value.resetPwd.genLoading = false
      }
    }
  });
}

const resetPwdFunc = () => {
  resetPwdRef.value.validate(async (errors) => {
    if (!errors) {
      try {
        state.value.resetPwd.loading = true
        await apiAuthResetPwd(state.value.resetPwd.data)
        $message.create('重置成功，新密码已发送至您的邮箱，请使用新密码登录', {
          type: "success",
          closable: true,
          duration: 1500,
        })
      } catch (err) {

      } finally {
        state.value.resetPwd.loading = false
      }
    }
  });
}

onBeforeMount(() => {
  refreshCaptchaFunc()
})

</script>

<template>
  <div class="resetPwd">
    <n-card class="bg">
      <div class="title" style="user-select: none;cursor: pointer;font-size: 1.5em">
        {{ appStore().siteConfig.title }}
      </div>
      <n-divider/>
      <n-tabs animated>
        <n-tab-pane tab="重置密码" name="resetPwd">
          <n-form
              ref="resetPwdRef"
              :model="state.resetPwd.data"
              :rules="state.resetPwd.dataRules"
              :show-label="false"
              size="large"
          >
            <n-form-item path="account" label="账号">
              <n-input
                  v-model:value="state.resetPwd.data.account"
                  placeholder="账号"
              />
            </n-form-item>
            <n-form-item path="captchaValue" label="验证码">
              <n-input
                  v-model:value="state.resetPwd.data.captchaValue"
                  placeholder="验证码"
              />
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-spin :show="state.captcha.loading" size="small" :delay="200" style="height: 40px">
                    <n-image
                        height="40"
                        @click="refreshCaptchaFunc"
                        style="border-radius: 3px;margin-left: 8px"
                        :src="'data:image/jpeg;base64,' + state.captcha.src"
                        preview-disabled
                    />
                  </n-spin>
                </template>
                点击刷新
              </n-tooltip>
            </n-form-item>
            <n-button style="width: 100%" type="primary" @click="genResetPwdEmailCodeFunc"
                      :loading="state.resetPwd.genLoading"
                      :disabled="state.resetPwd.data.captchaValue===''">
              {{ state.resetPwd.data.key ? '重新发送' : '发送邮箱验证码' }}
            </n-button>

            <div v-if="state.resetPwd.data.key!==''" style="margin-top: 16px">
              <n-form-item path="code" label="邮箱验证码">
                <n-input v-model:value="state.resetPwd.data.code" placeholder="邮箱验证码"/>
              </n-form-item>
              <n-button style="width: 100%" type="primary" @click="resetPwdFunc" :loading="state.resetPwd.loading">
                重置密码
              </n-button>
            </div>
          </n-form>
        </n-tab-pane>
      </n-tabs>
      <n-space justify="end">
        <n-button text type="primary" @click="router.push({name:'Login'})">前往登录</n-button>
      </n-space>
      <n-alert type="info">未绑定邮箱，无法重置密码</n-alert>
    </n-card>
  </div>
</template>

<style lang="scss" scoped>
.n-divider:not(.n-divider--vertical) {
  margin-top: 10px;
  margin-bottom: 10px;
}

.resetPwd {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #303133;

  .bg {
    border-radius: 8px;
    overflow: hidden;
    padding: 16px;
  }
}

@media screen and (min-width: 700px) {
  .bg {
    width: 320px;
  }
}

@media screen and (max-width: 700px) {
  .bg {
    width: 320px;
  }
}

@media screen and (max-width: 500px) {
  .bg {
    width: 80%;
  }
}
</style>