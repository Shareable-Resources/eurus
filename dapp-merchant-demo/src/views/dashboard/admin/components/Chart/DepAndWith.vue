<template>
  <div
    :id="id"
    :class="className"
    :style="{height:height,width:width}"
  />
</template>
<script>
import echarts from 'echarts'
import resize from './../mixins/resize'
export default {
  mixins: [resize],
  props: {
    loading: {
      type: Boolean,
      require: true,
      default: false
    },
    options: {
      type: Object,
      require: true,
      default: () => {}
    },
    id: {
      type: String,
      require: true,
      default: 'depAndWith'
    },
    className: {
      type: String,
      default: 'line'
    }
  },
  data() {
    return {
      chart: null,
      width: '1000px',
      height: '500px',
      dafaultOption: {
        title: {
          text: '用戶充值及提現',
          left: 'center'
        },
        tooltip: {
          show: true,
          trigger: 'axis'
        },
        legend: {
          show: true,
          bottom: '15px'
        },
        xAxis: {
          data: [],
          type: 'category'
        },
        yAxis: [
          {
            type: 'value',
            name: 'USDT',
            position: 'left'
          },
          {
            type: 'value',
            name: 'USDT',
            position: 'right'
          }
        ],
        dataset: {
          dimensions: [],
          source: []
        },
        series: [
          {
            name: '充值總數',
            type: 'bar',
            stack: 'deposit',
            barGap: 0
          },
          {
            name: '提現總數',
            type: 'bar',
            stack: 'withdraw'
          },
          {
            name: '提現手續費',
            type: 'line',
            stack: 'gas',
            yAxisIndex: 1
          }
        ]
      }
    }
  },
  watch: {
    options: {
      handler(val) {
        this.updateChart()
      },
      deep: true
    },
    loading: {
      handler(val) {
        if (val) {
          this.chart.showLoading()
        } else {
          this.chart.hideLoading()
        }
      }
    }
  },
  mounted() {
    this.initChart()
  },
  beforeDestroy() {
    if (!this.chart) {
      return
    }
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    updateChart() {
      this.chart.setOption(this.options)
    },
    initChart() {
      this.chart = echarts.init(document.getElementById(this.id))
      this.chart.setOption(this.dafaultOption)
    }
  }
}
</script>
