import Vue from 'vue'
import VueRouter from 'vue-router'

import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'

import NotFound from '@/components/NotFound.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '*',
    name: 'NotFound',
    component: NotFound
  },
  {
    path: '/',
    name: 'login',
    component: Login,
    meta: {
      auth: false
    }
  },
  {
    path: '/home',
    name: 'home',
    component: Home,
    meta: {
      auth: false
    }
  },
]

const router = new VueRouter({
  mode: 'history',
  routes
})

export default router
