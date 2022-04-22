<template>
  <div id="deposit">
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
  name: 'Deposit',
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
            type: 'select',
            model: 'userGroup',
            label: 'common.userGroup',
            options: USER_GROUP
          },
          {
            type: 'input',
            model: 'orderID',
            label: 'common.orderID'
          },
          {
            type: 'select',
            model: 'asset',
            label: 'common.asset',
            options: MERCHANT_STATUS
          },
          {
            type: 'input',
            model: 'transactionHash',
            label: 'common.transactionHash'
          },
          {
            type: 'select',
            model: 'status',
            label: 'common.status',
            options: MERCHANT_STATUS
          },
          {
            type: 'input',
            model: 'fromAddr',
            label: 'order.fromAddr'
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
          }
        ],
        default: {
          userNo: 'userNo',
          userGroup: '',
          orderID: 'orderID',
          asset: 'asset',
          transactionHash: 'transactionHash',
          status: 'normal',
          fromAddr: 'fromAddr',
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
            fromAddr: 'fromAddr',
            createTime: new Date(),
            lastUpdateTime: new Date()
          }
        ],
        columns: [
          {
            type: 'show',
            label: 'common.userNo',
            prop: 'userNo'
          },
          {
            type: 'options',
            label: 'common.userGroup',
            prop: 'userGroup',
            options: USER_GROUP
          },
          {
            type: 'show',
            label: 'common.orderID',
            prop: 'orderID'
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
            label: 'order.publicChainHash',
            prop: 'publicChainHash'
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
        beforeClose: this.handleClose,
        source: {},
        form: [
          {
            type: 'input',
            prop: 'userNo',
            label: 'common.userNo'
          },
          {
            type: 'filter',
            prop: 'userGroup',
            label: 'common.userGroup',
            filterFn: (userGroup) => {
              return this.$t(`options.${USER_GROUP[userGroup]}`)
            }
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
            prop: 'amount',
            label: 'common.amount'
          },
          {
            type: 'input',
            prop: 'asset',
            label: 'common.asset'
          },
          {
            type: 'input',
            prop: 'publicChainFromAddr',
            label: 'order.publicChainFromAddr'
          },
          {
            type: 'input',
            prop: 'publicChainToAddr',
            label: 'order.publicChainToAddr'
          },
          {
            type: 'input',
            prop: 'fromBal',
            label: 'order.fromBal'
          },
          {
            type: 'input',
            prop: 'fromFrozen',
            label: 'order.fromFrozen'
          },
          {
            type: 'input',
            prop: 'toBal',
            label: 'order.toBal'
          },
          {
            type: 'input',
            prop: 'toFrozen',
            label: 'order.toFrozen'
          },
          {
            type: 'input',
            prop: 'publicChainHash',
            label: 'order.publicChainHash'
          },
          {
            type: 'select',
            prop: 'status',
            label: 'common.status',
            options: MERCHANT_STATUS
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
        cancelAction: this.handleCancel
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
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
