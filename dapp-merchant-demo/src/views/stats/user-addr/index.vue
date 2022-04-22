<template>
  <div id="UserAddrStats">
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
import { USER_GROUP, STATS_USER_ADDR_TYPE } from '@/utils/selection.const'
export default {
  name: 'UserAddrStats',
  data() {
    return {
      USER_GROUP: USER_GROUP,
      STATS_USER_ADDR_TYPE: STATS_USER_ADDR_TYPE,
      filterFormVal: {
        columns: [
          {
            type: 'select',
            model: 'userGroup',
            label: 'common.userGroup',
            options: USER_GROUP
          },
          {
            type: 'select',
            model: 'assetType',
            label: 'stats.assetType',
            options: STATS_USER_ADDR_TYPE
          },
          {
            type: 'select',
            model: 'addrType',
            label: 'common.addrType',
            options: STATS_USER_ADDR_TYPE
          }
        ],
        default: {
          userGroup: '',
          assetType: '',
          addrType: ''
        },
        searchAction: this.searchAction,
        exportAction: this.exportAction
      },
      tableData: {
        source: [
          {
            userGroup: 'decentralized',
            assetType: 'addrType',
            addrType: 'mainChain',
            usedAmount: 6,
            availableAmount: 1,
            totalAmount: 7
          },
          {
            userGroup: 'centralized',
            assetType: 'addrType2',
            addrType: 'mainChain',
            usedAmount: 2,
            availableAmount: 1,
            totalAmount: 3
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
            type: 'options',
            label: 'common.addrType',
            prop: 'addrType',
            options: STATS_USER_ADDR_TYPE
          },
          {
            type: 'show',
            label: 'stats.assetType',
            prop: 'assetType'
          },
          {
            type: 'show',
            label: 'stats.usedAmount',
            prop: 'usedAmount'
          },
          {
            type: 'show',
            label: 'stats.availableAmount',
            prop: 'availableAmount'
          },
          {
            type: 'show',
            label: 'stats.totalAmount',
            prop: 'totalAmount'
          }
        ],
        currentPage: 1,
        total: 100,
        handleSizeChange: this.handleSizeChange,
        handleCurrentChange: this.handleCurrentChange
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
    }
  }
}
</script>
