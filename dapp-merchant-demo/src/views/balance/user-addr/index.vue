<template>
  <div id="user-addr">
    <div class="filter-form">
      <real-time
        :real-time-refresh-action="refreshRealTime"
      />
      <daily-report-export
        :daily-report-export-action="exportDailyReport"
      />
      <filter-form
        :filter-form="filterFormVal"
      />
    </div>
  </div>
</template>

<script>
import { STATS_USER_ADDR_TYPE } from '@/utils/selection.const'
export default {
  data() {
    return {
      STATS_USER_ADDR_TYPE: STATS_USER_ADDR_TYPE,
      filterFormVal: {
        columns: [
          // {
          //   type: 'input',
          //   model: 'testInp',
          //   label: 'filter.snapShot',
          //   placeholder: ''
          // },
          // {
          //   type: 'datetime',
          //   model: 'snapShot',
          //   label: 'filter.snapShot',
          //   format: 'yyyy/MM/dd HH:mm:ss'
          // },
          {
            type: 'snapShot',
            model: 'snapShot',
            label: 'filter.snapShot'
          },
          {
            type: 'select',
            model: 'addrType',
            label: 'common.addrType',
            options: STATS_USER_ADDR_TYPE
          },
          {
            type: 'input',
            model: 'designatedAddr',
            label: 'walletBal.designatedAddr',
            placeholder: ''
          },
          {
            type: 'select',
            model: 'asset',
            label: 'common.asset',
            options: STATS_USER_ADDR_TYPE
          }
        ],
        default: {
          testInp: 'testinp',
          snapShot: '',
          group: 'merchant'
        },
        searchAction: this.searchAction,
        exportAction: this.exportAction
      }
    }
  },
  methods: {
    refreshRealTime(timestamp) {
      console.log(timestamp)
    },
    exportDailyReport(date) {
      console.log('daily report', date)
    },
    searchAction(data) {
      console.log('search', data)
    },
    exportAction(data) {
      console.log('export', data)
    }
  }
}
</script>
