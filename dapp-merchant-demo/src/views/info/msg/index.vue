<template>
  <div id="InfoMsg">
    <div class="create">
      <el-button
        size="small"
        type="primary"
        @click="addAction"
      >{{ $t(`common.create`) }}
      </el-button>
    </div>
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
          :label="$t(`info.msgContent`)"
        >
          <el-input
            v-model="dialogForm.msgContent"
            type="textarea"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.userNo`)"
        >
          <el-input
            v-model="dialogForm.userNo"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.startTime`)"
        >
          <el-date-picker
            v-model="dialogForm.startTime"
            :disabled="dialogForm.id && !isEdit"
            type="datetime"
            placeholder=""
          />
        </el-form-item>
        <el-form-item
          :label="$t(`common.endTime`)"
        >
          <el-date-picker
            v-model="dialogForm.endTime"
            :disabled="dialogForm.id && !isEdit"
            type="datetime"
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
  name: 'InfoMsg',
  data() {
    return {
      tableData: {
        source: [
          {
            msgContent: 'msgContent',
            userNo: 'userNo',
            sendTime: new Date(),
            updateTime: new Date(),
            updator: 'updator',
            id: 1233
          }
        ],
        columns: [
          {
            type: 'show',
            label: 'info.msgContent',
            prop: 'msgContent'
          },
          {
            type: 'show',
            label: 'common.userNo',
            prop: 'userNo'
          },
          {
            type: 'time',
            label: 'info.sendTime',
            prop: 'sendTime'
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
            width: '200',
            btns: [
              {
                type: 'details',
                fn: this.handleClick
              },
              {
                type: 'btn',
                label: 'common.del',
                fn: this.handleDel
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
    addAction() {
      this.dialogForm = {}
      this.dialogFormVisible = true
    },
    handleSizeChange(val) {
      console.log(`每页 ${val} 条`)
    },
    handleCurrentChange(val) {
      console.log(`当前页: ${val}`)
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
    handleDel(row) {
      console.log(row, 'del')
    }
  }
}
</script>

<style lang="scss" scoped>
.create {
  margin-bottom: 10px;
}
</style>
