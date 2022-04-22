import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import i18n from './i18n'
import VueCookies from 'vue-cookies'


import './quasar'
import globalCompoent from '@/components/globalCom'

Vue.config.productionTip = false
Vue.config.devtools = true
Vue.use(globalCompoent);
Vue.use(VueCookies)

var app = new Vue({
  i18n,
  router,
  store,
  render: h => h(App)
})
app.$mount('#app')
