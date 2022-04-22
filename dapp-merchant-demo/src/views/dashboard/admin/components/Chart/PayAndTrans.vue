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
      default: 'payAndTrans'
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
          text: '用戶支付及轉賬',
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
            name: 'EUN',
            position: 'right'
          }
        ],
        dataset: {
          dimensions: [],
          source: []
        },
        series: [
          {
            name: '支付總數',
            type: 'line',
            stack: 'deposit',
            barGap: 0
          },
          {
            name: '轉賬總數',
            type: 'line',
            stack: 'withdraw'
          },
          {
            name: '私鏈支出gas fee',
            type: 'bar',
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
