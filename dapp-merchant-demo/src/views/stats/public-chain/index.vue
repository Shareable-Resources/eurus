<template>
  <div id="PublicChainStats">
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
import { STATS_PC_WALLET_GROUP } from '@/utils/selection.const'
export default {
  name: 'PublicChainStats',
  data() {
    return {
      STATS_PC_WALLET_GROUP: STATS_PC_WALLET_GROUP,
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
            options: STATS_PC_WALLET_GROUP
          },
          {
            type: 'input',
            model: 'walletAddr',
            label: 'stats.walletAddr',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'asset',
            label: 'common.asset',
            placeholder: ''
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
            walletType: 'genAddrWallet',
            walletName: 'walletName',
            walletAddr: 'walletAddr',
            mainChainGas_ETH: 6,
            mainChainGas_BTC: 1,
            mainChainGas_TRX: 2
          },
          {
            walletType: 'sweepWallet',
            walletName: 'walletName2',
            walletAddr: 'walletAddr2',
            mainChainGas_ETH: 3,
            mainChainGas_BTC: 5,
            mainChainGas_TRX: 1
          }
        ],
        columns: [
          {
            type: 'options',
            label: 'stats.walletType',
            prop: 'walletType',
            options: STATS_PC_WALLET_GROUP
          },
          {
            type: 'show',
            label: 'walletBal.walletName',
            prop: 'walletName'
          },
          {
            type: 'show',
            label: 'stats.walletAddr',
            prop: 'walletAddr'
          },
          {
            type: 'show',
            label: 'stats.mainChainGas_ETH',
            prop: 'mainChainGas_ETH'
          },
          {
            type: 'show',
            label: 'stats.mainChainGas_BTC',
            prop: 'mainChainGas_BTC'
          },
          {
            type: 'show',
            label: 'stats.mainChainGas_TRX',
            prop: 'mainChainGas_TRX'
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
        mainChainGas_ETH: true,
        mainChainGas_BTC: true,
        mainChainGas_TRX: true
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
