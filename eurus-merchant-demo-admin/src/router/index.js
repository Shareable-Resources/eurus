import Vue from 'vue'
import VueRouter from 'vue-router'

import Login from '@/views/Login.vue'
import DepositList from '@/views/DepositList.vue'
import UserList from '@/views/UserList.vue'
import WithdrawList from '@/views/WithdrawList.vue'
import WithdrawHistory from '@/views/WithdrawHistory.vue'

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
    path: '/depositList',
    name: 'depositList',
    component: DepositList,
    meta: {
      auth: false
    }
  },
  {
    path: '/userList',
    name: 'userList',
    component: UserList,
    meta: {
      auth: false
    }
  },
  {
    path: '/withdrawList',
    name: 'withdrawList',
    component: WithdrawList,
    meta: {
      auth: false
    }
  },
  {
    path: '/withdrawHistory',
    name: 'withdrawHistory',
    component: WithdrawHistory,
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
