<script setup>
import {resizeDirective} from '../../utils/resize.js'
import {computed, onBeforeMount, ref, watch} from "vue";
import {appStore} from "../../store/app.js";
import AppCard from "../../layout/components/AppCard.vue";
import Obs from "../../components/Obs.vue";
import {apiNormalSystemNoticeList} from "../../api/normal/system_notice.js";
import {localStore} from "../../store/local.js";
import Modal from "../../components/Modal.vue";
import {requiredRule} from "../../utils/formDataRule.js";
import {
  apiAuthBindEmail,
  apiAuthCloseOtp,
  apiAuthGenBindEmailCode,
  apiAuthGenOtp,
  apiAuthOpenOtp,
  apiAuthPassword,
  apiAuthUnBindEmail,
  apiAuthUserCheckin,
  apiAuthUserInfo
} from "../../api/auth/index.js";
import {NButton, NPopconfirm} from "naive-ui";
import {flowFormat} from "../../utils/flow.js";
import moment from "moment";
import {apiNormalGostObsUserMonth} from "../../api/normal/gost_obs.js";
import {
  apiNormalDashboardClientForwardObsDate,
  apiNormalDashboardClientHostObsDate,
  apiNormalDashboardClientObsDate,
  apiNormalDashboardClientTunnelObsDate,
  apiNormalDashboardCount
} from "../../api/normal/dashboard.js";
import Online from "../../icon/online.vue";
import DatePicker from "../../components/DatePicker.vue";

const vResize = resizeDirective

const state = ref({
  baseUrl: window.location.protocol + '//' + window.location.host,
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
  obsData: [],
  bindEmail: {
    open: false,
    loading: false,
    genLoading: false,
    data: {
      target: '',
      key: '',
      code: '',
    },
    dataRules: {
      target: requiredRule('请输入绑定的邮箱'),
      code: requiredRule('请输入验证码'),
    },
  },
  unBindEmailLoading: false,
  clientObsDate: [],
  clientHostObsDate: [],
  clientForwardObsDate: [],
  clientTunnelObsDate: [],
})

const cardStyleComputed = computed(() => {
  return {
    borderColor: 'var(--border-color)',
    borderRadius: 'var(--border-radius)',
    borderWidth: '1px',
    borderStyle: 'solid',
    // border: '1px solid ' + (localStore().darkTheme ? 'rgba(255, 255, 255, 0.09)' : 'rgb(239, 239, 245)')
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

const countFunc = async () => {
  try {
    let res = await apiNormalDashboardCount()
    state.value.gost = res.data
  } finally {

  }
}


const clientObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiNormalDashboardClientObsDate(date)
    state.value.clientObsDate = res.data || []
  } finally {

  }
}

const clientHostObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiNormalDashboardClientHostObsDate(date)
    state.value.clientHostObsDate = res.data || []
  } finally {

  }
}

const clientForwardObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiNormalDashboardClientForwardObsDate(date)
    state.value.clientForwardObsDate = res.data || []
  } finally {

  }
}

const clientTunnelObsDateFunc = async (date) => {
  try {
    if (date) {
      date = moment(date).format('yyyy-MM-DD')
    }
    let res = await apiNormalDashboardClientTunnelObsDate(date)
    state.value.clientTunnelObsDate = res.data || []
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

const openBindEmailModal = () => {
  state.value.bindEmail.data = {
    target: '',
    key: '',
    code: ''
  }
  state.value.bindEmail.open = true
}

const closeBindEmailModal = () => {
  state.value.bindEmail.open = false
}

const genBindEmailCodeFunc = () => {
  bindEmailRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.bindEmail.genLoading = true
        state.value.bindEmail.data.code = ''
        let res = await apiAuthGenBindEmailCode(state.value.bindEmail.data)
        state.value.bindEmail.data.key = res.data
      } finally {
        state.value.bindEmail.genLoading = false
      }
    }
  })
}

const bindEmailRef = ref()
const bindEmailFunc = () => {
  if (state.value.bindEmail.data.key === '') {
    $message.create('请先发送验证码', {
      type: "warning",
      closable: true,
      duration: 1500,
    })
    return
  }
  if (state.value.bindEmail.data.code === '') {
    $message.create('请输入验证码', {
      type: "warning",
      closable: true,
      duration: 1500,
    })
    return
  }
  bindEmailRef.value.validate(async valid => {
    if (!valid) {
      try {
        state.value.bindEmail.loading = true
        await apiAuthBindEmail(state.value.bindEmail.data)
        window.location.reload()
      } finally {
        state.value.bindEmail.loading = false
      }
    }
  })
}

const unBindEmailFunc = async () => {
  try {
    await apiAuthUnBindEmail()
    window.location.reload()
  } finally {

  }
}


onBeforeMount(() => {
  state.value.userInfo = appStore().userInfo
  noticeListFunc()
  countFunc()
  obsUserMonthFunc()
  let date = new Date()
  clientObsDateFunc(date)
  clientHostObsDateFunc(date)
  clientForwardObsDateFunc(date)
  clientTunnelObsDateFunc(date)
})

