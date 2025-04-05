<script setup>
import {computed, onBeforeMount, ref, watch} from "vue";
import {appStore} from "../../store/app.js";
import AppCard from "../../layout/components/AppCard.vue";
import Obs from "../../components/Obs.vue";
import {apiNormalSystemNoticeList} from "../../api/normal/system_notice.js";
import {localStore} from "../../store/local.js";
import Modal from "../../components/Modal.vue";
import {requiredRule} from "../../utils/formDataRule.js";
import {
  apiAuthCloseOtp,
  apiAuthGenOtp,
  apiAuthOpenOtp,
  apiAuthPassword,
  apiAuthUserCheckin,
  apiAuthUserInfo
} from "../../api/auth/index.js";
import {NButton, NPopconfirm} from "naive-ui";
import {apiNormalGostInfo} from "../../api/normal/gost.js";
import {flowFormat} from "../../utils/flow.js";
import moment from "moment";
import {apiNormalGostObsUserMonth} from "../../api/normal/gost_obs.js";

const state = ref({
  userInfo: {},
  gost: {},
  notices: [],
  password: {
    data: {
      oldPwd: '',
      newPwd: '',
    },
    dataRules: {
      oldPwd: requiredRule('请输入原密码'),
      newPwd: requiredRule('请输入新密码'),
    },
    open: false,
    loading: false,
  },
  openOtp: {
    data: {
      key: '',
      value: '',
    },
    img: '',
    open: false,
    loading: false,
    genOtpLoading: false,
  },
  closeOtp: {
    loading: false,
  },
  checkinLoading: false,
  obsDataRange: 1,
  obsLoading: false,
  obsData: []
})

const cardStyleComputed = computed(() => {
  return {
    borderColor: 'var(--border-color)',
    borderRadius: 'var(--border-radius)',
    borderWidth: '1px',
    borderStyle: 'solid',
  }
})

const gridColsComputed = computed(() => {
  if (appStore().width > 768) {
    return 2
  }
  return 1
})

const noticeListFunc = async () => {
  try {
    let res = await apiNormalSystemNoticeList()
    state.value.notices = res.data || []
  } finally {

  }
}

const gostInfoFunc = async () => {
  try {
    let res = await apiNormalGostInfo()
    state.value.gost = res.data
  } finally {

  }
}

const logoutFunc = () => {
  localStore().auth.token = ''
  localStore().auth.expAt = ''
  window.location.reload()
}

const passwordRef = ref()
const passwordFunc = () => {
  passwordRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.password.loading = true
        await apiAuthPassword(state.value.password.data)
        closePwdModal()
        localStore().auth.expAt = ''
        localStore().auth.token = ''
        window.location.reload()
      } finally {
        state.value.password.loading = false
      }
    }
  })
}

const openPwdModal = () => {
  state.value.password.data = {
    oldPwd: '',
    newPwd: '',
  }
  state.value.password.open = true
}

const closePwdModal = () => {
  state.value.password.open = false
}

const genOtpFunc = async () => {
  try {
    state.value.openOtp.genOtpLoading = true
    let res = await apiAuthGenOtp()
    state.value.openOtp.data.key = res.data.key
    state.value.openOtp.img = res.data.img
  } finally {
    state.value.openOtp.genOtpLoading = false
  }
}

const otpImageComputed = computed(() => {
  return 'data:image/png;base64,' + state.value.openOtp.img
})

const openOtpFunc = async () => {
  if (!state.value.openOtp.data.value) {
    $message.create('请输入TOTP工具生成的数字', {
      type: "warning",
      closable: true,
      duration: 1500,
    })
  } else {
    try {
      state.value.openOtp.loading = true
      await apiAuthOpenOtp(state.value.openOtp.data)
      state.value.userInfo.otp = 1
      appStore().userInfo.otp = 1
      closeOtpModal()
    } catch (err) {
      await genOtpFunc()
    } finally {
      state.value.openOtp.loading = false
    }
  }
}

const closeOtpFunc = async () => {
  try {
    state.value.closeOtp.loading = true
    await apiAuthCloseOtp()
    state.value.userInfo.otp = 2
    appStore().userInfo.otp = 2
  } finally {
    state.value.closeOtp.loading = false
  }
}

