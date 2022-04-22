<template>
  <div>
    <el-form
      :inline="true"
      size="small"
      :model="filterFormData"
      class="demo-form-inline"
    >
      <el-form-item
        :label="$t(`filter.realTime`)"
        class="form-item"
      >
        <el-input
          readonly
          prefix-icon="el-icon-date"
          :value="filterFormData.realTimeVal | moment('timezone', tz, 'YYYY-MM-DD HH:mm:ss')"
        />
        <el-button
          type="danger"
          icon="el-icon-refresh"
          class="refresh-icon"
          circle
          @click="onRefresh"
        />
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  name: 'RealTime',
  props: {
    realTimeRefreshAction: {
      type: Function,
      required: true
    }
  },
  data() {
    return {
      filterFormData: {
        realTimeVal: this.getTime()
      }
    }
  },
  computed: {
    ...mapGetters([
      'tz'
    ])
  },
  methods: {
    getTime() {
      return new Date().getTime()
    },
    onRefresh() {
      const newTime = this.getTime()
      this.$set(this.filterFormData, 'realTimeVal', newTime)
      return this.realTimeRefreshAction(newTime)
    }
  }
}
</script>

<style lang="scss" scoped>
.form-item {
  position: relative;
}
.refresh-icon {
  margin-left: 5px;
  position: absolute;
}
</style>
