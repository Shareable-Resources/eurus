import Vue from 'vue'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import App from './App.vue';
import "@/style/index.scss"
import router from './router'
import Moment from 'vue-moment'
import axios from 'axios'
import VueAxios from 'vue-axios'
Vue.config.productionTip = false
Vue.use(ElementUI);
Vue.use(Moment);
Vue.use(VueAxios, axios);
new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
