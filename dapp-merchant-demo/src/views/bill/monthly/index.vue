<template>
  <div id="MonthlyBill">
    <filter-form
      :filter-form="filterFormVal"
    />
  </div>
</template>

<script>
export default {
  name: 'MonthlyBill',
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
            type: 'month',
            model: 'selectMonth',
            label: 'common.selectMonth',
            format: 'yyyy/MM'
          }
        ],
        default: {
          userNo: '',
          appidOrMerchant: '',
          selectMonth: ''
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
