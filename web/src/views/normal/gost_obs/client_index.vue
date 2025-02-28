<script setup>
import highcharts from 'highcharts'
import {onMounted, onUnmounted, ref, watch} from "vue";
import AppCard from "../../../layout/components/AppCard.vue";
import router from "../../../router/index.js";
import {apiNormalGostObsClientMonth} from "../../../api/normal/gost_obs.js";
import {apiNormalGostClientList} from "../../../api/normal/gost_client.js";
import {NSpace} from "naive-ui";
import Online from "../../../icon/online.vue";

const state = ref({
  code: '',
  data: [],
  loading: false,
  clients: [],
})


const chartRef = ref()
const chart = ref()

const obsMonthFunc = async () => {
  try {
    let res = await apiNormalGostObsClientMonth({code: state.value.code})
    state.value.data = res.data || []
    let dates = state.value.data.map(item => item?.date)
    let ins = state.value.data.map(item => (item?.in / (1024 * 1024)))
    let outs = state.value.data.map(item => (item?.out / (1024 * 1024)))
    const options = {
      chart: {
        type: "area",
        style: {
          fontSize: "12px",
          fontWeight: "bold",
          color: "white"
        },
        borderWidth: 0,
        backgroundColor: "#00000000"
      },
      accessibility: {
        enabled: false,
      },
      title: {
        text: '最近30天流量使用趋势' // 图表标题
      },
      xAxis: {
        categories: dates,
      },
      yAxis: {
        title: {
          text: '流量(MB)' // Y轴标题
        }
      },
      tooltip: {
        formatter: function () {
          return `<div>上行流量：${outs[this.x].toFixed(2)} MB</div><br><div>下行流量：${ins[this.x].toFixed(2)} MB</div>`;
        },
        backgroundColor: '#fff',
        borderColor: '#fafafa',
        borderRadius: 6,
        borderWidth: 1
      },
      series: [
        {
          name: 'In',
          color: '#db9145',
          fillColor: '#db914520',
          data: ins,
          showInLegend: false,
        },
        {
          name: 'Out',
          color: '#a45bd4',
          fillColor: '#a45bd420',
          data: outs,
          showInLegend: false,
        },
      ],
      credits: {
        enabled: false
      }
    }
    highcharts.setOptions({global: {useUTC: false}});
    options.chart.renderTo = chartRef.value
    chart.value = highcharts.chart(options);
  } catch (e) {
    console.log(e)
    router.back()
  } finally {
  }
}

const clientListFunc = async () => {
  try {
    let res = await apiNormalGostClientList()
    state.value.clients = res.data || []
    if (state.value.clients.length > 0) {
      state.value.code = state.value.clients[0].code
    } else {
      router.back()
      $message.create('暂无客户端信息', {
        type: "warning",
        closable: true,
        duration: 1500,
      })
    }
  } finally {

  }
}

onMounted(async () => {
  await clientListFunc()
  let code = router.currentRoute.value.query?.code
  if (!code) {
    state.value.code = state.value.clients[0].code
  } else {
    state.value.code = code
  }
  await obsMonthFunc()
})

watch(() => ({code: state.value.code}), (val, oldVal) => {
  if (oldVal.code !== '') {
    router.replace({name: 'NormalGostClientObs', query: {code: val.code}})
  }
})

onUnmounted(() => {
})
</script>

<template>
  <AppCard :show-border="false">
    <n-h4 style="font-weight: bold">客户端：</n-h4>
    <n-alert v-if="state.clients.length===0" type="warning">没有客户端，请先新增一个客户端</n-alert>
    <n-radio-group v-else v-model:value="state.code"
                   style="width: 100%">
      <n-grid x-gap="12" y-gap="12" cols="550:2 800:3 1400:4 1">
        <n-grid-item v-for="client in state.clients">
          <n-alert
              type="info"
              :show-icon="false"
              :bordered="false"
              style="height: 100%;cursor: pointer"
              @click="state.code = client.code"
          >
            <n-radio
                :key="client.code"
                :value="client.code"
                style="width: 100%;"
            >
              <n-space justify="space-between" style="width: 100%">
                <Online :online="client.online===1"></Online>
                <span>{{ client.name }}</span>
              </n-space>
            </n-radio>
          </n-alert>
        </n-grid-item>
      </n-grid>
    </n-radio-group>
    <div ref="chartRef"></div>
  </AppCard>
</template>

<style scoped lang="scss">

</style>