<template>
  <div class="select-time-zone flex">
    <div class="avatar-wrapper pointer" @click="showDialog">
      {{ $t(`common.tz`)+': '+TimeZoneList[tz] }}
    </div>
    <div class="avatar-wrapper pl20">
      {{ date | moment('timezone', tz, "YYYY-MM-DD HH:mm") }}
    </div>
    <el-dialog
      :title="$t(`common.tz`)"
      :visible.sync="dialogVisible"
      width="30%"
    >
      <el-select v-model="selectedTZ">
        <el-option
          v-for="(val,key) in TimeZoneList"
          :key="key"
          :value="key"
          :label="val"
        />
      </el-select>
      <span slot="footer" class="dialog-footer">
        <common-btn
          :btn-data="cancel"
        />
        <common-btn
          :btn-data="confirm"
        />
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import TimeZoneList from './time-zone.json'
export default {
  name: 'SelectTimeZone',
  data() {
    return {
      dialogVisible: false,
      selectedTZ: '',
      TimeZoneList: TimeZoneList,
      confirm: {
        type: 'success',
        label: 'common.confirm',
        action: this.submitAction
      },
      cancel: {
        type: 'info',
        label: 'common.cancel',
        action: this.cancelAction
      },
      date: new Date(),
      dateTimer: ''
    }
  },
  computed: {
    ...mapGetters([
      'tz'
    ])
  },
  mounted() {
    this.setTimer()
  },
  beforeDestroy() {
    clearInterval(this.dateTimer)
  },
  methods: {
    submitAction() {
      const { selectedTZ } = this
      this.$store.dispatch('user/changeTZ', selectedTZ)
      this.dialogVisible = false
    },
    cancelAction() {
      this.dialogVisible = false
    },
    showDialog() {
      this.selectedTZ = this.tz
      this.dialogVisible = true
    },
    setTimer() {
      const that = this
      const INTERVAL = 60 * 60
      this.dateTimer = setInterval(() => {
        that.date = new Date()
      }, INTERVAL)
    }
  }
}
</script>

<style lang="scss" scoped>
.flex {
 display: flex;
}
.pl20 {
 padding-left: 20px;
}
</style>
