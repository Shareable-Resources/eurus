<template>
  <div id="configFee">
    <filter-form
      :filter-form="filterFormVal"
    />

    <basic-table :table-data="tableData" />

    <el-dialog
      width="80%"
      :close-on-click-modal="false"
      :title="dialogForm.id ? $t(`common.details`):$t(`common.create`)"
      :visible.sync="dialogFormVisible"
      :before-close="handleClose"
    >
      <el-form
        label-position="left"
        :model="dialogForm"
      >
        <el-form-item
          :label="$t(`config.appidOrMerchant`)"
        >
          <el-input
            v-model="dialogForm.appidOrMerchant"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.asset`)"
        >
          <el-input
            v-model="dialogForm.asset"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`order.txType`)"
        >
          <el-select
            v-model="dialogForm.txType"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_FEE_TRANS_TYPE"
              :key="label"
              :label="$t(`options.${CONFIG_FEE_TRANS_TYPE[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.feeRate`)"
        >
          <el-input
            v-model="dialogForm.feeRate"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.fixFee`)"
        >
          <el-input
            v-model="dialogForm.fixFee"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.feePaidBy`)"
        >
          <el-select
            v-model="dialogForm.feePaidBy"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_FEE_PAID_BY"
              :key="label"
              :label="$t(`options.${CONFIG_FEE_PAID_BY[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.maxAutoApprovalAmount`)"
        >
          <el-input
            v-model="dialogForm.maxAutoApprovalAmount"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.status`)"
        >
          <el-select
            v-model="dialogForm.status"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_FEE_STATUS"
              :key="label"
              :label="$t(`options.${CONFIG_FEE_STATUS[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <div v-if="dialogForm.id && !isEdit">
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
            @click="dialogForm.id? isEdit = false : dialogFormVisible = false"
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
import { CONFIG_FEE_STATUS, CONFIG_FEE_PAID_BY, CONFIG_FEE_TRANS_TYPE } from '@/utils/selection.const'
export default {
  name: 'ConfigAsset',
  data() {
    return {
      CONFIG_FEE_STATUS: CONFIG_FEE_STATUS,
      CONFIG_FEE_PAID_BY: CONFIG_FEE_PAID_BY,
      CONFIG_FEE_TRANS_TYPE: CONFIG_FEE_TRANS_TYPE,
      filterFormVal: {
        columns: [
          {
            type: 'input',
            model: 'appidOrMerchant',
            label: 'config.appidOrMerchant',
            placeholder: ''
          },
          {
            type: 'select',
            model: 'asset',
            label: 'common.asset',
            options: CONFIG_FEE_STATUS
          },
          {
            type: 'select',
            model: 'txType',
            label: 'order.txType',
            options: CONFIG_FEE_TRANS_TYPE
          },
          {
            type: 'select',
            model: 'status',
            label: 'common.status',
            options: CONFIG_FEE_STATUS
          }
        ],
        default: {
          roleName: '',
          status: ''
        },
        addAction: this.addAction,
        searchAction: this.searchAction
      },
      tableData: {
        source: [
          {
            appidOrMerchant: 'appid',
            asset: 'asset',
            createTime: new Date(),
            creator: 'creator',
            updateTime: new Date(),
            updator: 'updator',
            id: '12311',
            status: 'normal',
            txType: 'payment'
          }
        ],
        columns: [
          {
            type: 'show',
            label: 'config.appidOrMerchant',
            prop: 'appidOrMerchant'
          },
          {
            type: 'show',
            label: 'common.asset',
            prop: 'asset'
          },
          {
            type: 'options',
            label: 'order.txType',
            prop: 'txType',
            options: CONFIG_FEE_TRANS_TYPE
          },
          {
            type: 'show',
            label: 'config.feeRate',
            prop: 'feeRate'
          },
          {
            type: 'show',
            label: 'config.fixFee',
            prop: 'fixFee'
          },
          {
            type: 'options',
            label: 'order.feePaidBy',
            prop: 'feePaidBy',
            options: CONFIG_FEE_PAID_BY
          },
          {
            type: 'show',
            label: 'config.maxAutoApprovalAmount',
            prop: 'maxAutoApprovalAmount'
          },
          {
            type: 'options',
            label: 'common.status',
            prop: 'status',
            options: CONFIG_FEE_STATUS
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
      isEdit: false
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
      console.log(data, 'search')
    },
    addAction() {
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
    }
  }
}
</script>
