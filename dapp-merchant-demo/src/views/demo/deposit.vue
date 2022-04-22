<template>
  <div id="deposit">
    <el-form ref="form" :model="form" status-icon>
      <el-form-item :label="$t(`common.amount`)" prop="amount">
        <el-input v-model="form.amount" autocomplete="off" />
      </el-form-item>

      <el-form-item>
        <el-button @click="resetForm('form')">{{
          $t(`common.reset`)
        }}</el-button>
        <el-button :loading="isLoading" type="primary" @click="submitForm()">{{
          $t(`common.confirm`)
        }}</el-button>
      </el-form-item>
    </el-form>

    <div>
      {{ this.resultText }}
    </div>
  </div>
</template>

<script>
import { depositToDapp } from '@/web3'
import { getUsdtAddress, getDappAddress } from '@/utils/auth'

export default {
  name: 'Deposit',
  data() {
    return {
      form: {
        amount: 0
      },
      isLoading: false,
      resultText: ''
    }
  },
  methods: {
    submitForm() {
      var amount = this.form.amount
      if (isNaN(amount)) {
        alert('Not a number!')
      } else if (amount <= 0) {
        alert('Less than 0!')
      } else {
        // alert(this.form.amount)
        this.isLoading = true
        depositToDapp(getUsdtAddress(), getDappAddress(), amount)
          .then((txnHash) => {
            this.resultText = 'Txn hash: ' + txnHash
            console.log('###### txnHash: ', txnHash)
            this.isLoading = false
            alert('Success!')
          })
          .catch((error) => {
            this.isLoading = false
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
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    }
  }
}
</script>
