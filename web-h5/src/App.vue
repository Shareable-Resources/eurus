<template>
<div id="app">
      <router-view/>
</div>
</template>

<script>
import { getToken , getAccounttype } from '@/utils/auth.js'


export default {
   name : 'App',
   data : function(){
     return {

     }
   },
   beforeCreate() {
     if(!getToken() || getToken() == undefined) {
       this.$router.push("/login")
     }

     if(!getAccounttype() || getAccounttype() == undefined || getAccounttype() === "centralized") {
       document.body.classList.add("centralized")
     } else {
       document.body.classList.add("decentralized")
     }


   },
}
</script>
