<template>
  <b-container class="login-container mt-5">
    <b-overlay :show="isLoading" rounded="sm">
      <h3 style="text-align: center"> Eurus Merchant Demo</h3>
      <h5 style="text-align: center"> Admin Portal</h5>

      <b-form @submit="onSubmit">
        <b-form-group id="input-group-1" label="Username:" label-for="input-1">
          <b-form-input
            id="input-1"
            v-model="form.username"
            placeholder=""
            required
          ></b-form-input>
        </b-form-group>

        <b-form-group id="input-group-1" label="Password:" label-for="input-1">
          <b-form-input
            id="input-1"
            v-model="form.password"
            type="password"
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
import { getAccount, switchEthereumChain, checkDappBrowser } from "../web3";
import { setUsername, setWalletAddress, setOperatorId } from "../store";
import { merchantLogin } from "../api";

export default {
  data() {
    return {
      form: {
        username: "MC1",
        password: "123456",
      },
      isLoading: false,
    };
  },
  methods: {
    onSubmit(event) {
      event.preventDefault();
      // alert(JSON.stringify(this.form));
      this.merchantLogin(this.form.username, this.form.password);
    },
    async merchantLogin(username, password) {
      this.isLoading = true;
      console.log("merchantLogin start");
      try {
        let merchantLoginResponse = await merchantLogin(username, password);
        console.log("merchantLoginResponse:", merchantLoginResponse);
        if (
          merchantLoginResponse &&
          merchantLoginResponse.success &&
          merchantLoginResponse.data &&
          merchantLoginResponse.data.username &&
          merchantLoginResponse.data.operatorId
        ) {
          let walletAddress = await this.loginWithMetamask();
          if (walletAddress) {
            setOperatorId(merchantLoginResponse.data.operatorId);
            setUsername(username);
            setWalletAddress(walletAddress);
            this.isLoading = false;
            this.$router.push({
              path: "/userList",
            });
          }
        }
      } catch (error) {
        alert("Network error!");
        console.log("Network error:", error);
      } finally {
        this.isLoading = false;
      }
    },
    async loginWithMetamask() {
      let accountAddress = "";
      try {
        let isDappbrowser = false;
        if (checkDappBrowser()) {
          isDappbrowser = true;
        }
        if (isDappbrowser) {
          const switchResult = await switchEthereumChain();
          if (switchResult) {
            const accounts = await getAccount();
            accountAddress = accounts[0];
          }
        }
      } catch (err) {
        console.error(err);
      }
      return accountAddress;
    },
  },
};
</script>

<style scoped>
.login-container {
  max-width: 500px;
}
</style>