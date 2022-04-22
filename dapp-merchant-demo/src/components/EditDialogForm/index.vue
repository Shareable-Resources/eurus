<template>
  <el-dialog
    class="edit-dialog-form"
    :close-on-click-modal="false"
    :width="dialogData.width || '80%'"
    :title="dialogData.title ? $t(dialogData.title) : $t(`common.details`)"
    :visible.sync="dialogData.visible"
    :before-close="dialogData.beforeClose ? dialogData.beforeClose : null"
  >
    <el-form
      :model="dialogData.source"
    >
      <el-form-item
        v-for="item in dialogData.form"
        :key="item.prop"
        :label="$t(item.label)"
      >
        <el-input
          v-if="item.type==='input'"
          v-model="dialogData.source[item.prop]"
          :readonly="item.readonly || true"
          :placeholder="item.placeholder || ''"
        />

        <el-select
          v-if="item.type==='select'"
          v-model="dialogData.source[item.prop]"
          :disabled="item.disabled || true"
          :placeholder="item.placeholder || ''"
          style="width: 100%"
        >
          <el-option
            v-for="(label, val) in item.options"
            :key="label"
            :label="$t(`options.${item.options[label]}`)"
            :value="val"
          />
        </el-select>

        <el-input
          v-if="item.type==='time'"
          :readonly="item.readonly || true"
          :value="dialogData.source[item.prop] | moment('timezone', tz, 'YYYY-MM-DD HH:mm:ss')"
        />

        <el-input
          v-if="item.type==='filter'"
          :readonly="item.readonly || true"
          :value="item.filterFn(dialogData.source[item.prop])"
        />
      </el-form-item>
      <el-form-item
        v-if="dialogData.refundAction"
        :label="$t(`common.submitFinancePw`)"
      >
        <el-input
          v-model="submitFinancePw"
          show-password
        />
      </el-form-item>
      <el-form-item
        v-if="dialogData.auditAction"
        :label="$t(`common.auditorFinancePw`)"
      >
        <el-input
          v-model="auditorFinancePw"
          show-password
        />
      </el-form-item>
    </el-form>

    <div slot="footer" class="dialog-footer">
      <el-button
        v-if="dialogData.refundAction"
        @click="dialogData.refundAction(submitFinancePw)"
      >{{ $t(`common.refund`) }}
      </el-button>
      <el-button
        v-if="dialogData.auditAction"
        @click="dialogData.auditAction(auditorFinancePw)"
      >{{ $t(`common.audit`) }}
      </el-button>
      <el-button
        v-if="dialogData.cancelAction"
        @click="dialogData.cancelAction"
      >{{ $t(`common.cancel`) }}
      </el-button>
      <el-button
        v-if="dialogData.rejectAction"
        @click="dialogData.rejectAction"
      >{{ $t(`common.reject`) }}
      </el-button>
      <el-button
        v-if="dialogData.approveAction"
        @click="dialogData.approveAction"
      >{{ $t(`common.approve`) }}
      </el-button>
    </div>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex'
export default {
  name: 'EditDialogForm',
  props: {
    dialogData: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      submitFinancePw: '',
      auditorFinancePw: ''
    }
  },
  computed: {
    ...mapGetters([
      'tz'
    ])
  },
  methods: {
    handleClose() {
      const { beforeClose } = this.dialogData.beforeClose
      if (beforeClose) {
        return beforeClose()
      }
    }
  }
}
</script>
