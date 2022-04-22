// import { Cookies } from "quasar"
// import Cookies from 'js-cookies'
import Vue from 'vue'
const Username = 'name'
const Tokenkey = 'token'
const Locale = 'locale'
const Accounttype = 'accounttype'

export function setName(){
    return sessionStorage.setItem(Username, 'testing')
}

export function getName(){
    return sessionStorage.getItem(Username)
}
export function clearName() {
    return sessionStorage.removeItem(Username)
}

export function setToken(){
    return Vue.$cookies.set(Tokenkey, '12345678')
}
export function getToken() {
    return Vue.$cookies.get(Tokenkey)
}

export function clearToken() {
    return Vue.$cookies.remove(Tokenkey)
}
export function getLocale() {
    return sessionStorage.getItem(Locale)
}

export function setLocale(locale) {
    return sessionStorage.setItem(Locale, locale)
}


export function setAccounttype(type) {
    return sessionStorage.setItem(Accounttype, type)
}

export function getAccounttype() {
    return sessionStorage.getItem(Accounttype)
}
