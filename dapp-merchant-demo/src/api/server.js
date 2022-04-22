import request from '@/utils/request'
// import loginToken
const Authorization = 'string'

export function getServerList() {
  return request({
    url: '/admin/server/list',
    method: 'get',
    headers: { Authorization }
  })
}
