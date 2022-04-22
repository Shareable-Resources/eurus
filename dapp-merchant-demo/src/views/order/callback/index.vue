<template>
  <div id="callback">
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
import { USER_GROUP, MERCHANT_STATUS } from '@/utils/selection.const'
export default {
  name: 'Withdraw',
  data() {
    return {
      USER_GROUP: USER_GROUP,
      MERCHANT_STATUS: MERCHANT_STATUS,
      filterFormVal: {
        columns: [
          {
            type: 'input',
            model: 'userNo',
            label: 'common.userNo',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'orderID',
            label: 'common.orderID'
          },
          {
            type: 'input',
            model: 'originID',
            label: 'common.originID'
          },
          {
            type: 'input',
            model: 'merchantID',
            label: 'common.merchantID'
          },
          {
            type: 'input',
            model: 'privateChainHash',
            label: 'order.privateChainHash'
          },
          {
            type: 'input',
            model: 'privateChainFromAddr',
            label: 'order.privateChainFromAddr'
          },
          {
            type: 'input',
            model: 'callbackAddr',
            label: 'order.callbackAddr'
          },
          {
            type: 'select',
            model: 'type',
            label: 'common.type',
            options: MERCHANT_STATUS
          },
          {
            type: 'select',
            model: 'status',
            label: 'common.status',
            options: MERCHANT_STATUS
          },
          {
            type: 'datetimerange',
            model: 'createTime',
            label: 'common.createTime',
            format: 'yyyy/MM/dd HH:mm:ss'
          },
          {
            type: 'datetimerange',
            model: 'lastUpdateTime',
            label: 'order.lastUpdateTime',
            format: 'yyyy/MM/dd HH:mm:ss'
          },
          {
            type: 'input',
            model: 'privateChainToAddr',
            label: 'order.privateChainToAddr'
          }
        ],
        default: {
          userNo: 'userNo',
          userGroup: '',
          orderID: 'orderID',
          asset: 'asset',
          transactionHash: 'transactionHash',
          status: 'normal',
          recevAddr: 'recevAddr',
          createTime: new Date(),
          lastUpdateTime: new Date()
        },
        searchAction: this.searchAction,
        exportAction: this.exportAction
      },
      tableData: {
        source: [
          {
            userNo: 'userNo',
            userGroup: 'decentralized',
            orderID: 'orderID',
            amount: 'BTC',
            publicChainHash: '0x14b8ccee700ee7907bf8bf9ee08e64961dfe21ef7bf90971c382cd70be451ca3',
            asset: 'asset',
            transactionHash: 'transactionHash',
            status: 'normal',
            recevAddr: 'recevAddr',
            createTime: new Date(),
            lastUpdateTime: new Date()
          }
        ],
        columns: [
          {
            type: 'show',
            label: 'common.merchantID',
            prop: 'merchantID'
          },
          {
            type: 'show',
            label: 'common.userNo',
            prop: 'userNo'
          },
          {
            type: 'show',
            label: 'common.orderID',
            prop: 'orderID'
          },
          {
            type: 'show',
            label: 'common.originID',
            prop: 'originID'
          },
          {
            type: 'show',
            label: 'common.amount',
            prop: 'amount'
          },
          {
            type: 'show',
            label: 'order.otherSide',
            prop: 'otherSide'
          },
          {
            type: 'show',
            label: 'order.callbackAddr',
            prop: 'callbackAddr'
          },
          {
            type: 'options',
            label: 'common.type',
            prop: 'type',
            options: MERCHANT_STATUS
          },
          {
            type: 'time',
            label: 'common.createTime',
            prop: 'createTime'
          },
          {
            type: 'time',
            label: 'order.lastUpdateTime',
            prop: 'lastUpdateTime'
          },
          {
            type: 'options',
            label: 'common.status',
            prop: 'status',
            options: MERCHANT_STATUS
          },
          {
            type: 'show',
            label: 'common.msg',
            prop: 'msg'
          },
          {
            type: 'show',
            label: 'order.privateChainFromAddr',
            prop: 'privateChainFromAddr'
          },
          {
            type: 'show',
            label: 'order.privateChainToAddr',
            prop: 'privateChainToAddr'
          },
          {
            type: 'show',
            label: 'order.privateChainHash',
            prop: 'privateChainHash'
          },
          {
            type: 'btns',
            label: 'common.operation',
            btns: [
              {
                type: 'retry',
                fn: this.handleClick
              }
            ]
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
    handleSizeChange(val) {
      console.log(`每页 ${val} 条`)
    },
    handleCurrentChange(val) {
      console.log(`当前页: ${val}`)
    },
    exportDailyReport(date) {
      console.log('daily report', date)
    },
    searchAction(data) {
      console.log('search', data)
    },
    exportAction(data) {
      console.log('export', data)
    },
    handleClick(row) {
      console.log('retry', row)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
