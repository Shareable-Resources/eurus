import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/admin/login',
    method: 'post',
    data
  })
}

export function getPermission(token) {
  return request({
    url: '/admin/permission',
    method: 'get',
    params: { token }
  })
}

// outstanding
// export function logout() {
//   return request({
//     url: '/vue-element-admin/user/logout',
//     method: 'post'
//   })
// }
