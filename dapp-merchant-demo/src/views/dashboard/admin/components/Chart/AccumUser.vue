
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
      default: 'accumUser'
    },
    className: {
      type: String,
      default: 'bar'
    }
  },
  data() {
    return {
      chart: null,
      width: '1000px',
      height: '500px',
      dafaultOption: {
        title: {
          text: this.$t(`dashboard.cumNoUser`),
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
          dimensions: [],
          source: []
        },
        series: [
          {
            name: 'decentralized',
            type: 'bar',
            stack: 'user',
            label: {
              show: true
            }
          },
          {
            name: 'centrailzed',
            type: 'bar',
            stack: 'user',
            label: {
              show: true
            }
          },
          {
            name: 'merchant',
            type: 'bar',
            stack: 'user',
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
