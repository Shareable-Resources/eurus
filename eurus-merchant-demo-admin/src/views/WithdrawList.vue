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
          <template #cell(operation)="data">
            <b-button
              type="submit"
              variant="primary"
              @click="approveItem(data.item)"
              v-if="data.item.status === 10"
            >
              Approve
            </b-button>
          </template>
        </b-table>
      </div>
    </b-overlay>
  </b-container>
</template>

<script>
import NavBar from "../components/navbar";
import { getRefundRequestList } from "../api";
import { refund, getTransaction, amountDivideDecimals } from "../web3";

export default {
  components: {
    NavBar,
  },
  data() {
    return {
      fields: [
        { key: "id", label: "ID" },
        { key: "userId", label: "UserId" },
        { key: "username", label: "Username" },
        { key: "toAssetId", label: "Asset" },
        { key: "fromTokenAmt", label: "Amount" },
        { key: "fromWalletAddr", label: "Wallet Address" },
        { key: "status", label: "Status" },
        { key: "lastModifiedDate", label: "Last Modified Date" },
        { key: "operation", label: "Operation" },
      ],
      items: [],
      isLoading: false,
    };
  },
  mounted: function () {
    this.getRefundRequestList();
  },
  methods: {
    getDisplayBalance(inputBalance) {
      return amountDivideDecimals(inputBalance, 6);
    },
    getDisplayStatus(inputStatus) {
      return inputStatus === 10
        ? "Pending"
        : inputStatus === 20
        ? "Approved"
        : "";
    },
    approveItem(item) {
      // alert(item);
      this.submitRefund(item.id, item.fromWalletAddr, item.fromTokenAmt);
    },
    async getRefundRequestList() {
      this.isLoading = true;
      console.log("getRefundRequestList start");
      try {
        let getRefundRequestListResponse = await getRefundRequestList();
        if (
          getRefundRequestListResponse &&
          getRefundRequestListResponse.success &&
          getRefundRequestListResponse.data
        ) {
          const responseList = getRefundRequestListResponse.data;
          this.items = responseList;
        }
      } catch (error) {
        alert("Network error!");
        console.log("Network error:", error);
      } finally {
        this.isLoading = false;
      }
    },
    async submitRefund(requestID, targetAddress, amount) {
      this.isLoading = true;
      try {
        let txnHash = await refund(requestID, targetAddress, amount);
        if (txnHash) {
          this.resultText = "Txn hash: " + txnHash;
          console.log("###### txnHash: ", txnHash);
          await this.getTransactionReceipt(requestID, txnHash);
          console.log("###### after getTransactionReceipt ");
        }
      } catch (error) {
        this.isLoading = false;
        if (error.code === 4001) {
          // EIP-1193 userRejectedRequest error
          console.log("Please connect to MetaMask.");
          alert("Please connect to MetaMask.");
        } else if (error === 9991) {
          alert("Metamask wrong network!");
          console.log("Metamask wrong network!");
        } else {
          console.error(error);
        }
      }
    },
    async getTransactionReceipt(requestID, txnHash) {
      console.log("getTransactionReceipt start");
      let txnData = await getTransaction(txnHash);
      console.log("txnData:", txnData);
      if (txnData && txnData.txn_id) {
        if (txnData.status) {
          alert("Success!");
          this.getRefundRequestList();
        } else {
          if (txnData.revertReason) {
            alert("Failed! RevertReason:" + txnData.revertReason);
          } else {
            alert("Failed!");
          }
        }
        this.depositTxnId = txnHash;
        this.isLoading = false;
      } else {
        setTimeout(() => {
          this.getTransactionReceipt(requestID, txnHash);
        }, 3000);
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