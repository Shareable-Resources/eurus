<template>
  <div id="configMarketPrice">
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
          :label="$t(`config.trandingPair`)"
        >
          <el-select
            v-model="dialogForm.trandingPair"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_MKP_TRANS_PAIR"
              :key="label"
              :label="$t(`options.${label}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.basicAsset`)"
        >
          <el-input
            :value="dialogForm.trandingPair ? `${dialogForm.trandingPair.split('/')[0]}`: ''"
            readonly
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.quoteAsset`)"
        >
          <el-input
            :value="dialogForm.trandingPair ? `${dialogForm.trandingPair.split('/')[1]}`: ''"
            readonly
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.decimalDisplay`)"
        >
          <el-input
            v-model="dialogForm.decimalDisplay"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.feeType`)"
        >
          <el-select
            v-model="dialogForm.feeType"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_MKP_FEE_TYPE"
              :key="label"
              :label="$t(`options.${label}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.marketFeed`)"
        >
          <el-select
            v-model="dialogForm.marketFeed"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_MKP_SOURCE"
              :key="label"
              :label="$t(`options.${label}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.marketFeedPrice`)"
        >
          <el-input
            v-model="dialogForm.marketFeedPrice"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.genPrice`)"
        >
          <el-input
            v-model="dialogForm.genPrice"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.marketFeedStatus`)"
        >
          <el-select
            v-model="dialogForm.marketFeedStatus"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in CONFIG_MKP_SOURCE_STATUS"
              :key="label"
              :label="$t(`options.${CONFIG_MKP_SOURCE_STATUS[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`config.marketPriceDisplaySort`)"
        >
          <el-input
            v-model="dialogForm.marketPriceDisplaySort"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.marketQuoVisible`)"
        >
          <el-select
            v-model="dialogForm.marketQuoVisible"
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
import {
  CONFIG_MKP_FEE_TYPE,
  CONFIG_MKP_SOURCE_STATUS,
  CONFIG_FEE_STATUS,
  CONFIG_MKP_TRANS_PAIR,
  CONFIG_MKP_SOURCE
} from '@/utils/selection.const'
export default {
  name: 'ConfigMarketPrice',
  data() {
    return {
      CONFIG_MKP_FEE_TYPE: CONFIG_MKP_FEE_TYPE,
      CONFIG_MKP_SOURCE_STATUS: CONFIG_MKP_SOURCE_STATUS,
      CONFIG_FEE_STATUS: CONFIG_FEE_STATUS,
      CONFIG_MKP_TRANS_PAIR: CONFIG_MKP_TRANS_PAIR,
      CONFIG_MKP_SOURCE: CONFIG_MKP_SOURCE,
      filterFormVal: {
        columns: [
          {
            type: 'select',
            model: 'trandingPair',
            label: 'config.trandingPair',
            options: CONFIG_MKP_TRANS_PAIR
          },
          {
            type: 'select',
            model: 'marketFeed',
            label: 'config.marketFeed',
            options: CONFIG_MKP_SOURCE
          },
          {
            type: 'select',
            model: 'marketQuoVisible',
            label: 'config.marketQuoVisible',
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
            trandingPair: 'BTC/USDT',
            asset: 'asset',
            unit: '8',
            addrType: 'TYPE_ETH',
            creator: 'creator',
            lastUpdateTime: new Date(),
            updator: 'updator',
            id: '12311',
            status: 'listed',
            feeType: 'autoFee',
            marketFeedStatus: 'connected',
            marketQuoVisible: 'normal'
          }
        ],
        columns: [
          {
            type: 'show',
            label: 'config.trandingPair',
            prop: 'trandingPair'
          },
          {
            type: 'show',
            label: 'config.basicAsset',
            prop: 'basicAsset'
          },
          {
            type: 'show',
            label: 'config.quoteAsset',
            prop: 'quoteAsset'
          },
          {
            type: 'show',
            label: 'config.decimalDisplay',
            prop: 'decimalDisplay'
          },
          {
            type: 'options',
            label: 'config.feeType',
            prop: 'feeType',
            options: CONFIG_MKP_FEE_TYPE
          },
          {
            type: 'show',
            label: 'config.marketFeed',
            prop: 'marketFeed'
          },
          {
            type: 'show',
            label: 'config.marketFeedPrice',
            prop: 'marketFeedPrice'
          },
          {
            type: 'show',
            label: 'config.genPrice',
            prop: 'genPrice'
          },
          {
            type: 'options',
            label: 'config.marketFeedStatus',
            prop: 'marketFeedStatus',
            options: CONFIG_MKP_SOURCE_STATUS
          },
          {
            type: 'show',
            label: 'config.marketPriceDisplaySort',
            prop: 'marketPriceDisplaySort'
          },
          {
            type: 'options',
            label: 'config.marketQuoVisible',
            prop: 'marketQuoVisible',
            options: CONFIG_FEE_STATUS
          },
          {
            type: 'time',
            label: 'order.lastUpdateTime',
            prop: 'lastUpdateTime'
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
