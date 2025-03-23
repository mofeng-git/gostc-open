<script setup>
import highcharts from 'highcharts'
import {onUnmounted, onUpdated, ref} from "vue";
import AppCard from "../layout/components/AppCard.vue";

const props = defineProps({
  loading: {
    default: false,
    type: Boolean,
  },
  dark:{
    default: false,
    type: Boolean,
  },
  data: {
    required: false,
    type: Array,
    default: [
      {
        data: '2024-01-01',
        in: 1024,
        out: 2048,
      },
      {
        data: '2024-01-02',
        in: 2048,
        out: 1024,
      },
      {
        data: '2024-01-03',
        in: 1024,
        out: 2048,
      },
      {
        data: '2024-01-04',
        in: 2048,
        out: 1024,
      },
    ]
  }
})

const chartRef = ref()
const chart = ref()

onUpdated(() => {
  let dates = props.data.map(item => item?.date)
  let ins = props.data.map(item => (item?.in / (1024 * 1024)))
  let outs = props.data.map(item => (item?.out / (1024 * 1024)))
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
      text: ''
    },
    xAxis: {
      categories: dates,
      labels: {
        style: {
          color: props.dark ?'rgb(240, 240, 240)':'rgb(51, 51, 51)'
        }
      }
    },
    yAxis: {
      title: {
        text: null,
        // text: '流量(MB)' // Y轴标题
      },
      labels: {
        style: {
          color: props.dark ?'rgb(240, 240, 240)':'rgb(51, 51, 51)'
        },
        formatter: function() {
          return this.value +'M'
        }
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
    plotOptions: {
      series: {
        boostThreshold: 1000,
        marker: {
          enabled: false,
          states: {
            hover: {
              enabled: true
            }
          }
        }
      }
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
})

onUnmounted(() => {
})
</script>

<template>
  <AppCard :show-border="false" :loading="props.loading" style="margin: 12px 0 !important">
    <div ref="chartRef"></div>
  </AppCard>
</template>

<style scoped lang="scss">

</style>