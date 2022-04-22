<template>
  <div id="configAsset">
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
          :label="$t(`config.unit`)"
        >
          <el-input
            v-model="dialogForm.unit"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.addrType`)"
        >
          <el-select
            v-model="dialogForm.addrType"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_ASSET_ADDR_TYPE"
              :key="label"
              :label="val"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.withdrawFee`)"
        >
          <el-select
            v-model="dialogForm.withdrawFee"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_ASSET_WITHDRAW_FEE"
              :key="label"
              :label="$t(`options.${CONFIG_ASSET_WITHDRAW_FEE[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.withdrawFeeAmount`)"
        >
          <el-input
            v-model="dialogForm.withdrawFeeAmount"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.minWithdrawAmount`)"
        >
          <el-input
            v-model="dialogForm.minWithdrawAmount"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.maxWithdrawAmount`)"
        >
          <el-input
            v-model="dialogForm.maxWithdrawAmount"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.maxWithdrawAmountDaily0`)"
        >
          <el-input
            v-model="dialogForm.maxWithdrawAmountDaily0"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.maxWithdrawAmountDaily1`)"
        >
          <el-input
            v-model="dialogForm.maxWithdrawAmountDaily1"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.maxWithdrawAmountDaily2`)"
        >
          <el-input
            v-model="dialogForm.maxWithdrawAmountDaily2"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.maxWithdrawAmountDaily3`)"
        >
          <el-input
            v-model="dialogForm.maxWithdrawAmountDaily3"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.minDepositAmount`)"
        >
          <el-input
            v-model="dialogForm.minDepositAmount"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.maxWithdrawAmountApproval`)"
        >
          <el-input
            v-model="dialogForm.maxWithdrawAmountApproval"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.maxDepositAmountApproval`)"
        >
          <el-input
            v-model="dialogForm.maxDepositAmountApproval"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.displaySort`)"
        >
          <el-input
            v-model="dialogForm.displaySort"
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
              v-for="(label, val) in CONFIG_ASSET_STATUS"
              :key="label"
              :label="$t(`options.${CONFIG_ASSET_STATUS[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.appVisible`)"
        >
          <el-select
            v-model="dialogForm.appVisible"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_ASSET_APP_VISIBLE"
              :key="label"
              :label="$t(`options.${CONFIG_ASSET_APP_VISIBLE[label]}`)"
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
import { CONFIG_ASSET_STATUS, CONFIG_ASSET_ADDR_TYPE, CONFIG_ASSET_APP_VISIBLE, CONFIG_ASSET_WITHDRAW_FEE } from '@/utils/selection.const'
export default {
  name: 'ConfigAsset',
  data() {
    return {
      CONFIG_ASSET_STATUS: CONFIG_ASSET_STATUS,
      CONFIG_ASSET_ADDR_TYPE: CONFIG_ASSET_ADDR_TYPE,
      CONFIG_ASSET_APP_VISIBLE: CONFIG_ASSET_APP_VISIBLE,
      CONFIG_ASSET_WITHDRAW_FEE: CONFIG_ASSET_WITHDRAW_FEE,
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
            options: CONFIG_ASSET_STATUS
          },
          {
            type: 'select',
            model: 'addrType',
            label: 'common.addrType',
            options: CONFIG_ASSET_ADDR_TYPE
          },
          {
            type: 'select',
            model: 'status',
            label: 'common.status',
            options: CONFIG_ASSET_STATUS
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
            unit: '8',
            addrType: 'TYPE_ETH',
            createTime: new Date(),
            creator: 'creator',
            lastUpdateTime: new Date(),
            updator: 'updator',
            id: '12311',
            status: 'listed',
            appVisible: 'visible',
            withdrawFee: 'fixAmount'
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
            type: 'show',
            label: 'config.unit',
            prop: 'unit'
          },
          {
            type: 'show',
            label: 'common.addrType',
            prop: 'addrType'
          },
          {
            type: 'options',
            label: 'config.withdrawFee',
            prop: 'withdrawFee',
            options: CONFIG_ASSET_WITHDRAW_FEE
          },
          {
            type: 'show',
            label: 'config.minWithdrawAmount',
            prop: 'minWithdrawAmount'
          },
          {
            type: 'show',
            label: 'config.maxWithdrawAmount',
            prop: 'maxWithdrawAmount'
          },
          {
            type: 'show',
            label: 'config.maxWithdrawAmountDaily0',
            prop: 'maxWithdrawAmountDaily0'
          },
          {
            type: 'show',
            label: 'config.maxWithdrawAmountDaily1',
            prop: 'maxWithdrawAmountDaily1'
          },
          {
            type: 'show',
            label: 'config.maxWithdrawAmountDaily2',
            prop: 'maxWithdrawAmountDaily2'
          },
          {
            type: 'show',
            label: 'config.maxWithdrawAmountDaily3',
            prop: 'maxWithdrawAmountDaily3'
          },
          {
            type: 'show',
            label: 'config.minDepositAmount',
            prop: 'minDepositAmount'
          },
          {
            type: 'show',
            label: 'config.maxWithdrawAmountApproval',
            prop: 'maxWithdrawAmountApproval'
          },
          {
            type: 'show',
            label: 'config.maxDepositAmountApproval',
            prop: 'maxDepositAmountApproval'
          },
          {
            type: 'show',
            label: 'config.displaySort',
            prop: 'displaySort'
          },
          {
            type: 'options',
            label: 'common.status',
            prop: 'status',
            options: CONFIG_ASSET_STATUS
          },
          {
            type: 'options',
            label: 'config.appVisible',
            prop: 'appVisible',
            options: CONFIG_ASSET_APP_VISIBLE
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
            label: 'order.lastUpdateTime',
            prop: 'lastUpdateTime'
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
