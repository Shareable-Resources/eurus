<template>
  <div id="DailyBill">
    <filter-form
      :filter-form="filterFormVal"
    />
  </div>
</template>

<script>
export default {
  name: 'DailyBill',
  data() {
    return {
      filterFormVal: {
        columns: [
          {
            type: 'input',
            model: 'userNo',
            label: 'common.userNo',
            placeholder: ''
          },
          {
            type: 'input',
            model: 'appidOrMerchant',
            label: 'config.appidOrMerchant'
          },
          {
            type: 'date',
            model: 'startDate',
            label: 'common.startDate',
            format: 'yyyy/MM/dd'
          },
          {
            type: 'date',
            model: 'endDate',
            label: 'common.endDate',
            format: 'yyyy/MM/dd'
          }
        ],
        default: {
          userNo: '',
          appidOrMerchant: '',
          startDate: '',
          endDate: ''
        },
        searchAction: this.searchAction,
        exportAction: this.exportAction
      }
    }
  },
  methods: {
    searchAction(data) {
      console.log(data, 'search')
    },
    exportAction(data) {
      console.log(data, 'export')
    }
  }
}
</script>
