<template>
  <b-container class="home-container">
    <nav-bar />
    <b-overlay :show="isLoading" rounded="sm">
      <div>
        <b-table responsive striped hover :fields="fields" :items="items">
          <template #cell(balance)="data">
            {{ getDisplayBalance(data.item.balance) }}
          </template>
        </b-table>
      </div>
    </b-overlay>
  </b-container>
</template>

<script>
import NavBar from "../components/navbar";
import { getMerchantClients } from "../api";
import { amountDivideDecimals } from "../web3";

export default {
  components: {
    NavBar,
  },
  data() {
    return {
      fields: [
        { key: "id", label: "ID" },
        { key: "username", label: "Username" },
        { key: "balance", label: "Balance" },
        { key: "walletAddress", label: "Wallet Address" },
        { key: "createdDate", label: "Created Date" },
        { key: "lastModifiedDate", label: "Last Modified Date" },
        { key: "status", label: "Status" },
      ],
      items: [],
      isLoading: false,
    };
  },
  mounted: function () {
    this.getMerchantClients();
  },
  methods: {
    getDisplayBalance(inputBalance) {
      return amountDivideDecimals(inputBalance, 6);
    },
    async getMerchantClients() {
      this.isLoading = true;
      console.log("getMerchantClients start");
      try {
        let getMerchantClientsResponse = await getMerchantClients();
        if (
          getMerchantClientsResponse &&
          getMerchantClientsResponse.success &&
          getMerchantClientsResponse.data
        ) {
          const responseList = getMerchantClientsResponse.data;
          this.items = responseList;
        }
      } catch (error) {
        alert("Network error!");
        console.log("Network error:", error);
      } finally {
        this.isLoading = false;
      }
    },
  },
};
</script>

<style scoped>
.home-container {
  max-width: 1400px;
}
</style>