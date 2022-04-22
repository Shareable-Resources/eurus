// import Cookies from 'js-cookie'

const TokenKey = 'EurusMerchantDappDemoToken'
const AddressKey = 'EurusMerchantDappDemoAddress'

const EurusRpcDomainKey = 'EurusMerchantDappDemoRpcDomain'
const EurusRpcPortKey = 'EurusMerchantDappDemoRpcPort'
const EurusChainIdKey = 'EurusMerchantDappDemoChainId'

export function getToken() {
  return sessionStorage.getItem(TokenKey)
}

export function setToken(token) {
  return sessionStorage.setItem(TokenKey, token)
}

export function removeToken() {
  return sessionStorage.removeItem(TokenKey)
}

export function getAddress() {
  return sessionStorage.getItem(AddressKey)
}

export function setAddress(address) {
  return sessionStorage.setItem(AddressKey, address)
}

export function removeAddress() {
  return sessionStorage.removeItem(AddressKey)
}

export function getEurusRpcDomain() {
  return sessionStorage.getItem(EurusRpcDomainKey)
}

export function setEurusRpcDomain(eurusRpcDomain) {
  return sessionStorage.setItem(EurusRpcDomainKey, eurusRpcDomain)
}

export function removeEurusRpcDomain() {
  return sessionStorage.removeItem(EurusRpcDomainKey)
}

export function getEurusRpcPort() {
  return sessionStorage.getItem(EurusRpcPortKey)
}

export function setEurusRpcPort(eurusRpcPort) {
  return sessionStorage.setItem(EurusRpcPortKey, eurusRpcPort)
}

export function removeEurusRpcPort() {
  return sessionStorage.removeItem(EurusRpcPortKey)
}

export function getEurusChainId() {
  return sessionStorage.getItem(EurusChainIdKey)
}

export function setEurusChainId(eurusChainId) {
  return sessionStorage.setItem(EurusChainIdKey, eurusChainId)
}

export function removeEurusChainId() {
  return sessionStorage.removeItem(EurusChainIdKey)
}

export function getUsdtAddress() {
  return '0xa54Dee79c3bB34251DEbf86C1BA7D21898FFb7AC'
}

export function getDappAddress() {
  return '0xcdeE76833e5991289f87a1123fBAa3F0F56E2409'
}

export function getDappStockAddress() {
  return '0xBa571455DCEE4fed5C43669af9DeBfDC63a52508'
}
