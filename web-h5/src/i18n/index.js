
import Vue from 'vue'
import VueI18n from 'vue-i18n'
import EN from './en'
import ZH from './zh'
import { getLocale } from '@/utils/auth'

Vue.use(VueI18n)

const lang = getLocale() || 'en'
const messages = {
    en : EN,
    zh : ZH
}
const i18n = new VueI18n({
    locale : lang,
    messages : messages
})

export default i18n