const alertSystemConfigBaseUrl = computed(() => {
  if (appStore().userInfo.admin === 1) {
    return appStore().siteConfig.baseUrl !== state.value.baseUrl;
  } else {
    return false
  }
})

const heightSync = ref(0)
const userInfoResize = (arg) => {
  heightSync.value = arg.height
}

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-grid :cols="gridColsComputed" :x-gap="12" :y-gap="12">
        <n-grid-item>
          <n-el tag="div" :style="cardStyleComputed">
            <AppCard :show-border="false" v-resize="userInfoResize">
              <n-h4 style="font-weight: bold">个人信息</n-h4>
              <n-descriptions :column="1" label-placement="left" label-class="userinfo-label">
                <n-descriptions-item label="账号">{{ state.userInfo.account }}</n-descriptions-item>
                <n-descriptions-item label="积分">{{ state.userInfo.amount }}</n-descriptions-item>
                <n-descriptions-item label="注册时间">{{ state.userInfo.createdAt }}</n-descriptions-item>

                <n-descriptions-item label="邮箱" v-if="state.userInfo.email!==''">
                  {{ state.userInfo.email }}
                  <n-popconfirm
                      :on-positive-click="unBindEmailFunc"
                      :positive-button-props="{loading:state.unBindEmailLoading}">
                    <template #trigger>
                      <n-button text type="info" :focusable="false" size="small">解绑</n-button>
                    </template>
                    确认解绑吗？
                  </n-popconfirm>
                </n-descriptions-item>
                <n-descriptions-item label="邮箱" v-if="state.userInfo.email===''">
                  <n-button text type="info" :focusable="false" size="small" @click="openBindEmailModal">绑定</n-button>
                </n-descriptions-item>
              </n-descriptions>

              <div v-if="appStore().siteConfig.checkIn==='1'">
                <p></p>
                <n-spin :show="state.checkinLoading">
                  <n-tag type="info" :bordered="false" v-if="state.userInfo.checkinAmount!=='0'">
                    今日已签到，获取{{ state.userInfo?.checkinAmount }}积分
                  </n-tag>
                  <n-alert type="info" v-else>
                    {{
                      `今日还未进行签到，签到可随机获得${appStore().siteConfig.checkInStart}-${appStore().siteConfig.checkInEnd}积分，`
                    }}
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
              <n-scrollbar :style="{maxHeight:heightSync-44+'px'}">
                <n-alert style="margin-bottom: 8px" type="info" v-for="(item,index) in state.notices" :key="item.code"
                         :bordered="false">
                  <template #header>
                    <div style="display: flex;justify-content: space-between;flex-direction: column;line-height: 1.5">
                      <span style="font-weight: bold">{{ item.title }}</span>
                      <span style="font-size: 0.8em">{{ item.date }}</span>
                    </div>
                  </template>
                  <span v-for="(content,contentIndex) in item.content.split('\n')" :key="contentIndex">
                    <br v-if="contentIndex!==0">
                    {{ content }}
                  </span>
                </n-alert>
              </n-scrollbar>
            </AppCard>
          </n-el>
        </n-grid-item>

        <n-grid-item>
          <n-el tag="div" :style="cardStyleComputed">
            <n-card size="small" :content-style="{padding:'6px 16px !important'}">
              <template #header>
                <n-space justify="space-between" align="center">
                  <span>客户端流量排行 ( IN | OUT )</span>
                  <DatePicker :max-value="new Date()" @on-before="args => clientObsDateFunc(args)"
                              @on-after="args => clientObsDateFunc(args)"></DatePicker>
                </n-space>
              </template>
              <n-scrollbar style="height: 300px">
                <n-list v-if="state.clientObsDate?.length">
                  <n-list-item v-for="(obs,index) in state.clientObsDate">
                    <div style="display: flex;justify-content: space-between;align-items: center">
                      <div style="width: 50px">{{ index + 1 }}</div>
                      <div style="flex: 1">
                        <Online :online="obs.online === 1"></Online>
                        {{ obs.name }}
                      </div>
                      <div>
                        <span>{{ flowFormat(obs.inputBytes) }}</span>
                        <n-divider vertical></n-divider>
                        <span>{{ flowFormat(obs.outputBytes) }}</span>
                      </div>
                    </div>
                  </n-list-item>
                </n-list>
                <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
              </n-scrollbar>
            </n-card>
          </n-el>
        </n-grid-item>

        <n-grid-item>
          <n-el tag="div" :style="cardStyleComputed">
            <n-card size="small" :content-style="{padding:'6px 16px !important'}">
              <template #header>
                <n-space justify="space-between" align="center">
                  <span>域名解析流量排行 ( IN | OUT )</span>
                  <DatePicker :max-value="new Date()" @on-before="args => clientHostObsDateFunc(args)"
                              @on-after="args => clientHostObsDateFunc(args)"></DatePicker>
                </n-space>
              </template>
              <n-scrollbar style="height: 300px">
                <n-list v-if="state.clientHostObsDate?.length">
                  <n-list-item v-for="(obs,index) in state.clientHostObsDate">
                    <div style="display: flex;justify-content: space-between;align-items: center">
                      <div style="width: 50px">{{ index + 1 }}</div>
                      <div style="flex: 1">
                        {{ obs.name }}
                      </div>
                      <div>
                        <span>{{ flowFormat(obs.inputBytes) }}</span>
                        <n-divider vertical></n-divider>
                        <span>{{ flowFormat(obs.outputBytes) }}</span>
                      </div>
                    </div>
                  </n-list-item>
                </n-list>
                <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
              </n-scrollbar>
            </n-card>
          </n-el>
        </n-grid-item>

        <n-grid-item>
          <n-el tag="div" :style="cardStyleComputed">
            <n-card size="small" :content-style="{padding:'6px 16px !important'}">
              <template #header>
                <n-space justify="space-between" align="center">
                  <span>端口转发流量排行 ( IN | OUT )</span>
                  <DatePicker :max-value="new Date()" @on-before="args => clientForwardObsDateFunc(args)"
                              @on-after="args => clientForwardObsDateFunc(args)"></DatePicker>
                </n-space>
              </template>
              <n-scrollbar style="height: 300px">
                <n-list v-if="state.clientForwardObsDate?.length">
                  <n-list-item v-for="(obs,index) in state.clientForwardObsDate">
                    <div style="display: flex;justify-content: space-between;align-items: center">
                      <div style="width: 50px">{{ index + 1 }}</div>
                      <div style="flex: 1">
                        {{ obs.name }}
                      </div>
                      <div>
                        <span>{{ flowFormat(obs.inputBytes) }}</span>
                        <n-divider vertical></n-divider>
                        <span>{{ flowFormat(obs.outputBytes) }}</span>
                      </div>
                    </div>
                  </n-list-item>
                </n-list>
                <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
              </n-scrollbar>
            </n-card>
          </n-el>
        </n-grid-item>

        <n-grid-item>
          <n-el tag="div" :style="cardStyleComputed">
            <n-card size="small" :content-style="{padding:'6px 16px !important'}">
              <template #header>
                <n-space justify="space-between" align="center">
                  <span>私有隧道流量排行 ( IN | OUT )</span>
                  <DatePicker :max-value="new Date()" @on-before="args => clientTunnelObsDateFunc(args)"
                              @on-after="args => clientTunnelObsDateFunc(args)"></DatePicker>
                </n-space>
              </template>
              <n-scrollbar style="height: 300px">
                <n-list v-if="state.clientTunnelObsDate?.length">
                  <n-list-item v-for="(obs,index) in state.clientTunnelObsDate">
                    <div style="display: flex;justify-content: space-between;align-items: center">
                      <div style="width: 50px">{{ index + 1 }}</div>
                      <div style="flex: 1">
                        {{ obs.name }}
                      </div>
                      <div>
                        <span>{{ flowFormat(obs.inputBytes) }}</span>
                        <n-divider vertical></n-divider>
                        <span>{{ flowFormat(obs.outputBytes) }}</span>
                      </div>
                    </div>
                  </n-list-item>
                </n-list>
                <n-empty v-else style="width: 100%;" description="暂无数据"></n-empty>
              </n-scrollbar>
            </n-card>
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
              <Obs :data="state.obsData" :loading="state.obsLoading" :dark="localStore().darkTheme"
                   style="width: 100%"></Obs>
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

    <Modal
        title="绑定邮箱"
        :show="state.bindEmail.open"
        confirm-text="绑定"
        cancel-text="取消"
        width="500px"
        @on-confirm="bindEmailFunc"
        @on-cancel="closeBindEmailModal"
        :confirm-loading="state.bindEmail.loading"
    >
      <n-form ref="bindEmailRef" :rules="state.bindEmail.dataRules" :model="state.bindEmail.data">
        <n-form-item label="邮箱" path="target">
          <n-input-group>
            <n-input v-model:value="state.bindEmail.data.target" placeholder="请输入邮箱"></n-input>
            <n-button :disabled="state.bindEmail.data.target===''" type="info" @click="genBindEmailCodeFunc"
                      :loading="state.bindEmail.genLoading">
              {{ state.bindEmail.data.key === '' ? '发送' : '重新发送' }}
            </n-button>
          </n-input-group>
        </n-form-item>

        <n-form-item label="验证码" path="code" v-if="state.bindEmail.data.key!==''">
          <n-input :disabled="state.bindEmail.data.key===''" v-model:value="state.bindEmail.data.code"
                   placeholder="请输入验证码"></n-input>
        </n-form-item>
      </n-form>
      <n-alert type="info">绑定邮箱后，如果忘记密码可以通过邮箱重置密码</n-alert>
    </Modal>
  </div>
</template>

<style scoped lang="scss">
:deep(.userinfo-label), :deep(.n-thing-header__title) {
  font-weight: bold !important;
  font-size: 14px !important;
}
:deep(.n-card){
  border: none;
}
</style>