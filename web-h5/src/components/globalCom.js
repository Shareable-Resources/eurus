import inputField from '@/components/input/outlined'
import titleHeader from '@/components/titleHeader'

const plugins = {
    install(Vue){
        Vue.component('inputField', inputField)
        Vue.component('titleHeader', titleHeader)
    }
}


export default plugins;