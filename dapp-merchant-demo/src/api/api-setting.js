import request from '@/utils/request'
// import loginToken
const Authorization = 'string'

export function getAPISettingList() {
  return request({
    url: '/admin/pendingConfirmAPISetting',
    method: 'get',
    headers: { Authorization }
  })
}

// outstanding
export function updateAPISetting({ pendingId, ...data }) {
  return request({
    url: `/admin/pendingConfirmAPISetting/${pendingId}`,
    method: 'put',
    headers: { Authorization },
    data
  })
}
