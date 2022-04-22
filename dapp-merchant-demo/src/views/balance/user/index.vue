<template>
  <div id="user-balance">
    <div class="filter-form">
      <real-time
        :real-time-refresh-action="refreshRealTime"
      />
      <daily-report-export
        :daily-report-export-action="exportDailyReport"
      />
      <filter-form
        :filter-form="filterFormVal"
      />
    </div>

    <div class="table">
      <el-table
        border
        :data="tableData"
        :summary-method="getSummaries"
        show-summary
        style="width: 100%"
      >
        <el-table-column type="expand">
          <template slot-scope="props">
            <el-table
              :data="props.row.sub"
            >
              <el-table-column
                :label="$t(`common.userGroup`)"
                prop="group"
              />
              <!-- <el-table-column
                label="Asset"
                prop="asset">
              </el-table-column> -->
              <el-table-column
                :label="$t(`common.available`)"
                prop="bal"
              />
              <el-table-column
                :label="$t(`common.frozen`)"
                prop="frozen"
              />
              <el-table-column
                :label="$t(`common.total`)"
                prop="total"
              />
              <el-table-column
                :label="$t(`common.totalUSDT`)"
                prop="usdt_total"
              />
            </el-table>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t(`common.asset`)"
          prop="asset"
        />
        <el-table-column
          :label="$t(`common.available`)"
          prop="bal"
        />
        <el-table-column
          :label="$t(`common.frozen`)"
          prop="frozen"
        />
        <el-table-column
          :label="$t(`common.total`)"
          prop="total"
        />
        <el-table-column
          :label="$t(`common.totalUSDT`)"
          prop="usdt_total"
        />
      </el-table>
    </div>
  </div>
</template>

<script>
import { USER_GROUP } from '@/utils/selection.const'
export default {
  name: 'UserBalance',
  data() {
    return {
      filterFormVal: {
        columns: [
          // {
          //   type: 'input',
          //   model: 'testInp',
          //   label: 'filter.snapShot',
          //   placeholder: ''
          // // },
          // {
          //   type: 'datetime',
          //   model: 'snapShot',
          //   label: 'filter.snapShot',
          //   format: 'yyyy/MM/dd HH:mm:ss'
          // },
          {
            type: 'snapShot',
            model: 'snapShot',
            label: 'filter.snapShot'
          },
          {
            type: 'select',
            model: 'group',
            label: 'common.userGroup',
            options: USER_GROUP
          }
        ],
        default: {
          testInp: 'testinp',
          snapShot: '',
          group: 'merchant'
        },
        // addAction: (data)=>{console.log('add',data)},
        searchAction: this.searchAction,
        exportAction: this.exportAction
      },
      tableData: [
        {
          asset: 'BTC',
          bal: '3',
          frozen: '0.657',
          total: '3.657',
          usdt_total: '52001',
          sub: [{
            asset: 'BTC',
            group: 'centralized',
            bal: '1',
            frozen: '0.3445',
            total: '1.3445',
            usdt_total: '16000'
          }, {
            asset: 'BTC',
            group: 'decentralized',
            bal: '2',
            frozen: '0.3125',
            total: '2.3125',
            usdt_total: '36001'
          }]
        },
        {
          asset: 'ETH',
          bal: '3',
          frozen: '0.657',
          total: '3.657',
          usdt_total: '52001',
          sub: [{
            asset: 'ETH',
            group: 'centralized',
            bal: '1',
            frozen: '0.3445',
            total: '1.3445',
            usdt_total: '16000'
          }, {
            asset: 'ETH',
            group: 'centralized',
            bal: '2',
            frozen: '0.3125',
            total: '2.3125',
            usdt_total: '36001'
          }]
        }
      ]
    }
  },
  methods: {
    refreshRealTime(timestamp) {
      console.log(timestamp)
    },
    exportDailyReport(date) {
      console.log('daily report', date)
    },
    getSummaries(param) {
      const { columns, data } = param
      const sums = []
      const COUNT_MAP = {
        usdt_total: true
      }
      columns.forEach((column, index) => {
        if (COUNT_MAP[column.property]) {
          const values = data.map(item => Number(item[column.property]))
          if (!values.every(value => isNaN(value))) {
            sums[index] = values.reduce((prev, curr) => {
              const value = Number(curr)
              if (!isNaN(value)) {
                return prev + curr
              } else {
                return prev
              }
            }, 0)
            sums[index] += ''
          } else {
            sums[index] = 'N/A'
          }
        }
      })
      return sums
    },
    searchAction(data) {
      console.log('search', data)
    },
    exportAction(data) {
      console.log('export', data)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
