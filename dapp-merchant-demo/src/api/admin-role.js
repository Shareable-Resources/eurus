import request from '@/utils/request'
// import loginToken
const Authorization = 'string'

export function getRoleList() {
  return request({
    url: '/admin/role',
    method: 'get',
    headers: { Authorization }
  })
}

export function getRole(roleId) {
  return request({
    url: `/admin/role/${roleId}`,
    method: 'get',
    headers: { Authorization }
  })
}

// outstanding - no API
export function addRole(data) {
  return request({
    url: '/admin/role',
    method: 'post',
    data
  })
}

export function updateRole({ roleId, ...data }) {
  return request({
    url: `/admin/role/${roleId}`,
    method: 'put',
    headers: { Authorization },
    data: {
      ...data,
      roleId
    }
  })
}

export function deleteRole(roleId) {
  return request({
    url: `/admin/role/${roleId}`,
    method: 'delete',
    headers: { Authorization }
  })
}
