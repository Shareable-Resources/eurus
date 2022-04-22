<template>
  <div id="PrivateChainStats">
    <div class="filter-form">
      <daily-report-export
        :daily-report-export-action="exportDailyReport"
      />
      <filter-form
        :filter-form="filterFormVal"
      />
    </div>
    <basic-table :table-data="tableData" />
  </div>
</template>

<script>
import { USER_GROUP } from '@/utils/selection.const'
export default {
  name: 'PrivateChainStats',
  data() {
    return {
      USER_GROUP: USER_GROUP,
      filterFormVal: {
        columns: [
          {
            type: 'datetime',
            model: 'startDate',
            label: 'common.startDate',
            format: 'yyyy/MM/dd hh:mm:ss'
          },
          {
            type: 'datetime',
            model: 'endDate',
            label: 'common.endDate',
            format: 'yyyy/MM/dd hh:mm:ss'
          },
          {
            type: 'select',
            model: 'walletType',
            label: 'stats.walletType',
            options: USER_GROUP
          }
        ],
        default: {
          startDate: '',
          endDate: ''
        },
        searchAction: this.searchAction,
        exportAction: this.exportAction
      },
      tableData: {
        source: [
          {
            type: 'centralized',
            sideChainGasFeeRevAmount: 6,
            sideChainGasFeePaidAmount: 1
          },
          {
            type: 'merchant',
            sideChainGasFeeRevAmount: 3,
            sideChainGasFeePaidAmount: 5
          }
        ],
        columns: [
          {
            type: 'options',
            label: 'stats.type',
            prop: 'type',
            options: USER_GROUP
          },
          {
            type: 'show',
            label: 'stats.sideChainGasFeeRevAmount',
            prop: 'sideChainGasFeeRevAmount'
          },
          {
            type: 'show',
            label: 'stats.sideChainGasFeePaidAmount',
            prop: 'sideChainGasFeePaidAmount'
          }
        ],
        currentPage: 1,
        total: 100,
        handleSizeChange: this.handleSizeChange,
        handleCurrentChange: this.handleCurrentChange,
        getSummaries: this.getSummaries
      }
    }
  },
  methods: {
    searchAction(data) {
      console.log(data, 'search')
    },
    exportAction(data) {
      console.log(data, 'export')
    },
    exportDailyReport(date) {
      console.log('daily report', date)
    },
    handleSizeChange(val) {
      console.log(`每页 ${val} 条`)
    },
    handleCurrentChange(val) {
      console.log(`当前页: ${val}`)
    },
    getSummaries(param) {
      const { columns, data } = param
      const sums = []
      const COUNT_MAP = {
        sideChainGasFeeRevAmount: true,
        sideChainGasFeePaidAmount: true
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
    }
  }
}
</script>
