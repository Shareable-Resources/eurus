
const USERNAME_KEY = 'eurus_merchant_admin_demo_username'
const OPERATOR_ID_KEY = 'eurus_merchant_admin_demo_operator_id'
const WALLET_ADDRESS_KEY = 'eurus_merchant_admin_demo_wallet_address'

export function setUsername(username) {
    return sessionStorage.setItem(USERNAME_KEY, username)
}

export function getUsername() {
    return sessionStorage.getItem(USERNAME_KEY)
}
export function clearUsername() {
    return sessionStorage.removeItem(USERNAME_KEY)
}

export function setOperatorId(operatorId) {
    return sessionStorage.setItem(OPERATOR_ID_KEY, operatorId)
}

export function getOperatorId() {
    return sessionStorage.getItem(OPERATOR_ID_KEY)
}
export function clearOperatorId() {
    return sessionStorage.removeItem(OPERATOR_ID_KEY)
}

export function setWalletAddress(walletAddress) {
    return sessionStorage.setItem(WALLET_ADDRESS_KEY, walletAddress)
}

export function getWalletAddress() {
    return sessionStorage.getItem(WALLET_ADDRESS_KEY)
}

export function clearWalletAddress() {
    return sessionStorage.removeItem(WALLET_ADDRESS_KEY)
}

export function clearAll() {
    clearWalletAddress()
    clearUsername()
    clearOperatorId()
}