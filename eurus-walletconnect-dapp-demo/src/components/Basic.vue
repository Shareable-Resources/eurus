<template>
  <div class="">
    <h1>{{ msg }}</h1>
    <!-- <button>vue-cli documentation</button> -->
    <button @click="connectWc">Connect to Wallet Connect</button>
  </div>
</template>

<script>

import { walletConnect } from "../../src/api";

export default {
  name: "HelloWorld",
  props: {
    msg: String,
  },
  data: function () {
    // connected: false
  },
  methods: {
    async connectWc() {
      console.log("connectWc Clicked")
      walletConnect();

      const connector = await walletConnect();

      if (!connector.connected) {
          connector.createSession();
          this.$router.push('connected')

      }
    },
  },
}
</script>

<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
