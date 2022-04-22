<template>
  <b-container class="home-container">
    <nav-bar />
    <b-overlay :show="isLoading" rounded="sm">
      <div>
        <b-table responsive striped hover :fields="fields" :items="items">
          <template #cell(fromTokenAmt)="data">
            {{ getDisplayBalance(data.item.fromTokenAmt) }}
          </template>
          <template #cell(status)="data">
            {{ getDisplayStatus(data.item.status) }}
          </template>
        </b-table>
      </div>
    </b-overlay>
  </b-container>
</template>

<script>
import NavBar from "../components/navbar";
import { getRefundRequestLedgers } from "../api";
import { amountDivideDecimals } from "../web3";

export default {
  components: {
    NavBar,
  },
  data() {
    return {
      fields: [
        { key: "reqId", label: "Request Id" },
        { key: "toUserId", label: "UserId" },
        {
          key: "txHash",
          label: "Tx Hash",
          tdClass: "withdraw-history-table-td",
        },
        {
          key: "toWalletAddr",
          label: "Wallet Address",
          tdClass: "withdraw-history-table-td",
        },
        { key: "toAssetId", label: "Asset" },
        { key: "fromTokenAmt", label: "Amount" },
        { key: "status", label: "status" },
        { key: "remarks", label: "remarks" },
        { key: "createdDate", label: "createdDate" },
        { key: "lastModifiedDate", label: "lastModifiedDate" },
      ],
      items: [],
      isLoading: false,
    };
  },
  mounted: function () {
    this.getRefundRequestLedgers();
  },
  methods: {
    getDisplayStatus(inputStatus) {
      return inputStatus === 80 ? "Success" : "";
    },
    getDisplayBalance(inputBalance) {
      return amountDivideDecimals(inputBalance, 6);
    },
    async getRefundRequestLedgers() {
      this.isLoading = true;
      console.log("getRefundRequestLedgers start");
      try {
        let getRefundRequestLedgersResponse = await getRefundRequestLedgers();
        if (
          getRefundRequestLedgersResponse &&
          getRefundRequestLedgersResponse.success &&
          getRefundRequestLedgersResponse.data
        ) {
          const responseList = getRefundRequestLedgersResponse.data;
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

<style>
.withdraw-history-table-td {
  word-break: break-word;
}
</style>