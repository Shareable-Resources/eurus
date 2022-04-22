<template>
  <div id="UserStats">
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
  name: 'UserStats',
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
            model: 'userGroup',
            label: 'common.userGroup',
            options: USER_GROUP
          },
          {
            type: 'input',
            model: 'registerFrom',
            label: 'common.registerFrom',
            placeholder: ''
          }
        ],
        default: {
          userNo: '',
          appidOrMerchant: '',
          startDate: '',
          endDate: ''
        },
        searchAction: this.searchAction,
        exportAction: this.exportAction
      },
      tableData: {
        source: [
          {
            userGroup: 'decentralized',
            registerFrom: 'registerFrom',
            newUserDuringPeriodCount: 6,
            cumUserCount: 1,
            kycUserDuringPeriodCount: 3,
            cumKycUserDuringPeriodCount: 2,
            activeUserDuringPeriodCount: 1,
            paidUserDuringPeriodCount: 5,
            transferOutUserDuringPeriodCount: 2,
            transferInUserDuringPeriodCount: 8
          },
          {
            userGroup: 'centralized',
            registerFrom: '1',
            newUserDuringPeriodCount: 2,
            cumUserCount: 1,
            kycUserDuringPeriodCount: 3,
            cumKycUserDuringPeriodCount: 2,
            activeUserDuringPeriodCount: 1,
            paidUserDuringPeriodCount: 4,
            transferOutUserDuringPeriodCount: 2,
            transferInUserDuringPeriodCount: 8
          }
        ],
        columns: [
          {
            type: 'options',
            label: 'common.userGroup',
            prop: 'userGroup',
            options: USER_GROUP
          },
          {
            type: 'show',
            label: 'common.registerFrom',
            prop: 'registerFrom'
          },
          {
            type: 'show',
            label: 'stats.newUserDuringPeriodCount',
            prop: 'newUserDuringPeriodCount'
          },
          {
            type: 'show',
            label: 'stats.cumUserCount',
            prop: 'cumUserCount'
          },
          {
            type: 'show',
            label: 'stats.kycUserDuringPeriodCount',
            prop: 'kycUserDuringPeriodCount'
          },
          {
            type: 'show',
            label: 'stats.cumKycUserDuringPeriodCount',
            prop: 'cumKycUserDuringPeriodCount'
          },
          {
            type: 'show',
            label: 'stats.activeUserDuringPeriodCount',
            prop: 'activeUserDuringPeriodCount'
          },
          {
            type: 'show',
            label: 'stats.paidUserDuringPeriodCount',
            prop: 'paidUserDuringPeriodCount'
          },
          {
            type: 'show',
            label: 'stats.transferOutUserDuringPeriodCount',
            prop: 'transferOutUserDuringPeriodCount'
          },
          {
            type: 'show',
            label: 'stats.transferInUserDuringPeriodCount',
            prop: 'transferInUserDuringPeriodCount'
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
        newUserDuringPeriodCount: true,
        cumUserCount: true,
        kycUserDuringPeriodCount: true,
        cumKycUserDuringPeriodCount: true,
        activeUserDuringPeriodCount: true,
        paidUserDuringPeriodCount: true,
        transferOutUserDuringPeriodCount: true,
        transferInUserDuringPeriodCount: true
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
