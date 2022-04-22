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
      default: 'userSta'
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
          text: '用戶統計',
          left: 'center'
        },
        tooltip: {
          show: false,
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
        yAxis: {},
        dataset: {
          dimensions: ['date', 'data', 'data2', 'data3'],
          source: []
        },
        series: [
          {
            name: '新注冊用戶數量',
            type: 'line',
            label: {
              show: true
            }
          },
          {
            name: '通過KYC用戶數量',
            type: 'line',
            label: {
              show: true
            }
          },
          {
            name: '活躍用戶數量',
            type: 'line',
            label: {
              show: true
            }
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
