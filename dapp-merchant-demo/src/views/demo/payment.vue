<template>
  <div id="payment">
    <basic-table :table-data="tableData" />
    <div>
      {{ this.resultText }}
    </div>
  </div>
</template>

<script>
import { getDAppProductList, purchaseToDapp } from '@/web3'
import { getUsdtAddress, getDappStockAddress } from '@/utils/auth'
import { Loading } from 'element-ui'

export default {
  name: 'Payment',
  data() {
    return {
      isLoading: false,
      resultText: '',
      tableData: {
        source: [],
        columns: [
          {
            type: 'show',
            label: 'Product Id',
            prop: 'productId'
          },
          {
            type: 'show',
            label: 'Name',
            prop: 'name'
          },
          {
            type: 'show',
            label: 'Price(USDT)',
            prop: 'price'
          },
          {
            type: 'show',
            label: 'Stock',
            prop: 'stock'
          },
          {
            type: 'show',
            label: 'OnShelf',
            prop: 'onShelf'
          },
          {
            type: 'btns',
            label: 'common.operation',
            btns: [
              {
                type: 'btn',
                label: 'Purchase',
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
  created() {},
  async mounted() {
    let productArr = await getDAppProductList(getDappStockAddress())
    this.tableData.source = productArr
    this.tableData.total = productArr.length
  },
  methods: {
    handleSizeChange(val) {
      console.log(`每页 ${val} 条`)
    },
    handleCurrentChange(val) {
      console.log(`当前页: ${val}`)
    },
    refreshRealTime(timestamp) {
      console.log(timestamp)
    },
    async handleClick(row) {
      console.log('handleClick start')
      console.log(row)
      console.log(row.productId)
      this.isLoading = true
      const loadingInstance = Loading.service({ fullscreen: true })
      purchaseToDapp(
        getUsdtAddress(),
        getDappStockAddress(),
        row.productId,
        1,
        row.price
      )
        .then((txnHash) => {
          this.resultText = 'Txn hash: ' + txnHash
          console.log('###### txnHash: ', txnHash)
          this.isLoading = false
          loadingInstance.close()
          alert('Success!')
        })
        .catch((error) => {
          this.isLoading = false
          loadingInstance.close()
          if (error.code === 4001) {
            // EIP-1193 userRejectedRequest error
            console.log('Please connect to MetaMask.')
            alert('Please connect to MetaMask.')
          } else if (error === 9991) {
            alert('Metamask wrong network!')
            console.log('Metamask wrong network!')
          } else {
            console.error(error)
          }
        })
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
