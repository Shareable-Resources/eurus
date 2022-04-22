<template>
  <div id="ConfigReportTime">
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
          :label="$t(`config.dailyReportStartTime`)"
        >
          <el-time-picker
            v-model="dialogForm.dailyReportStartTime"
            :readonly="dialogForm.id && !isEdit"
            value-format="HH:mm:ss"
            placeholder=""
          />
        </el-form-item>
        <el-form-item
          :label="$t(`config.dailyReportEndTime`)"
        >
          <el-time-picker
            v-model="dialogForm.dailyReportEndTime"
            readonly
            value-format="HH:mm:ss"
            placeholder=""
          />
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
export default {
  name: 'ConfigReportTime',
  data() {
    return {
      filterFormVal: {
        columns: [
          {
            type: 'input',
            model: 'appidOrMerchant',
            label: 'config.appidOrMerchant'
          }
        ],
        default: {
          appidOrMerchant: 'appidOrMerchant'
        },
        addAction: this.addAction,
        searchAction: this.searchAction
      },
      tableData: {
        source: [
          {
            appidOrMerchant: 'appidOrMerchant',
            dailyReportStartTime: '13:00:00',
            dailyReportEndTime: '12:59:59',
            updator: 'updator',
            creator: 'creator',
            createTime: new Date(),
            updateTime: new Date(),
            id: '1231231'
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
            label: 'config.dailyReportStartTime',
            prop: 'dailyReportStartTime'
          },
          {
            type: 'filter',
            label: 'config.dailyReportEndTime',
            prop: 'dailyReportStartTime',
            filterFn: (start) => {
              return '12:59:59'
            }
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
      console.log(this.dialogForm)
      this.dialogFormVisible = false
    }
  }
}
</script>
