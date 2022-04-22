import request from '@/utils/request'
// import loginToken
const Authorization = 'string'
// import concat query func
function joinQuery(query) {
  const arr = []
  for (var i in query) {
    arr.push(`${i}=${query[i]}`)
  }
  return arr.join(`&`)
}

export function getWithdrawList() {
  return request({
    url: '/admin/pendingWithdraw',
    method: 'get',
    headers: { Authorization }
  })
}

export function updateWithdraw({ pendingWithdrawId, ...data }) {
  return request({
    url: `/admin/pendingWithdraw/${pendingWithdrawId}`,
    method: 'put',
    headers: { Authorization },
    data
  })
}

export function getWithdrawHistory({ merchantId, ...query }) {
  let queryStr = ``
  if (query) {
    queryStr = `?${joinQuery(query)}`
  }
  return request({
    url: `/admin/pendingWithdrawHistory/${merchantId}${queryStr}`,
    method: 'get',
    headers: { Authorization }
  })
}
