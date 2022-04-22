import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    redirect: '/transactions',
    children: [
      {
        path: 'transactions',
        component: () => import('@/views/Home.vue'),
        name: 'transaction',
        meta: { title: 'dailyBill', icon: 'list', noCache: true }
      },
      {
        path: 'transactions/:id',
        component: () => import('@/views/Home.vue'),
        name: 'transactions',
        meta: { title: 'dailyBill', icon: 'list', noCache: true }
      }
    ]
  },
  {
    path: '/blocks/:id',
    name: 'Blocks',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/Block')
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  }
]

const router = new VueRouter({
  routes
})

export default router
