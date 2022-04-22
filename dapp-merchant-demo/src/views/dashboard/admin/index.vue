<template>
  <div class="dashboard-editor-container">
    <el-form
      ref="form"
      :inline="true"
      :model="form"
    >
      <el-form-item label="APPID">
        <el-select v-model="form.appid" placeholder="">
          <el-option
            v-for="(label, val) in APPID_OPTIONS"
            :key="val"
            :label="label"
            :value="val"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="顯示時間">
        <el-select v-model="form.timePeriod" placeholder="">
          <el-option
            v-for="(label, val) in TIME_PERIOD"
            :key="val"
            :label="label"
            :value="val"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          @click="onSearch"
        >{{ $t(`common.search`) }}
        </el-button>
      </el-form-item>
    </el-form>
    <div class="chart-group">
      <accum-user
        id="test"
        :options="barOptions"
        :loading="chartLoading"
      />
      <el-divider />
      <user-sta
        id="test2"
        :options="barOptions"
        :loading="chartLoading"
      />
      <el-divider />
      <dep-and-with
        id="test3"
        :options="barOptions"
        :loading="chartLoading"
      />
      <el-divider />
      <pay-and-trans
        id="test4"
        :options="barOptions"
        :loading="chartLoading"
      />
    </div>
  </div>
</template>

<script>
import AccumUser from './components/Chart/AccumUser'
import UserSta from './components/Chart/UserSta'
import DepAndWith from './components/Chart/DepAndWith'
import PayAndTrans from './components/Chart/PayAndTrans'
export default {
  name: 'DashboardAdmin',
  components: {
    AccumUser,
    UserSta,
    DepAndWith,
    PayAndTrans
  },
  data() {
    return {
      form: {
        appid: 'Eurus',
        timePeriod: 'month'
      },
      APPID_OPTIONS: {
        eurus: 'Eurus',
        merchant1: 'Merchant1',
        merchant2: 'Merchant2'
      },
      TIME_PERIOD: {
        month: '最近 30 天',
        twoMonth: '最近 60 天',
        year: '全年'
      },
      barOptions: {},
      chartLoading: false
    }
  },
  mounted() {
    this.onSearch()
  },
  methods: {
    onSearch() {
      this.chartLoading = true
      const { timePeriod } = this._clone(this.form)
      const map = {
        month: 30,
        twoMonth: 60,
        year: 12
      }
      const { date, source } = this.getData(map[timePeriod])
      setTimeout(() => {
        this.chartLoading = false
        this.$set(this, 'barOptions', {
          xAxis: {
            data: date,
            type: 'category'
          },
          dataset: {
            dimensions: ['date', 'data', 'data2', 'data3'],
            source: source
          }
        })
      }, 3000)
    },
    getData(inp) {
      const date = []
      const source = []
      for (let i = 0; i < inp; i++) {
        date.push(`01-${i + 1}`)
        source.push({
          date: `01-${i + 1}`,
          data: Math.round(Math.random() * 8 + i),
          data2: Math.round(Math.random() * 5 + i),
          data3: Math.round(Math.random() * 3 + i / 2)
        })
      }
      return {
        date,
        source
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}

@media (max-width:1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
