<template>
  <div id="merchant">
    <div class="filter-form">
      <daily-report-export
        :daily-report-export-action="exportDailyReport"
      />
      <filter-form
        :filter-form="filterFormVal"
      />
    </div>

    <basic-table :table-data="tableData" />

    <el-dialog
      width="80%"
      :close-on-click-modal="false"
      :title="dialogForm.merchantID ? $t(`common.details`):$t(`common.create`)"
      :visible.sync="dialogFormVisible"
      :before-close="handleClose"
    >
      <el-form
        label-position="left"
        :model="dialogForm"
      >
        <el-form-item
          :label="$t(`common.merchantID`)"
        >
          <el-input
            v-model="dialogForm.merchantID"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.userNo`)"
        >
          <el-input
            v-model="dialogForm.userNo"
            :readonly="dialogForm.merchantID && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.merchantName`)"
        >
          <el-input
            v-model="dialogForm.merchantName"
            :readonly="dialogForm.merchantID && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`merchant.businessType`)"
        >
          <el-select
            v-model="dialogForm.businessType"
            :disabled="dialogForm.merchantID && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in MERCHANT_BUSINESSTYPE"
              :key="label"
              :label="$t(`options.${MERCHANT_BUSINESSTYPE[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`common.status`)"
        >
          <el-select
            v-model="dialogForm.status"
            :disabled="dialogForm.merchantID && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in MERCHANT_STATUS"
              :key="label"
              :label="$t(`options.${MERCHANT_STATUS[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`merchant.pubkey`)"
        >
          <el-input
            v-model="dialogForm.pubkey"
            readonly
            type="textarea"
          />
        </el-form-item>
        <el-form-item>
          <el-upload
            v-show="!dialogForm.merchantID || isEdit"
            action=""
            :data="dataObj"
            :on-change="handleUploadChange"
            :on-exceed="handleUploadExceed"
            :on-remove="handleUploadRemove"
            :auto-upload="false"
            :multiple="false"
            :limit="1"
          >
            <el-button
              size="small"
              type="primary"
            >{{ $t(`common.upload`) }}
            </el-button>
            <div slot="tip" class="el-upload__tip">{{ $t(`merchant.uploadTip`) }}</div>
          </el-upload>
        </el-form-item>
        <el-form-item
          :label="$t(`merchant.assetCBURL`)"
        >
          <el-input
            v-model="dialogForm.assetCBURL"
            :readonly="dialogForm.merchantID && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`merchant.queryCBURL`)"
        >
          <el-input
            v-model="dialogForm.queryCBURL"
            :readonly="dialogForm.merchantID && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`merchant.MD5Key`)"
        >
          <el-input
            v-model="dialogForm.MD5Key"
            :readonly="dialogForm.merchantID && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`merchant.reportPW`)"
        >
          <el-input
            v-model="dialogForm.reportPW"
            :readonly="dialogForm.merchantID && !isEdit"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <div v-if="dialogForm.userNo && !isEdit">
          <el-button
            @click="dialogFormVisible = false"
          >{{ $t(`common.close`) }}
          </el-button>
          <el-button
            type="danger"
            @click="isEdit = true"
          >{{ $t(`common.edit`) }}
          </el-button>
        </div>
        <div v-else>
          <el-button
            @click="dialogForm.userNo? isEdit = false : dialogFormVisible = false"
          >{{ $t(`common.cancel`) }}
          </el-button>
          <el-button
            type="primary"
            @click="handleSubmit"
          >{{ $t(`common.confirm`) }}
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { MERCHANT_BUSINESSTYPE, MERCHANT_STATUS } from '@/utils/selection.const'
import { readTxtFile } from '@/utils'
export default {
  name: 'Merchant',
  data() {
    return {
      MERCHANT_BUSINESSTYPE: MERCHANT_BUSINESSTYPE,
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
            model: 'userNo',
            label: 'common.userNo'
          },
          {
            type: 'input',
            model: 'merchantName',
            label: 'common.merchantName'
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
            type: 'select',
            model: 'businessType',
            label: 'merchant.businessType',
            options: MERCHANT_BUSINESSTYPE
          }
        ],
        default: {
          merchantID: 'merchantID',
          userNo: 'userNo',
          merchantName: 'merchantName',
          status: 'normal',
          createTime: '',
          businessType: ''
        },
        addAction: this.addAction,
        searchAction: this.searchAction,
        exportAction: this.exportAction
      },
      tableData: {
        source: [
          {
            merchantID: 'merchantID',
            userNo: 'userNo',
            merchantName: 'merchantName',
            businessType: 'others',
            createTime: new Date(),
            creator: 'creator',
            updateTime: new Date(),
            updator: 'updator',
            status: 'normal'
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
            label: 'common.userNo',
            prop: 'userNo'
          },
          {
            type: 'show',
            label: 'common.merchantName',
            prop: 'merchantName'
          },
          {
            type: 'options',
            label: 'merchant.businessType',
            prop: 'businessType',
            options: MERCHANT_BUSINESSTYPE
          },
          {
            type: 'time',
            label: 'common.createTime',
            prop: 'createTime'
          },
          {
            type: 'show',
            label: 'common.creator',
            prop: 'creator'
          },
          {
            type: 'time',
            label: 'common.updateTime',
            prop: 'updateTime'
          },
          {
            type: 'show',
            label: 'common.updator',
            prop: 'updator'
          },
          {
            type: 'options',
            label: 'common.status',
            prop: 'status',
            options: MERCHANT_STATUS
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
      dialogFormVisible: false,
      dialogForm: {},
      isEdit: false,
      dataObj: {}
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
    addAction(data) {
      console.log('add', data)
      this.dialogForm = {}
      this.dialogFormVisible = true
    },
    handleClick(row) {
      this.isEdit = false
      this.dialogForm = Object.assign({}, row)
      this.dialogFormVisible = true
    },
    handleClose(done) {
      this.dialogForm = {}
      return done()
    },
    handleSubmit() {
      this.dialogFormVisible = false
    },
    handleUploadChange(file, fileList) {
      const that = this
      readTxtFile(file.raw)
        .then(data => {
          // console.log(data,'merchant')
          that.$set(that.dialogForm, 'pubkey', data)
        })
    },
    handleUploadRemove() {
      this.$set(this.dialogForm, 'pubkey', '')
    },
    handleUploadExceed(file, fileList) {
      this.$message.error(`Please remove the exist one before upload`)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
