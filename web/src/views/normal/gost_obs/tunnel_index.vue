<script setup>
import highcharts from 'highcharts'
import {onMounted, onUnmounted, ref} from "vue";
import AppCard from "../../../layout/components/AppCard.vue";
import router from "../../../router/index.js";
import {apiNormalGostObsTunnelMonth} from "../../../api/normal/gost_obs.js";

const state = ref({
  code: '',
  data: [],
  loading: false,
})


const chartRef = ref()
const chart = ref()

const obsMonthFunc = async () => {
  try {
    state.value.loading = true
    let res = await apiNormalGostObsTunnelMonth({code: state.value.code})
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
        // margin: [0, -8, 0, -8],
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
    router.back()
  } finally {
    state.value.loading = false
  }
}


onMounted(() => {
  state.value.code = router.currentRoute.value.query?.code
  if (state.value.code) {
    obsMonthFunc()
  } else {
    router.back()
  }
})

onUnmounted(() => {
})
</script>

<template>
  <AppCard :show-border="false" :loading="state.loading">
    <div ref="chartRef"></div>
  </AppCard>
</template>

<style scoped lang="scss">

</style>