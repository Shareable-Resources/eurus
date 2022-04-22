<template>
  <div id="changePW">
    <el-form
      ref="form"
      :model="form"
      :rules="rules"
      status-icon
    >
      <el-form-item
        :label="$t(`admin.oldPW`)"
        prop="oldPW"
      >
        <el-input
          v-model="form.oldPW"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        :label="$t(`admin.newPW`)"
        prop="newPW"
      >
        <el-input
          v-model="form.newPW"
          show-password
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        :label="$t(`admin.cfNewPW`)"
        prop="cfNewPW"
      >
        <el-input
          v-model="form.cfNewPW"
          show-password
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item>
        <el-button
          @click="resetForm('form')"
        >{{ $t(`common.reset`) }}</el-button>
        <el-button
          type="primary"
          @click="submitForm('form')"
        >{{ $t(`common.confirm`) }}</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: 'ChangePW',
  data() {
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(this.$t(`admin.inpPwTip`)))
      } else {
        if (this.form.cfNewPW !== '') {
          this.$refs.form.validateField('cfNewPW')
        }
        callback()
      }
    }
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(this.$t(`admin.inpPwTip2`)))
      } else if (value !== this.form.newPW) {
        callback(new Error(this.$t(`admin.pwErrTip`)))
      } else {
        callback()
      }
    }
    return {
      form: {
        oldPW: '',
        newPW: '',
        cfNewPW: ''
      },
      rules: {
        newPW: [
          { validator: validatePass, trigger: 'blur' }
        ],
        cfNewPW: [
          { validator: validatePass2, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          alert('submit!')
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    }
  }
}
</script>
