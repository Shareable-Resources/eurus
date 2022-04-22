<template>
  <div>
    <el-form
      :inline="true"
      size="small"
      :model="filterFormData"
      class="demo-form-inline"
    >
      <el-form-item
        :label="$t(`filter.dailyReport`)"
      >
        <el-date-picker
          v-model="dailyReportDate"
          type="date"
          format="yyyy/MM/dd"
          value-format="yyyy/MM/dd"
        />
        <el-button
          type="primary"
          icon="el-icon-download"
          class="download-icon"
          circle
          @click="emitExportFn"
        />
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: 'DailyReportExport',
  props: {
    dailyReportExportAction: {
      type: Function,
      required: true
    }
  },
  data() {
    return {
      filterFormData: {},
      dailyReportDate: ''
    }
  },
  methods: {
    emitExportFn() {
      return this.dailyReportExportAction(this.dailyReportDate)
    }
  }
}
</script>

<style lang="scss" scoped>
.download-icon {
  margin-left: 5px;
}
</style>
