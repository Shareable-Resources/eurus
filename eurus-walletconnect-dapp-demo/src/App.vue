<template>
  <img alt="" src="./assets/Eurus_Blue.png" style="height:10vh">
  <Basic msg="Wallet Connect Demo"/>
</template>

<script>
import Basic from './components/Basic.vue'

export default {
  name: 'App',
  components: {
    Basic
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
