<template>
  <b-container class="home-container">
    <nav-bar />
    <b-overlay :show="isLoading" rounded="sm">
      <div>
        <b-table responsive striped hover :fields="fields" :items="items">
          <template #cell(fromAssetAmt)="data">
            {{ getDisplayBalance(data.item.fromAssetAmt) }}
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
import { getDepositLedgers } from "../api";
import { amountDivideDecimals } from "../web3";

export default {
  components: {
    NavBar,
  },
  data() {
    return {
      fields: [
        { key: "txHash", label: "Tx Hash", tdClass: "deposit-table-td" },
        {
          key: "fromWalletAddr",
          label: "Wallet Address",
          tdClass: "deposit-table-td",
        },
        { key: "fromAssetId", label: "Asset" },
        { key: "fromAssetAmt", label: "Amount" },
        { key: "userId", label: "UserId" },
        { key: "status", label: "Status" },
        { key: "remarks", label: "Remarks" },
        { key: "createdDate", label: "created Date" },
        { key: "lastModifiedDate", label: "Last Modified Date" },
      ],
      items: [],
      isLoading: false,
    };
  },
  mounted: function () {
    this.getDepositLedgers();
  },
  methods: {
    getDisplayStatus(inputStatus) {
      return inputStatus === 10
        ? "Pending"
        : inputStatus === 40
        ? "Success"
        : "";
    },
    getDisplayBalance(inputBalance) {
      return amountDivideDecimals(inputBalance, 6);
    },
    async getDepositLedgers() {
      this.isLoading = true;
      console.log("getDepositLedgers start");
      try {
        let getDepositLedgersResponse = await getDepositLedgers();
        if (
          getDepositLedgersResponse &&
          getDepositLedgersResponse.success &&
          getDepositLedgersResponse.data
        ) {
          const responseList = getDepositLedgersResponse.data;
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
.deposit-table-td {
  word-break: break-word;
}
</style>