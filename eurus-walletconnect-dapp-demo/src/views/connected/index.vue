<template>
  <section class="main">
    <button @click="sendEth">Send Eth</button>
    <button @click="personalSign">Personal Sign</button>
    <button @click="typedSign">Typed Sign</button>
  </section>
</template>

<script>

import {EthCrypto} from "eth-crypto";
import {walletConnect} from "../../api"
import {uuidv4} from "../../index"
import {eip712} from "../../eip712"

export default {
  name: "walletConnect",
  data() {
    return {

    };
  },
  methods: {
    // async onDisconnect() {
    //   const connector = await walletConnect();
    //   connector.killSession();
    //   console.log("killed sessions");
    // },
    async sendEth() {
      const connector = await walletConnect();
      const address = connector.accounts[0];
      const chainId = connector.chainId;
      //   let txndetails = getTxnDetails(this.network, this.$route.query.txnHash);
      //   console.log("##### created txndetails:", txndetails);

      console.log("###sendEth", connector, address, chainId);
      if (!connector) {
        return;
      }

      //const from = getAddress();
      const from = "0x866d119ba3c74b888732d8c55e685925d548b9e6";
      //   const from = address;
      const to = address;
      const nonce = parseInt(uuidv4());
      console.log("####nonce", nonce);
      const gasPrice = 1;
      const gasLimit = 1;
      const value = 0;
      const data = "0x";

      // test transaction
      const tx = {
        from,
        to,
        nonce,
        gasPrice,
        gasLimit,
        value,
        data,
      };
      console.log("###tx", tx);

      try {
        const result = await connector.sendTransaction(tx);
        console.log("###result", result);
      } catch (error) {
        console.error(error);
      }
    },


    async personalSign() {
      const connector = await walletConnect();
      const address = connector.accounts[0];
      const message = "Testing Personal Signature, email eu16@18m.dev";

      // Currently follows the HashMsg for our Api's signatures
      let hashMsg = EthCrypto.hash.keccak256(message).substring(2);
      console.log("### hashMsg:", hashMsg);

      const msgParams = [hashMsg, address];
      try {
        const result = await connector.signPersonalMessage(msgParams);
        console.log("###result", result);
      } catch (error) {
        console.error(error);
      }
    },


    async typedSign() {
      const connector = await walletConnect();
      const address = connector.accounts[0];

      //Message has to be EIP712 standard
      const message = JSON.stringify(eip712.example);
      
      console.log("Message Stringify:", message);

      // eth_signTypedData params
      const msgParams = [address, message];

      try {
        const result = await connector.signTypedData(msgParams);
        const hash = EthCrypto.hash.keccak256(message).substring(2);

        console.log("####result", result, "####hash", hash);
      } catch (error) {
        console.error(error);
      }
    },
  },
};
</script>

<style>

</style>
