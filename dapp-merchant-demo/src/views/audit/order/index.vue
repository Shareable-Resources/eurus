<template>
  <div id="orderAudit">
    <div class="filter-form">
      <filter-form
        :filter-form="filterFormVal"
      />
    </div>
    <basic-table :table-data="tableData" />
    <basic-dialog-form :dialog-data="dialogData" />
  </div>
</template>

<script>
import { TX_TYPE, MERCHANT_STATUS, USER_GROUP } from '@/utils/selection.const'
export default {
  name: 'OrderAudit',
  data() {
    return {
      TX_TYPE: TX_TYPE,
      MERCHANT_STATUS: MERCHANT_STATUS,
      USER_GROUP: USER_GROUP,
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
            model: 'txType',
            label: 'order.txType',
            options: TX_TYPE
          },
          {
            type: 'select',
            model: 'status',
            label: 'common.status',
            options: MERCHANT_STATUS
          },
          {
            type: 'select',
            model: 'userGroup',
            label: 'common.userGroup',
            options: USER_GROUP
          },
          {
            type: 'input',
            model: 'recevAddr',
            label: 'order.recevAddr'
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
        searchAction: this.searchAction
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
            label: 'order.txType',
            prop: 'txType',
            options: TX_TYPE
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
            type: 'show',
            label: 'order.publicChainHash',
            prop: 'publicChainHash'
          },
          {
            type: 'show',
            label: 'order.publicChainGas',
            prop: 'publicChainGas'
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
            prop: 'txType',
            label: 'order.txType',
            filterFn: (txType) => {
              if (!txType) return ''
              return this.$t(`options.${TX_TYPE[txType]}`)
            }
          },
          {
            type: 'filter',
            prop: 'userGroup',
            label: 'common.userGroup',
            filterFn: (userGroup) => {
              if (!userGroup) return ''
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
            prop: 'publicChainHash',
            label: 'order.publicChainHash'
          },
          {
            type: 'input',
            prop: 'publicChainGas',
            label: 'order.publicChainGas'
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
        approveAction: this.approveAction,
        rejectAction: this.rejectAction
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
    searchAction(data) {
      console.log('search', data)
    },
    handleClick(row) {
      this.$set(this.dialogData, 'source', row)
      this.$set(this.dialogData, 'visible', true)
    },
    handleClose(done) {
      this.$set(this.dialogData, 'source', {})
      return done()
    },
    approveAction() {
      const { source } = this.dialogData
      console.log('approve', source)
    },
    rejectAction() {
      const { source } = this.dialogData
      console.log('reject', source)
    }
  }
}
</script>
