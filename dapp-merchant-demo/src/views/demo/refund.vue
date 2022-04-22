<template>
  <div id="refund">
    <el-form ref="form" :model="form" status-icon>
      <el-form-item :label="$t(`common.asset`)" prop="assetName">
        <el-input v-model="form.assetName" autocomplete="off" :disabled="true"/>
      </el-form-item>
      <el-form-item :label="$t(`common.amount`)" prop="amount">
        <el-input v-model="form.amount" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Reason" prop="reason">
        <el-input v-model="form.reason" autocomplete="off" />
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
import { requestRefund } from '@/api'

export default {
  name: 'Refund',
  data() {
    return {
      form: {
        assetName: 'USDT',
        amount: 0,
        reason: ''
      },
      isLoading: false,
      resultText: ''
    }
  },
  methods: {
    async submitForm() {
      var amount = this.form.amount
      if (isNaN(amount)) {
        alert('Not a number!')
      } else if (amount <= 0) {
        alert('Less than 0!')
      } else {
        // alert(this.form.amount)
        this.isLoading = true
        try {
          let requestRefundResponse = await requestRefund(
            this.form.assetName, 
            this.form.amount, 
            this.form.reason
          );
          console.log('requestRefundResponse:', requestRefundResponse)
          if (
            requestRefundResponse &&
            requestRefundResponse.returnCode === 0
          ) {
            alert('Success!')
          } else {
            alert('Network error!')
          }
        } catch (error) {
          alert('Network error!')
          console.log('Network error:', error)
        }
        this.isLoading = false
      }
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    }
  }
}
</script>
