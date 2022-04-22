<template>
  <div id="adminList">
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
          :label="$t(`admin.adminName`)"
        >
          <el-input
            v-model="dialogForm.adminName"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`admin.newPW`)"
        >
          <el-input
            v-model="dialogForm.newPW"
            show-password
            autocomplete="off"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`admin.cfNewPW`)"
        >
          <el-input
            v-model="dialogForm.cfNewPW"
            show-password
            autocomplete="off"
            :readonly="dialogForm.id && !isEdit"
          />
        </el-form-item>
        <el-form-item
          :label="$t(`admin.roleName`)"
        >
          <el-select
            v-model="dialogForm.roleName"
            :disabled="dialogForm.id && !isEdit"
            placeholder=""
            style="width: 100%"
          >
            <el-option
              v-for="val in ['role1','role2','role3']"
              :key="val"
              :label="val"
              :value="val"
            />
          </el-select>
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
              v-for="(label, val) in MERCHANT_STATUS"
              :key="label"
              :label="$t(`options.${MERCHANT_STATUS[label]}`)"
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
import { MERCHANT_STATUS } from '@/utils/selection.const'
export default {
  name: 'AdminList',
  data() {
    return {
      MERCHANT_STATUS: MERCHANT_STATUS,
      filterFormVal: {
        columns: [
          {
            type: 'input',
            model: 'adminName',
            label: 'admin.adminName',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'roleName',
            label: 'admin.roleName',
            placeholder: ''
          },
          {
            type: 'select',
            model: 'status',
            label: 'common.status',
            options: MERCHANT_STATUS
          }
        ],
        default: {
          adminName: 'adminName',
          roleName: 'roleName',
          status: ''
        },
        addAction: this.addAction,
        searchAction: this.searchAction
      },
      tableData: {
        source: [
          {
            adminName: 'adminName',
            roleName: 'role1',
            status: 'normal',
            createTime: new Date(),
            creator: 'creator',
            updateTime: new Date(),
            lastLogin: new Date(),
            updator: 'updator',
            id: '12311',
            lastLoginIP: 'lastLoginIP',
            newPW: '123',
            cfNewPW: '123'
          }
        ],
        columns: [
          {
            type: 'show',
            label: 'admin.adminName',
            prop: 'adminName'
          },
          {
            type: 'show',
            label: 'admin.roleName',
            prop: 'roleName'
          },
          {
            type: 'options',
            label: 'common.status',
            prop: 'status',
            options: MERCHANT_STATUS
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
            type: 'time',
            label: 'user.lastLogin',
            prop: 'lastLogin'
          },
          {
            type: 'show',
            label: 'admin.lastLoginIP',
            prop: 'lastLoginIP'
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
    }
  }
}
</script>
