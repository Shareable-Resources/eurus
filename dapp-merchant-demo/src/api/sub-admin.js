import request from '@/utils/request'
// import loginToken
const Authorization = 'string'
// import nonce method

export function getSubadminList() {
  return request({
    url: '/admin/subAdmin',
    method: 'get',
    headers: { Authorization }
  })
}

export function addSubadmin(data) {
  return request({
    url: '/admin/subAdmin',
    method: 'post',
    headers: { Authorization },
    data
  })
}

export function updateSubadminRole({ adminId, ...data }) {
  return request({
    url: `/admin/subAdmin/role/${adminId}`,
    method: 'put',
    headers: { Authorization },
    data: {
      ...data,
      adminId
    }
  })
}

export function updateSubadminPw({ adminId, ...data }) {
  return request({
    url: `/admin/subAdmin/password/${adminId}`,
    method: 'put',
    headers: { Authorization },
    data
  })
}

export function updateSubadminStatus({ adminId, ...data }) {
  return request({
    url: `/admin/subAdmin/status/${adminId}`,
    method: 'put',
    headers: { Authorization },
    data
  })
}

export function deleteSubadmin(adminId) {
  return request({
    url: `/admin/role/${adminId}`,
    method: 'delete',
    headers: { Authorization }
  })
}
