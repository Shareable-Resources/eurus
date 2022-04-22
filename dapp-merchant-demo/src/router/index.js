import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/* Router Modules */
// import componentsRouter from './modules/components'
// import chartsRouter from './modules/charts'
// import tableRouter from './modules/table'
// import nestedRouter from './modules/nested'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index')
      }
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/auth-redirect',
    component: () => import('@/views/login/auth-redirect'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error-page/401'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/balance',
    children: [
      {
        path: 'balance',
        component: () => import('@/views/demo/balance'),
        name: 'Balance',
        meta: { title: 'balance', icon: 'chart', affix: true }
      }
    ]
  },
  {
    path: '/deposit',
    component: Layout,
    redirect: '/deposit',
    children: [
      {
        path: '/deposit',
        component: () => import('@/views/demo/deposit'),
        name: 'Deposit',
        meta: { title: 'deposit', icon: 'chart', noCache: true }
      }
    ]
  },
  {
    path: '/payment',
    component: Layout,
    redirect: '/payment',
    children: [
      {
        path: '/payment',
        component: () => import('@/views/demo/payment'),
        name: 'Payment',
        meta: { title: 'payment', icon: 'chart', noCache: true }
      }
    ]
  },
  {
    path: '/refund',
    component: Layout,
    redirect: '/refund',
    children: [
      {
        path: '/refund',
        component: () => import('@/views/demo/refund'),
        name: 'Refund',
        meta: { title: 'refund', icon: 'chart', noCache: true }
      }
    ]
  },
  // {
  //   path: '/',
  //   component: Layout,
  //   redirect: '/dashboard',
  //   children: [
  //     {
  //       path: 'dashboard',
  //       component: () => import('@/views/dashboard/index'),
  //       name: 'Dashboard',
  //       meta: { title: 'dashboard', icon: 'dashboard', affix: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/balance',
  //   name: 'Balance',
  //   meta: {
  //     title: 'balance',
  //     icon: 'chart'
  //   },
  //   component: Layout,
  //   redirect: '/balance/user',
  //   children: [
  //     {
  //       path: 'user',
  //       component: () => import('@/views/balance/user/index'),
  //       name: 'UserBalance',
  //       meta: { title: 'userBalance', icon: 'chart', noCache: true }
  //     },
  //     {
  //       path: 'wallet',
  //       component: () => import('@/views/balance/wallet/index'),
  //       name: 'WalletBalance',
  //       meta: { title: 'walletBalance', icon: 'chart', noCache: true }
  //     },
  //     {
  //       path: 'user-addr',
  //       component: () => import('@/views/balance/user-addr/index'),
  //       name: 'UserAddrBalance',
  //       meta: { title: 'userAddrBalance', icon: 'chart', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/user',
  //   component: Layout,
  //   redirect: '/user/index',
  //   children: [
  //     {
  //       path: 'index',
  //       component: () => import('@/views/user/index'),
  //       name: 'User',
  //       meta: { title: 'user', icon: 'user', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/merchant',
  //   component: Layout,
  //   redirect: '/merchant/index',
  //   children: [
  //     {
  //       path: 'index',
  //       component: () => import('@/views/merchant/index'),
  //       name: 'Merchant',
  //       meta: { title: 'merchant', icon: 'nested', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/order',
  //   name: 'Order',
  //   meta: {
  //     title: 'order',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/order/deposit',
  //   children: [
  //     {
  //       path: 'deposit',
  //       component: () => import('@/views/order/deposit/index'),
  //       name: 'Deposit',
  //       meta: { title: 'deposit', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'withdraw',
  //       component: () => import('@/views/order/withdraw/index'),
  //       name: 'Withdraw',
  //       meta: { title: 'withdraw', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'payment',
  //       component: () => import('@/views/order/payment/index'),
  //       name: 'Payment',
  //       meta: { title: 'payment', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'transfer',
  //       component: () => import('@/views/order/transfer/index'),
  //       name: 'Transfer',
  //       meta: { title: 'transfer', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'transaction',
  //       component: () => import('@/views/order/transaction/index'),
  //       name: 'Transaction',
  //       meta: { title: 'transaction', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'callback',
  //       component: () => import('@/views/order/callback/index'),
  //       name: 'Callback',
  //       meta: { title: 'callback', icon: 'list', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/bill',
  //   name: 'Bill',
  //   meta: {
  //     title: 'bill',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/bill/daily',
  //   children: [
  //     {
  //       path: 'daily',
  //       component: () => import('@/views/bill/daily/index'),
  //       name: 'Daily',
  //       meta: { title: 'dailyBill', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'monthly',
  //       component: () => import('@/views/bill/monthly/index'),
  //       name: 'Monthly',
  //       meta: { title: 'monthlyBill', icon: 'list', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/stats',
  //   name: 'Stats',
  //   meta: {
  //     title: 'statistics',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/stats/user',
  //   children: [
  //     {
  //       path: 'user',
  //       component: () => import('@/views/stats/user/index'),
  //       name: 'UserSta',
  //       meta: { title: 'userSta', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'merchant',
  //       component: () => import('@/views/stats/merchant/index'),
  //       name: 'MerchantSta',
  //       meta: { title: 'merchantSta', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'order',
  //       component: () => import('@/views/stats/order/index'),
  //       name: 'OrderSta',
  //       meta: { title: 'orderSta', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'user-addr',
  //       component: () => import('@/views/stats/user-addr/index'),
  //       name: 'UserAddrSta',
  //       meta: { title: 'userAddrSta', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'public-chain',
  //       component: () => import('@/views/stats/public-chain/index'),
  //       name: 'PublicChainSta',
  //       meta: { title: 'publicChainSta', icon: 'list', noCache: true }
  //     },
  //     {
  //       path: 'private-chain',
  //       component: () => import('@/views/stats/private-chain/index'),
  //       name: 'PrivateChainSta',
  //       meta: { title: 'privateChainSta', icon: 'list', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/admin',
  //   name: 'Admin',
  //   meta: {
  //     title: 'admin',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/admin/role',
  //   children: [
  //     {
  //       path: 'role',
  //       component: () => import('@/views/admin/role'),
  //       name: 'Role',
  //       meta: { title: 'role', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'admin-list',
  //       component: () => import('@/views/admin/admin-list'),
  //       name: 'AdminList',
  //       meta: { title: 'adminList', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'change-pw',
  //       component: () => import('@/views/admin/change-pw'),
  //       name: 'ChangePW',
  //       meta: { title: 'changePW', icon: 'nested', noCache: true }
  //     }
  //   ]
  // },
  // not in M2
  // {
  //   path: '/wallet',
  //   name: 'Wallet',
  //   meta: {
  //     title: 'wallet',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/wallet/group-wallet',
  //   children: [
  //     {
  //       path: 'group-wallet',
  //       component: () => import('@/views/wallet/group-wallet'),
  //       name: 'GroupWallet',
  //       meta: { title: 'groupWallet', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'export-wallet',
  //       component: () => import('@/views/wallet/export-wallet'),
  //       name: 'ExportWallet',
  //       meta: { title: 'exportWallet', icon: 'nested', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/audit',
  //   name: 'Audit',
  //   meta: {
  //     title: 'audit',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/audit/kyc',
  //   children: [
  //     {
  //       path: 'kyc',
  //       component: () => import('@/views/audit/kyc'),
  //       name: 'KycAudit',
  //       meta: { title: 'kycAudit', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'order',
  //       component: () => import('@/views/audit/order'),
  //       name: 'OrderAudit',
  //       meta: { title: 'orderAudit', icon: 'nested', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/config',
  //   name: 'Config',
  //   meta: {
  //     title: 'config',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/config/asset',
  //   children: [
  //     {
  //       path: 'asset',
  //       component: () => import('@/views/config/asset'),
  //       name: 'AssetConfig',
  //       meta: { title: 'assetConfig', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'fee',
  //       component: () => import('@/views/config/fee'),
  //       name: 'FeeConfig',
  //       meta: { title: 'feeConfig', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'market-price',
  //       component: () => import('@/views/config/market-price'),
  //       name: 'MarketPriceConfig',
  //       meta: { title: 'marketPriceConfig', icon: 'nested', noCache: true }
  //     }
  //     // not in M2
  //     // {
  //     //   path: 'report-time',
  //     //   component: () => import('@/views/config/report-time'),
  //     //   name: 'ReportTimeConfig',
  //     //   meta: { title: 'reportTimeConfig', icon: 'nested', noCache: true }
  //     // }
  //   ]
  // },
  // not in M2
  // {
  //   path: '/alert',
  //   name: 'Alert',
  //   meta: {
  //     title: 'alertNotification',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/alert/order',
  //   children: [
  //     {
  //       path: 'order',
  //       component: () => import('@/views/alert/order-audit'),
  //       name: 'OrderAduitAlert',
  //       meta: { title: 'orderAduitAlert', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'wallet',
  //       component: () => import('@/views/alert/wallet-bal'),
  //       name: 'WalletBalAlert',
  //       meta: { title: 'walletBalAlert', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'export',
  //       component: () => import('@/views/alert/export-status'),
  //       name: 'ExportStatusAlert',
  //       meta: { title: 'exportStatusAlert', icon: 'nested', noCache: true }
  //     }
  //   ]
  // },
  // {
  //   path: '/log',
  //   name: 'Log',
  //   meta: {
  //     title: 'log',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/log/user',
  //   children: [
  //     {
  //       path: 'user',
  //       component: () => import('@/views/log/user'),
  //       name: 'UserLoginLog',
  //       meta: { title: 'userLoginLog', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'admin',
  //       component: () => import('@/views/log/admin'),
  //       name: 'AdminOperationLog',
  //       meta: { title: 'adminOperationLog', icon: 'nested', noCache: true }
  //     }
  //   ]
  // },
  // not in M2
  // {
  //   path: '/notice',
  //   name: 'notice',
  //   meta: {
  //     title: 'information',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/notice/ann',
  //   children: [
  //     {
  //       path: 'ann',
  //       component: () => import('@/views/notice/announcement'),
  //       name: 'Announcement',
  //       meta: { title: 'announcement', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'msg',
  //       component: () => import('@/views/notice/message'),
  //       name: 'MsgNotification',
  //       meta: { title: 'msgNotification', icon: 'nested', noCache: true }
  //     }
  //   ]
  // },
  // not in M2
  // {
  //   path: '/manual',
  //   name: 'ManualOrder',
  //   meta: {
  //     title: 'manualOrder',
  //     icon: 'list'
  //   },
  //   component: Layout,
  //   redirect: '/manual/create',
  //   children: [
  //     {
  //       path: 'create',
  //       component: () => import('@/views/manual/create'),
  //       name: 'CreateManualOrder',
  //       meta: { title: 'createManualOrder', icon: 'nested', noCache: true }
  //     },
  //     {
  //       path: 'audit',
  //       component: () => import('@/views/manual/audit'),
  //       name: 'AuditManualOrder',
  //       meta: { title: 'auditManualOrder', icon: 'nested', noCache: true }
  //     }
  //   ]
  // },
  {
    path: '/profile',
    component: Layout,
    redirect: '/profile/index',
    hidden: true,
    children: [
      {
        path: 'index',
        component: () => import('@/views/profile/index'),
        name: 'Profile',
        meta: { title: 'profile', icon: 'user', noCache: true }
      }
    ]
  }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  /** when your routing map is too long, you can split it into small modules **/
  // componentsRouter,
  // chartsRouter,
  // nestedRouter,
  // tableRouter,
  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}
export default router
