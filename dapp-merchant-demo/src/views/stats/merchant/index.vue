<template>
  <div id="MerchantStats">
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
import { MERCHANT_BUSINESSTYPE } from '@/utils/selection.const'
export default {
  name: 'MerchantStats',
  data() {
    return {
      MERCHANT_BUSINESSTYPE: MERCHANT_BUSINESSTYPE,
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
            businessType: 'restaurant',
            newMerchantDuringPeriodCount: 6,
            cumMerchantCount: 1
          },
          {
            businessType: 'others',
            newMerchantDuringPeriodCount: 4,
            cumMerchantCount: 2
          }
        ],
        columns: [
          {
            type: 'options',
            label: 'merchant.businessType',
            prop: 'businessType',
            options: MERCHANT_BUSINESSTYPE
          },
          {
            type: 'show',
            label: 'stats.newMerchantDuringPeriodCount',
            prop: 'newMerchantDuringPeriodCount'
          },
          {
            type: 'show',
            label: 'stats.cumMerchantCount',
            prop: 'cumMerchantCount'
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
        newMerchantDuringPeriodCount: true,
        cumMerchantCount: true
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
