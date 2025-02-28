<script setup>
import {onBeforeMount, ref} from "vue";
import {apiAuthCaptcha, apiAuthLogin, apiAuthLoginOtp, apiAuthRegister} from "../../api/auth/index.js";
import {localStore} from "../../store/local.js";
import {appStore} from "../../store/app.js";
import {requiredRule} from "../../utils/formDataRule.js";

const state = ref({
  tabsValue: 'login',
  captcha: {
    src: '',
    loading: false,
    security: true,
  },
  login: {
    loading: false,
    data: {
      account: '',
      password: '',
      captchaKey: '',
      captchaValue: ''
    },
    dataRules: {
      account: requiredRule('请输入账号'),
      password: requiredRule('请输入密码'),
      captchaValue: requiredRule('请输入验证码'),
    }
  },
  register: {
    loading: false,
    data: {
      account: '',
      password: '',
      captchaKey: '',
      captchaValue: '',
    },
    dataRules: {
      account: requiredRule('请输入账号'),
      password: requiredRule('请输入密码'),
      captchaValue: requiredRule('请输入验证码'),
    }
  },
  loginOtp: {
    data: {
      key: '',
      value: '',
    },
    loading: false,
    dataRules: {
      value: requiredRule('请输入验证码'),
    }
  }
})

const refreshCaptchaFunc = async () => {
  try {
    state.value.captcha.loading = true
    let res = await apiAuthCaptcha()
    state.value.login.data.captchaKey = res.data.key
    state.value.register.data.captchaKey = res.data.key
    state.value.captcha.src = res.data?.bs64
    state.value.captcha.security = res.data?.security
  } finally {
    state.value.captcha.loading = false
  }
}

const otpRef = ref()
const loginOtpFunc = () => {
  otpRef.value.validate(async (valid) => {
    if (!valid) {
      try {
        state.value.loginOtp.loading = true
        let res = await apiAuthLoginOtp(state.value.loginOtp.data)
        localStore().auth.token = res.data.token
        localStore().auth.tokenExpAt = res.data.tokenExpAt
        window.location.reload()
      } finally {
        state.value.loginOtp.loading = false
      }
    }
  });
}

const loginRef = ref()
const loginFunc = () => {
  loginRef.value.validate(async (errors) => {
    if (!errors) {
      try {
        state.value.login.loading = true
        let res = await apiAuthLogin(state.value.login.data)
        if (res.data.otp === 1) {
          state.value.loginOtp.data.key = res.data.token
        } else {
          localStore().auth.token = res.data.token
          localStore().auth.tokenExpAt = res.data.tokenExpAt
          window.location.reload()
        }
      } catch (err) {
        await refreshCaptchaFunc()
      } finally {
        state.value.login.loading = false
      }
    }
  });
}

const registerRef = ref()
const registerFunc = () => {
  registerRef.value.validate(async (valid) => {
    if (!valid) {
      try {
        state.value.register.loading = true
        await apiAuthRegister(state.value.register.data)
        $message.create('注册成功', {
          type: "success",
          showIcon: false,
          duration: 3000
        })
        state.value.login.data.account = state.value.register.data.account
        state.value.login.data.password = state.value.register.data.password
        state.value.tabsValue = 'login'
        state.value.register.data.account = ''
        state.value.register.data.password = ''
      } finally {
        state.value.register.loading = false
      }
    }
  });
}

onBeforeMount(() => {
  refreshCaptchaFunc()
})

</script>

<template>
  <div class="login">
    <n-card class="bg">
      <div class="title" style="user-select: none;cursor: pointer;font-size: 1.5em">
        {{ appStore().siteConfig.title }}
      </div>
      <n-divider/>
      <n-tabs animated v-model:value="state.tabsValue">
        <n-tab-pane name="login" tab="登录">
          <n-form
              v-if="state.loginOtp.data.key===''"
              ref="loginRef"
              :model="state.login.data"
              :rules="state.login.dataRules"
              :show-label="false"
              size="large"
          >
            <n-form-item path="account" label="账号">
              <n-input
                  v-model:value="state.login.data.account"
                  placeholder="账号"
              />
            </n-form-item>
            <n-form-item path="password" label="密码">
              <n-input
                  v-model:value="state.login.data.password"
                  type="password"
                  placeholder="密码"
                  show-password-on="click"
              />
            </n-form-item>
            <n-form-item path="captchaValue" label="验证码" v-if="!state.captcha.security">
              <n-input
                  v-model:value="state.login.data.captchaValue"
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
            <n-button style="width: 100%" type="primary" @click="loginFunc" :loading="state.login.loading">
              登录
            </n-button>
            <!--            <n-space style="margin-top: 8px">-->
            <!--              <n-checkbox v-model:checked="state.rememberAccountFlag" :focusable="false">记住</n-checkbox>-->
            <!--            </n-space>-->
          </n-form>
          <n-form
              v-else
              ref="otpRef"
              :model="state.loginOtp.data"
              :rules="state.loginOtp.dataRules"
              :show-label="false"
              size="large"
          >
            <n-form-item path="value" label="验证码">
              <n-input
                  v-model:value="state.loginOtp.data.value"
                  placeholder="请输入TOTP验证码"
              />
            </n-form-item>
            <n-button style="width: 100%" type="primary" @click="loginOtpFunc" :loading="state.loginOtp.loading">
              验证登录
            </n-button>
          </n-form>
        </n-tab-pane>

        <n-tab-pane name="register" tab="注册">
          <n-form
              ref="registerRef"
              :model="state.register.data"
              :rules="state.register.dataRules"
              :show-label="false"
              size="large"
          >
            <n-form-item path="account" label="账号">
              <n-input
                  v-model:value="state.register.data.account"
                  placeholder="账号"
              />
            </n-form-item>
            <n-form-item path="password" label="密码">
              <n-input
                  v-model:value="state.register.data.password"
                  type="password"
                  placeholder="密码"
                  show-password-on="click"
              />
            </n-form-item>
            <n-form-item path="captchaValue" label="验证码" v-if="!state.captcha.security">
              <n-input
                  v-model:value="state.register.data.captchaValue"
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
            <n-button style="width: 100%" type="primary" @click="registerFunc" :loading="state.register.loading">
              注册
            </n-button>
          </n-form>
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<style lang="scss" scoped>
.n-divider:not(.n-divider--vertical) {
  margin-top: 10px;
  margin-bottom: 10px;
}

.login {
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