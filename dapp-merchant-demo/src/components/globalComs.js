import RealTime from '@/components/RealTime'
import DailyReportExport from '@/components/DailyReportExport'
import FilterForm from '@/components/FilterForm'
import BasicTable from '@/components/BasicTable'
import BasicDialogForm from '@/components/BasicDialogForm'
import CommonBtn from '@/components/CommonBtn'

const plugins = {
  install(Vue) {
    Vue.component('RealTime', RealTime)
    Vue.component('DailyReportExport', DailyReportExport)
    Vue.component('FilterForm', FilterForm)
    Vue.component('BasicTable', BasicTable)
    Vue.component('BasicDialogForm', BasicDialogForm)
    Vue.component('CommonBtn', CommonBtn)
  }
}
export default plugins
