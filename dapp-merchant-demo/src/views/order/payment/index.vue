<template>
  <div id="payment">
    <div class="filter-form">
      <daily-report-export
        :daily-report-export-action="exportDailyReport"
      />
      <filter-form
        :filter-form="filterFormVal"
      />
    </div>

    <basic-table :table-data="tableData" />

    <basic-dialog-form :dialog-data="dialogData" />
  </div>
</template>

<script>
import { USER_GROUP, MERCHANT_STATUS } from '@/utils/selection.const'
export default {
  name: 'Payment',
  data() {
    return {
      USER_GROUP: USER_GROUP,
      MERCHANT_STATUS: MERCHANT_STATUS,
      filterFormVal: {
        columns: [
          {
            type: 'input',
            model: 'merchantID',
            label: 'common.merchantID',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'fromUser',
            label: 'order.fromUser',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'toUser',
            label: 'order.toUser',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'orderID',
            label: 'common.orderID',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'originID',
            label: 'common.originID',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'privateChainHash',
            label: 'order.privateChainHash',
            placeholder: ''
          },
          {
            type: 'select',
            model: 'asset',
            label: 'common.asset',
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
            model: 'privateChainFromAddr',
            label: 'order.privateChainFromAddr',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'privateChainToAddr',
            label: 'order.privateChainToAddr',
            placeholder: ''
          }
        ],
        default: {
          merchantID: 'merchantID',
          fromUser: 'fromUser',
          orderID: 'orderID',
          privateChainHash: 'privateChainHash',
          asset: 'asset',
          status: 'normal',
          createTime: new Date(),
          lastUpdateTime: new Date(),
          privateChainFromAddr: 'privateChainFromAddr',
          privateChainToAddr: 'privateChainToAddr'
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
            label: 'common.merchantName',
            prop: 'merchantName'
          },
          {
            type: 'show',
            label: 'order.fromUser',
            prop: 'fromUser'
          },
          {
            type: 'show',
            label: 'order.toUser',
            prop: 'toUser'
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
            label: 'common.asset',
            prop: 'asset'
          },
          {
            type: 'show',
            label: 'order.fee',
            prop: 'fee'
          },
          {
            type: 'show',
            label: 'order.paidTotal',
            prop: 'paidTotal'
          },
          {
            type: 'options',
            label: 'common.status',
            prop: 'status',
            options: MERCHANT_STATUS
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
            type: 'show',
            label: 'order.privateChainGas',
            prop: 'privateChainGas'
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
            type: 'btns',
            label: 'common.operation',
            btns: [
              {
                type: 'details',
                fn: this.handleClick
              }
            ]
          }
        ],
        currentPage: 1,
        total: 100,
        handleSizeChange: this.handleSizeChange,
        handleCurrentChange: this.handleCurrentChange
      },
      dialogData: {
        visible: false,
        labelWidth: '150px',
        beforeClose: this.handleClose,
        source: {},
        form: [
          {
            type: 'input',
            prop: 'merchantID',
            label: 'common.merchantID'
          },
          {
            type: 'input',
            prop: 'merchantName',
            label: 'common.merchantName'
          },
          {
            type: 'input',
            prop: 'orderID',
            label: 'common.orderID'
          },
          {
            type: 'input',
            prop: 'originID',
            label: 'common.originID'
          },
          {
            type: 'input',
            prop: 'fromUser',
            label: 'order.fromUser'
          },
          {
            type: 'input',
            prop: 'toUser',
            label: 'order.toUser'
          },
          {
            type: 'input',
            prop: 'asset',
            label: 'common.asset'
          },
          {
            type: 'input',
            prop: 'amount',
            label: 'common.amount'
          },
          {
            type: 'input',
            prop: 'fee',
            label: 'order.fee'
          },
          {
            type: 'input',
            prop: 'paidTotal',
            label: 'order.paidTotal'
          },
          {
            type: 'input',
            prop: 'feePaidBy',
            label: 'order.feePaidBy'
          },
          {
            type: 'select',
            prop: 'status',
            label: 'common.status',
            options: MERCHANT_STATUS
          },
          {
            type: 'input',
            prop: 'fromFromBal',
            label: 'order.fromFromBal'
          },
          {
            type: 'input',
            prop: 'fromFromFrozen',
            label: 'order.fromFromFrozen'
          },
          {
            type: 'input',
            prop: 'fromToBal',
            label: 'order.fromToBal'
          },
          {
            type: 'input',
            prop: 'fromToFrozen',
            label: 'order.fromToFrozen'
          },
          {
            type: 'input',
            prop: 'toFromBal',
            label: 'order.toFromBal'
          },
          {
            type: 'input',
            prop: 'toFromFrozen',
            label: 'order.toFromFrozen'
          },
          {
            type: 'input',
            prop: 'toToBal',
            label: 'order.toToBal'
          },
          {
            type: 'input',
            prop: 'toToFrozen',
            label: 'order.toToFrozen'
          },
          {
            type: 'input',
            prop: 'privateChainFromAddr',
            label: 'order.privateChainFromAddr'
          },
          {
            type: 'input',
            prop: 'privateChainToAddr',
            label: 'order.privateChainToAddr'
          },
          {
            type: 'input',
            prop: 'privateChainHash',
            label: 'order.privateChainHash'
          },
          {
            type: 'input',
            prop: 'privateChainGas',
            label: 'order.privateChainGas'
          },
          {
            type: 'time',
            prop: 'createTime',
            label: 'common.createTime'
          },
          {
            type: 'time',
            prop: 'lastUpdateTime',
            label: 'order.lastUpdateTime'
          }
        ],
        cancelAction: this.handleCancel,
        refundAction: this.handleRefund,
        auditAction: this.handleAudit
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
      this.$set(this.dialogData, 'source', row)
      this.$set(this.dialogData, 'visible', true)
    },
    handleClose(done) {
      this.$set(this.dialogData, 'source', {})
      return done()
    },
    handleCancel() {
      this.$set(this.dialogData, 'source', {})
      this.$set(this.dialogData, 'visible', false)
    },
    handleRefund(pw) {
      console.log(pw, 'pw refund')
    },
    handleAudit(pw) {
      console.log(pw, 'pw audit')
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
