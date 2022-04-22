<template>
  <div id="role">
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
          :label="$t(`admin.roleName`)"
        >
          <el-input
            v-model="dialogForm.roleName"
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
              v-for="(label, val) in MERCHANT_STATUS"
              :key="label"
              :label="$t(`options.${MERCHANT_STATUS[label]}`)"
              :value="val"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t(`admin.rights`)">
          <div>{{ `模拟數據，真實數據等接API後替換` }}</div>
          <el-tree
            :data="treeTemp"
            node-key="id"
            default-expand-all
            :expand-on-click-node="false"
          >
            <span slot-scope="{ node, data }" class="custom-tree-node">
              <span>{{ node.label }}</span>
              <span v-if="node.isLeaf">
                <el-checkbox
                  v-show="!dialogForm.id || isEdit"
                  v-model="data.read"
                  :label="$t(`common.check`)"
                  :disabled="!data.hasOwnProperty('read')"
                />
                <el-checkbox
                  v-show="!dialogForm.id || isEdit"
                  v-model="data.create"
                  :label="$t(`common.create`)"
                  :disabled="!data.hasOwnProperty('create')"
                />
                <el-checkbox
                  v-show="!dialogForm.id || isEdit"
                  v-model="data.edit"
                  :label="$t(`common.edit`)"
                  :disabled="!data.hasOwnProperty('edit')"
                />
                <el-checkbox
                  v-show="!dialogForm.id || isEdit"
                  v-model="data.audit"
                  :label="$t(`common.audit`)"
                  :disabled="!data.hasOwnProperty('audit')"
                />
              </span>
            </span>
          </el-tree>
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
  name: 'Role',
  data() {
    return {
      MERCHANT_STATUS: MERCHANT_STATUS,
      filterFormVal: {
        columns: [
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
          roleName: 'roleName',
          status: ''
        },
        addAction: this.addAction,
        searchAction: this.searchAction
      },
      tableData: {
        source: [
          {
            roleName: 'roleName',
            status: 'normal',
            createTime: new Date(),
            creator: 'creator',
            updateTime: new Date(),
            updator: 'updator',
            id: '12311'
          }
        ],
        columns: [
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
      treeTemp: [
        {
          id: 1,
          label: '餘額管理',
          children: [
            {
              id: 2,
              label: '用戶餘額列表',
              read: false,
              create: false,
              edit: false
            },
            {
              id: 3,
              label: '錢包餘額列表',
              read: false,
              create: false,
              edit: false
            },
            {
              id: 4,
              label: '用戶地扯餘額列表',
              read: false,
              create: false,
              edit: false,
              audit: false
            }
          ]
        },
        {
          id: 5,
          label: '餘額管理2',
          children: [
            {
              id: 6,
              label: '用戶餘額列表2',
              read: false
            },
            {
              id: 7,
              label: '錢包餘額列表2',
              read: false,
              create: false,
              edit: false
            }
          ]
        }
      ]
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
