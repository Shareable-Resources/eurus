<template>
  <div id="wallet-balance">
    <div class="table">
      <el-table border :data="tableData">
        <el-table-column :label="$t(`common.addrType`)" prop="addrType" />
        <el-table-column :label="$t(`common.total`)" prop="usdt_total" />
      </el-table>
    </div>
  </div>
</template>

<script>
import { getErc20Balance } from '@/web3'
import { getAddress, getUsdtAddress, getDappAddress } from '@/utils/auth'

export default {
  name: 'WalletBalance',
  data() {
    return {
      tableData: [
        // {
        //   addrType: 'USDT',
        //   usdt_total: '100.2'
        // }
      ]
    }
  },
  created() {},
  async mounted() {
    // USDT
    const usdt = await getErc20Balance(
      getAddress(),
      getUsdtAddress()
    )
    this.tableData.push({
      addrType: 'USDT',
      usdt_total: usdt
    })

    const dapp = await getErc20Balance(
      getAddress(),
      getDappAddress()
    )
    this.tableData.push({
      addrType: 'DAPP token',
      usdt_total: dapp
    })
  },
  methods: {
    refreshRealTime(timestamp) {
      console.log(timestamp)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
