import request from '@/utils/request'
// import loginToken
const Authorization = 'string'

export function getMerchantList() {
  return request({
    url: '/admin/merchant',
    method: 'get',
    headers: { Authorization }
  })
}

export function addMerchant(data) {
  return request({
    url: '/admin/merchant',
    method: 'post',
    headers: { Authorization },
    data
  })
}

export function getMerchant(merchantId) {
  return request({
    url: `/admin/merchant/${merchantId}`,
    method: 'get',
    headers: { Authorization }
  })
}
