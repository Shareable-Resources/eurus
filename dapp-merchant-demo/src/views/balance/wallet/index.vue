<template>
  <div id="wallet-balance">
    <div class="filter-form">
      <real-time
        :real-time-refresh-action="refreshRealTime"
      />
      <daily-report-export
        :daily-report-export-action="exportDailyReport"
      />
      <filter-form
        :filter-form="filterFormVal"
      />
    </div>

    <div class="table">
      <el-table
        border
        :data="tableData"
        :summary-method="getSummaries"
        show-summary
        style="width: 100%"
      >
        <el-table-column type="expand">
          <template slot-scope="props">
            <el-table
              :data="props.row.sub"
            >
              <el-table-column
                :label="$t(`walletBal.walletName`)"
                prop="walletName"
              />
              <el-table-column
                :label="$t(`common.addr`)"
                prop="addr"
              />
              <el-table-column
                :label="$t(`common.asset`)"
                prop="asset"
              />
              <el-table-column
                :label="$t(`common.total`)"
                prop="total"
              />
              <el-table-column
                :label="$t(`common.totalUSDT`)"
                prop="usdt_total"
              />
            </el-table>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t(`common.addrType`)"
          prop="addrType"
        >
          <template slot-scope="props">
            {{ $t(`options.${WALLET_BAL_ADDRTYPE[props.row.addrType]}`) }}
          </template>
        </el-table-column>
        <el-table-column
          width="200"
          :label="$t(`common.totalUSDT`)"
          prop="usdt_total"
        />
      </el-table>
    </div>
  </div>
</template>

<script>
import { WALLET_BAL_ADDRTYPE } from '@/utils/selection.const'
export default {
  name: 'WalletBalance',
  data() {
    return {
      WALLET_BAL_ADDRTYPE: WALLET_BAL_ADDRTYPE,
      filterFormVal: {
        columns: [
          {
            type: 'snapShot',
            model: 'snapShot',
            label: 'filter.snapShot'
          },
          {
            type: 'select',
            model: 'addrType',
            label: 'common.addrType',
            options: WALLET_BAL_ADDRTYPE
          },
          {
            type: 'input',
            model: 'designatedAddr',
            label: 'walletBal.designatedAddr',
            placeholder: ''
          },
          {
            type: 'select',
            model: 'asset',
            label: 'common.asset',
            options: WALLET_BAL_ADDRTYPE
          }
        ],
        default: {
          addrType: 'createWallet',
          designatedAddr: '12'
        },
        // addAction: (data)=>{console.log('add',data)},
        searchAction: this.searchAction,
        exportAction: this.exportAction
      },
      tableData: [
        {
          addrType: 'createWallet',
          usdt_total: '100.2',
          sub: [{
            walletName: 'name3',
            asset: 'ABC',
            addr: '213131231a2131312ddr',
            total: '1.3445',
            usdt_total: '16000'
          }, {
            walletName: 'name4',
            asset: 'CCC',
            addr: '213131231ad1313d221212r',
            total: '2.3125',
            usdt_total: '36001'
          }]
        },
        {
          addrType: 'groupWallet',
          usdt_total: '100.3',
          sub: [{
            walletName: 'name5',
            asset: 'CD',
            addr: '21313123112112a2131312ddr',
            total: '1.3445',
            usdt_total: '16000'
          }, {
            walletName: 'name6',
            asset: 'AC',
            addr: '213131231ad1313d221212r',
            total: '2.3125',
            usdt_total: '36001'
          }]
        },
        {
          addrType: 'feeWallet',
          usdt_total: '100.4',
          sub: [{
            walletName: 'name5',
            asset: 'CD',
            addr: '21313123112112a2131312ddr',
            total: '1.3445',
            usdt_total: '16000'
          }, {
            walletName: 'name6',
            asset: 'AC',
            addr: '213131231ad1313d221212r',
            total: '2.3125',
            usdt_total: '36001'
          }]
        },
        {
          addrType: 'withdrawWallet',
          usdt_total: '100.5',
          sub: [{
            walletName: 'name5',
            asset: 'CD',
            addr: '21313123112112a2131312ddr',
            total: '1.3445',
            usdt_total: '16000'
          }, {
            walletName: 'name6',
            asset: 'AC',
            addr: '213131231ad1313d221212r',
            total: '2.3125',
            usdt_total: '36001'
          }]
        }
      ]
    }
  },
  methods: {
    refreshRealTime(timestamp) {
      console.log(timestamp)
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
    getSummaries(param) {
      const { columns, data } = param
      const sums = []
      const COUNT_MAP = {
        usdt_total: true
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

<style lang="scss" scoped>

</style>
