<template>
  <div id="kyc">
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
      :title="dialogForm.userNo ? $t(`common.details`):$t(`common.create`)"
      :visible.sync="dialogFormVisible"
      :before-close="handleClose"
    >
      <el-form
        label-position="left"
        :model="dialogForm"
      >
        <el-form-item
          :label="$t(`common.userNo`)"
        >
          <el-input
            v-model="dialogForm.userNo"
            readonly
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.email`)"
        >
          <el-input
            v-model="dialogForm.email"
            readonly
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.status`)"
        >
          <el-select
            v-model="dialogForm.status"
            :disabled="dialogForm.userNo && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="(label, val) in USER_STATUS"
              :key="`status_${label}`"
              :label="$t(`options.${USER_STATUS[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t(`user.kycStatus`)"
        >
          <el-input
            readonly
            :value="dialogForm.kycStatus? $t(`options.${USER_KYCSTATUS[dialogForm.kycStatus]}`) : ''"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.registerFrom`)"
        >
          <el-input
            v-model="dialogForm.registerFrom"
            readonly
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.userGroup`)"
        >
          <el-input
            readonly
            :value="dialogForm.userGourp? $t(`options.${USER_GROUP[dialogForm.userGroup]}`): ''"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`user.inviteCode`)"
        >
          <el-input
            v-model="dialogForm.inviteCode"
            readonly
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.registerDate`)"
        >
          <el-input
            readonly
            :value="dialogForm.registerDate | moment('timezone', tz, 'YYYY-MM-DD HH:mm:ss')"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`user.lastLogin`)"
        >
          <el-input
            readonly
            :value="dialogForm.lastLogin | moment('timezone', tz, 'YYYY-MM-DD HH:mm:ss')"
          />
        </el-form-item>
      </el-form>
      <el-table
        border
        :data="gridData"
      >
        <el-table-column
          property="asset"
          :label="$t(`common.asset`)"
        />
        <el-table-column
          property="name"
          :label="$t(`common.available`)"
        />
        <el-table-column
          property="address"
          :label="$t(`common.frozen`)"
        />
        <el-table-column
          property="address"
          :label="$t(`common.total`)"
        />
        <el-table-column
          property="address"
          :label="$t(`common.addrType`)"
        >
          <template slot-scope="props">
            <span>{{ $t(`options.${WALLET_BAL_ADDRTYPE[props.row.addrType]}`) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          property="address"
          :label="$t(`common.addr`)"
        />
        <el-table-column
          align="center"
        >
          <template slot="header" slot-scope="scope">
            <el-button
              size="mini"
              type="danger"
              :disabled="dialogForm.userNo && !isEdit"
              @click="handleGenAddr(scope.$index, scope.row)"
            >{{ $t(`user.genAddr`) }}</el-button>
          </template>
          <template>
            {{ `` }}
          </template>
        </el-table-column>
      </el-table>
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
import { mapGetters } from 'vuex'
import { USER_STATUS, USER_KYCSTATUS, USER_GROUP } from '@/utils/selection.const'
export default {
  name: 'KYC',
  data() {
    return {
      USER_STATUS: USER_STATUS,
      USER_KYCSTATUS: USER_KYCSTATUS,
      USER_GROUP: USER_GROUP,
      filterFormVal: {
        columns: [
          {
            type: 'input',
            model: 'userNo',
            label: 'common.userNo',
            placeholder: ''
          },
          {
            type: 'select',
            model: 'status',
            label: 'common.status',
            options: USER_STATUS
          },
          {
            type: 'select',
            model: 'kycStatus',
            label: 'user.kycStatus',
            options: USER_KYCSTATUS
          },
          {
            type: 'select',
            model: 'kycLevel',
            label: 'user.kycLevel',
            options: USER_KYCSTATUS
          },
          {
            type: 'datetimerange',
            model: 'registerDate',
            label: 'common.registerDate',
            format: 'yyyy/MM/dd HH:mm:ss'
          },
          {
            type: 'datetimerange',
            model: 'kycSubmitTime',
            label: 'user.kycSubmitTime',
            format: 'yyyy/MM/dd HH:mm:ss'
          },
          {
            type: 'datetimerange',
            model: 'kycApproveTime',
            label: 'user.kycApproveTime',
            format: 'yyyy/MM/dd HH:mm:ss'
          }
        ],
        default: {
          userNo: 'userNo',
          email: '',
          registerDate: '',
          status: '',
          kycStatus: '',
          inviteCode: ''
        },
        searchAction: this.searchAction,
        exportAction: this.exportAction
      },
      tableData: {
        source: [
          {
            userNo: 'userNo',
            status: 'normal',
            kycStatus: 'pending',
            registerFrom: 'registerFrom',
            userGroup: 'merchant',
            inviteCode: 'inviteCode',
            registerDate: new Date(),
            lastLogin: new Date()
          }
        ],
        columns: [
          {
            type: 'show',
            label: 'common.userNo',
            prop: 'userNo'
          },
          {
            type: 'options',
            label: 'common.status',
            prop: 'status',
            options: USER_STATUS
          },
          {
            type: 'options',
            label: 'user.kycStatus',
            prop: 'kycStatus',
            options: USER_KYCSTATUS
          },
          {
            type: 'show',
            label: 'user.kycLevel',
            prop: 'kycLevel'
          },
          {
            type: 'options',
            label: 'common.userGroup',
            prop: 'userGroup',
            options: USER_GROUP
          },
          {
            type: 'show',
            label: 'common.registerFrom',
            prop: 'registerFrom'
          },
          {
            type: 'time',
            label: 'common.registerDate',
            prop: 'registerDate'
          },
          {
            type: 'time',
            label: 'user.kycSubmitTime',
            prop: 'kycSubmitTime'
          },
          {
            type: 'time',
            label: 'user.kycApproveTime',
            prop: 'kycApproveTime'
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
      isEdit: false,
      gridData: []
    }
  },
  computed: {
    ...mapGetters([
      'tz'
    ])
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
    handleClick(row) {
      this.isEdit = false
      this.dialogForm = Object.assign({}, row)
      this.dialogFormVisible = true
    },
    handleGenAddr() {
      const that = this
      console.log('gen Addr')
      this.$prompt(this.$t(`user.genAddrTip`), this.$t(`common.tip`), {
        confirmButtonText: that.$t(`common.confirm`),
        cancelButtonText: that.$t(`common.cancel`)
      }).then(({ value }) => {
        that.$message({
          type: 'success',
          message: 'Asset of Address: ' + value
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: 'cancel'
        })
      })
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

<style lang="scss" scoped>
</style>