const openOtpModal = () => {
  genOtpFunc()
  state.value.openOtp.data.value = ''
  state.value.openOtp.open = true
}

const closeOtpModal = () => {
  state.value.openOtp.open = false
}

const checkinFunc = async () => {
  try {
    state.value.checkinLoading = true
    await apiAuthUserCheckin()
    let res = await apiAuthUserInfo()
    state.value.userInfo = res.data
    appStore().userInfo = res.data
  } finally {
    state.value.checkinLoading = false
  }
}

watch(() => ({type: state.value.obsDataRange}), () => {
  obsUserMonthFunc()
})

const obsUserMonthFunc = async () => {
  try {
    state.value.obsLoading = true
    let data = {
      start: moment().add(-29, 'days').format('yyyy-MM-DD'),
      end: moment().format('yyyy-MM-DD')
    }
    if (state.value.obsDataRange === 1) {
      data = {
        start: moment().add(-6, 'days').format('yyyy-MM-DD'),
        end: moment().format('yyyy-MM-DD')
      }
    }
    let res = await apiNormalGostObsUserMonth(data)
    state.value.obsData = res.data || []
  } finally {
    state.value.obsLoading = false
  }
}

onBeforeMount(() => {
  state.value.userInfo = appStore().userInfo
  noticeListFunc()
  gostInfoFunc()
  obsUserMonthFunc()
})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-grid :cols="gridColsComputed" :x-gap="12" :y-gap="12">
        <n-grid-item>
          <n-el tag="div" :style="cardStyleComputed">
            <AppCard :show-border="false">
              <n-h4 style="font-weight: bold">个人信息</n-h4>
              <n-descriptions :column="1" label-placement="left" label-class="userinfo-label">
                <n-descriptions-item label="账号">{{ state.userInfo.account }}</n-descriptions-item>
                <n-descriptions-item label="积分">{{ state.userInfo.amount }}</n-descriptions-item>
                <n-descriptions-item label="注册时间">{{ state.userInfo.createdAt }}</n-descriptions-item>
              </n-descriptions>
              <div v-if="appStore().siteConfig.checkIn==='1'">
                <p></p>
                <n-spin :show="state.checkinLoading">
                  <n-tag type="info" :bordered="false" v-if="state.userInfo.checkinAmount!=='0'">
                    今日已签到，获取{{ state.userInfo?.checkinAmount }}积分
                  </n-tag>
                  <n-alert type="info" v-else>
                    今日还未进行签到，签到可随机获取1-5积分，
                    <n-button
                        text
                        type="info"
                        size="small"
                        :focusable="false"
                        @click="checkinFunc">
                      点击
                    </n-button>
                    进行签到
                  </n-alert>
                </n-spin>
              </div>
              <p></p>
              <n-grid :cols="2">
                <n-grid-item>
                  <n-thing title="客户端">{{ state.gost.client }}</n-thing>
                </n-grid-item>
                <n-grid-item>
                  <n-thing title="域名解析">{{ state.gost.host }}</n-thing>
                </n-grid-item>
                <n-grid-item>
                  <n-thing title="端口转发">{{ state.gost.forward }}</n-thing>
                </n-grid-item>
                <n-grid-item>
                  <n-thing title="私有隧道">{{ state.gost.tunnel }}</n-thing>
                </n-grid-item>
                <n-grid-item>
                  <n-thing title="代理隧道">{{ state.gost.proxy }}</n-thing>
                </n-grid-item>
                <n-grid-item>
                  <n-thing title="P2P隧道">{{ state.gost.p2p }}</n-thing>
                </n-grid-item>
                <n-grid-item>
                  <n-thing title="下行流量">{{ flowFormat(state.gost.inputBytes) }}</n-thing>
                </n-grid-item>
                <n-grid-item>
                  <n-thing title="上行流量">{{ flowFormat(state.gost.outputBytes) }}</n-thing>
                </n-grid-item>
              </n-grid>
              <n-space style="margin-top: 8px">
                <n-button text type="info" @click="openPwdModal">修改密码</n-button>
                <n-button text type="warning" v-if="state.userInfo.otp===2" @click="openOtpModal">
                  二步验证(未开启)
                </n-button>
                <n-popconfirm
                    v-else
                    @positive-click="closeOtpFunc"
                    :positive-button-props="{loading:state.closeOtp.loading}"
                >
                  <template #trigger>
                    <n-button text type="success">二步验证(已开启)</n-button>
                  </template>
                  确认关闭吗？
                </n-popconfirm>
                <n-popconfirm @positive-click="logoutFunc">
                  <template #trigger>
                    <n-button text type="error">退出登录</n-button>
                  </template>
                  确认退出吗？
                </n-popconfirm>
              </n-space>
            </AppCard>
          </n-el>
        </n-grid-item>
        <n-grid-item>
          <n-el tag="div" :style="cardStyleComputed">
            <AppCard :show-border="false">
              <n-h4 style="font-weight: bold">通知公告</n-h4>
              <n-scrollbar style="max-height: 600px">
                <n-alert style="margin-bottom: 8px" type="info" v-for="item in state.notices" :key="item.code"
                         :bordered="false">
                  <template #header>
                    <div style="display: flex;justify-content: space-between;flex-direction: column;line-height: 1.5">
                      <span style="font-weight: bold">{{ item.title }}</span>
                      <span style="font-size: 0.8em">{{ item.date }}</span>
                    </div>
                  </template>
                  {{ item.content }}
                </n-alert>
              </n-scrollbar>
            </AppCard>
          </n-el>
        </n-grid-item>
        <n-grid-item :span="gridColsComputed">
          <n-el tag="div" :style="cardStyleComputed">
            <AppCard :show-border="false" style="min-height: 300px">
              <n-space justify="space-between">
                <n-h4 style="font-weight: bold">最近{{ state.obsDataRange === 1 ? '7' : '30' }}天流量使用趋势</n-h4>
                <n-radio-group size="small" v-model:value="state.obsDataRange">
                  <n-radio-button :value="1">最近7天</n-radio-button>
                  <n-radio-button :value="2">最近30天</n-radio-button>
                </n-radio-group>
              </n-space>
              <Obs :data="state.obsData" :loading="state.obsLoading" :dark="localStore().darkTheme" style="width: 100%"></Obs>
            </AppCard>
          </n-el>
        </n-grid-item>
      </n-grid>
    </AppCard>

    <Modal
        :show="state.password.open"
        title="修改密码"
        width="400px"
        confirm-text="保存"
        cancel-text="取消"
        :confirm-loading="state.password.loading"
        @on-confirm="passwordFunc"
        @on-cancel="closePwdModal"
        mask-close
    >
      <n-form ref="passwordRef" :rules="state.password.dataRules" :model="state.password.data">
        <n-form-item path="oldPwd" label="原密码">
          <n-input v-model:value.trim="state.password.data.oldPwd" type="password" show-password-on="click"></n-input>
        </n-form-item>
        <n-form-item path="newPwd" label="新密码">
          <n-input v-model:value.trim="state.password.data.newPwd" type="password" show-password-on="click"></n-input>
        </n-form-item>
      </n-form>
    </Modal>

    <Modal
        title="二步验证"
        :show="state.openOtp.open"
        confirm-text="绑定"
        width="500px"
        @on-cancel="closeOtpModal"
        @on-confirm="openOtpFunc"
        :confirm-loading="state.openOtp.loading"
    >
      <n-grid :x-gap="12" :y-gap="12" :cols="2">
        <n-grid-item :span="6">
          <div style="display: flex;align-items: center;justify-content: center">
            <img alt="" style="width: 100%;max-width:200px;object-fit: fill" :src="otpImageComputed">
          </div>
        </n-grid-item>
        <n-grid-item :span="18">
          <span>请使用任意二步验证APP或者支持二步验证的密码管理软件扫描左侧二维码添加本站。扫描完成后请填写二步验证APP给出的6位验证码以开启二步验证。</span>
          <br>
          <br>
          <n-input placeholder="TOTP验证码" v-model:value="state.openOtp.data.value"></n-input>
        </n-grid-item>
      </n-grid>
    </Modal>
  </div>
</template>

<style scoped lang="scss">
:deep(.userinfo-label), :deep(.n-thing-header__title) {
  font-weight: bold !important;
  font-size: 14px !important;
}

:deep(.n-grid) {
  & > div > div:first-child {
    //height: 100% !important;
  }
}

</style>