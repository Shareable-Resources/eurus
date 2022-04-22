const EURUS_TOKEN_KEY = 'eurus_merchant_demo_token'
const USER_ID_KEY = 'eurus_merchant_demo_user_id'
const USERNAME_KEY = 'eurus_merchant_demo_username'
const WALLET_ADDRESS_KEY = 'eurus_merchant_demo_wallet_address'
const BALANCE_KEY = 'eurus_merchant_demo_balance'

export function setEurusToken(eurusToken) {
    return sessionStorage.setItem(EURUS_TOKEN_KEY, eurusToken)
}

export function getEurusToken() {
    return sessionStorage.getItem(EURUS_TOKEN_KEY)
}

export function clearEurusToken() {
    return sessionStorage.removeItem(EURUS_TOKEN_KEY)
}

export function setUserId(userId) {
    return sessionStorage.setItem(USER_ID_KEY, userId)
}

export function getUserId() {
    return sessionStorage.getItem(USER_ID_KEY)
}

export function clearUserId() {
    return sessionStorage.removeItem(USER_ID_KEY)
}

export function setUsername(username) {
    return sessionStorage.setItem(USERNAME_KEY, username)
}

export function getUsername() {
    return sessionStorage.getItem(USERNAME_KEY)
}

export function clearUsername() {
    return sessionStorage.removeItem(USERNAME_KEY)
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

export function setBalance(balance) {
    return sessionStorage.setItem(BALANCE_KEY, balance)
}

export function getBalance() {
    return sessionStorage.getItem(BALANCE_KEY)
}

export function clearBalance() {
    return sessionStorage.removeItem(BALANCE_KEY)
}

export function clearAll() {
    clearEurusToken()
    clearWalletAddress()
    clearUserId()
    clearUsername()
    clearBalance()
}