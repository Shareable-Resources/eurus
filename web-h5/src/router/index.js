import Vue from 'vue'
import VueRouter from 'vue-router'
import layout from '@/layout/layout'

Vue.use(VueRouter)

export const routes = [
  {
     path: '/', 
     name: 'home',
     component: layout,
     children: [
  {
    path: '/dashboard',
    name: 'Dashboard',
    title: 'dashboard',
    component: () => import('../views/dashboard')
  },
  {
    path : '/login',
    name : 'Login',
    title: 'login',
    component: () => import('../views/login')
  },
  {
    path : '/cam',
    name : 'cam',
    component: () => import('../views/cam')
  }

]
  }]

const router = new VueRouter({
  mode: 'history',
  // base: process.env.BASE_URL,
  routes : routes
})

router.beforeEach((to, from, next) => {
  next(vm => {
      vm.$router.replace(from.path);
  })
 })

export default router
