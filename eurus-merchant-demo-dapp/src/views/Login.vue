<template>
  <b-container class="login-container mt-5">
    <b-overlay :show="loading" rounded="sm">
      <h3 style="text-align: center">Eurus Merchant Demo</h3>
      <h5 style="text-align: center">Dapp</h5>

      <b-form @submit="onSubmit">
        <b-form-group id="input-group-1" label="Username:" label-for="input-1">
          <b-form-input
            id="input-1"
            v-model="form.username"
            placeholder=""
            required
          ></b-form-input>
        </b-form-group>

        <b-button class="mt-3 w-100" type="submit" variant="primary"
          >Login with Metamask</b-button
        >
      </b-form>
    </b-overlay>
  </b-container>
</template>

<script>
import {
  getAccount,
  switchEthereumChain,
  checkDappBrowser,
  // importWallet,
  importWalletForEurus,
} from "../web3";

import { setUsername, setWalletAddress, setEurusToken } from "../store";

export default {
  data() {
    return {
      form: {
        username: "",
      },
      loading: false,
    };
  },
  methods: {
    onSubmit(event) {
      event.preventDefault();
      //   alert(JSON.stringify(this.form.username));
      this.loginWithMetamask(this.form.username);
    },
    async loginWithMetamask(username) {
      this.loading = true;
      try {
        let isDappbrowser = false;
        if (checkDappBrowser()) {
          isDappbrowser = true;
        }
        if (isDappbrowser) {
          const switchResult = await switchEthereumChain();
          if (switchResult) {
            const accounts = await getAccount();
            const accountAddress = accounts[0];

            console.log("accountAddress: ", accountAddress);
            const importWalletResponse = await importWalletForEurus(
              accountAddress,
              username
            );
            if (
              importWalletResponse &&
              importWalletResponse.data &&
              importWalletResponse.data.data &&
              importWalletResponse.data.data.token
            ) {
              setEurusToken(importWalletResponse.data.data.token);
              setWalletAddress(accountAddress);
              setUsername(username);

              this.loading = false;
              // window.location.href = '/home'
              this.$router.push({
                path: "/home",
              });
            } else {
              if (importWalletResponse && importWalletResponse.errorMsg) {
                alert("Error: " + importWalletResponse.errorMsg);
              } else {
                alert("Network error!!");
              }
            }
          }
        }
      } catch (err) {
        console.error(err);
        this.loading = false;
        alert("Metamask Failed!");
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.login-container {
  max-width: 500px;
}
</style>