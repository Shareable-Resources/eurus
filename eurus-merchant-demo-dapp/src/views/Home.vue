<template>
  <b-container class="home-container mt-5">
    <b-overlay :show="isLoading" rounded="sm">
      <h3 style="text-align: center">Eurus Merchant Demo</h3>
      <h5 style="text-align: center">Dapp</h5>

      <b-form>
        <b-form-group label="Username:">
          <b-form-input
            v-model="form.username"
            :readonly="true"
            required
          ></b-form-input>
        </b-form-group>

        <b-form-group label="Wallet Address:">
          <b-form-input
            v-model="form.walletAddress"
            :readonly="true"
            required
          ></b-form-input>
        </b-form-group>

        <b-form-group label="Balance:">
          <b-form-input
            v-model="balance"
            :readonly="true"
            required
          ></b-form-input>
        </b-form-group>

        <b-form-group class="mt-5" label="Deposit Amount(USDT):">
          <b-form-input
            v-model="form.depositAmount"
            type="number"
            placeholder=""
            required
          ></b-form-input>
        </b-form-group>

        <b-button
          class="mt-2 w-100"
          type="submit"
          variant="primary"
          @click="depositClick"
          >Deposit</b-button
        >

        <div
          class="txn-id-display bg-secondary mt-1 text-light"
          v-if="depositTxnId != ''"
        >
          <p>Status: {{ depositStatus }}</p>
          <p>Txn Id : {{ depositTxnId }}</p>
        </div>

        <b-form-group class="mt-5" label="Withdraw Amount(USDT):">
          <b-form-input
            v-model="form.withdrawAmount"
            type="number"
            placeholder=""
            required
          ></b-form-input>
        </b-form-group>

        <b-button
          class="mt-2 w-100"
          type="submit"
          variant="primary"
          @click="withdrawClick"
          >Request Withdraw</b-button
        >

        <div
          class="txn-id-display bg-secondary mt-1 text-light"
          v-if="refundRequestId != ''"
        >
          <p>Status: {{ refundStatus }}</p>
          <p>Request ID : {{ refundRequestId }}</p>
        </div>

        <b-button
          class="mt-5 w-100"
          type="submit"
          variant="primary"
          @click="logout"
          >Logout</b-button
        >
      </b-form>
    </b-overlay>
  </b-container>
</template>

<script>
import {
  depositToDapp,
  getTransaction,
  amountDivideDecimals,
  amountMultipleDecimals,
} from "../web3";
import { requestRefund, getUserBalance } from "../api";
import {
  getBalance,
  getUsername,
  getWalletAddress,
  setBalance,
  clearAll,
} from "../store";

export default {
  data() {
    return {
      form: {
        username: getUsername(),
        walletAddress: getWalletAddress(),
        depositAmount: 0,
        withdrawAmount: 0,
      },
      depositStatus: "",
      depositTxnId: "",
      refundStatus: "",
      refundRequestId: "",
      balance: getBalance(),
      isLoading: false,
    };
  },
  mounted: function () {
    this.refreshUserBalance();
  },
  methods: {
    logout(event) {
      event.preventDefault();
      clearAll();
      this.$router.push({
        path: "/",
      });
    },
    depositClick(event) {
      event.preventDefault();
      // alert(JSON.stringify(this.form.depositAmount));
      if (isNaN(this.form.depositAmount)) {
        alert("Not a number!");
      } else if (this.form.depositAmount <= 0) {
        alert("Less than 0!");
      } else {
        this.submitDeposit(this.form.depositAmount);
      }
    },
    withdrawClick(event) {
      event.preventDefault();
      // alert(JSON.stringify(this.form.withdrawAmount));
      if (isNaN(this.form.withdrawAmount)) {
        alert("Not a number!");
      } else if (this.form.withdrawAmount <= 0) {
        alert("Less than 0!");
      } else {
        this.submitRefund(this.form.withdrawAmount);
      }
    },
    async submitDeposit(depositAmount) {
      this.isLoading = true;
      this.depositStatus = "";
      this.depositTxnId = "";
      try {
        let txnHash = await depositToDapp(depositAmount);
        if (txnHash) {
          this.resultText = "Txn hash: " + txnHash;
          console.log("###### txnHash: ", txnHash);
          await this.getTransactionReceipt(txnHash);
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
    async submitRefund(amount) {
      this.isLoading = true;
      this.refundStatus = "";
      this.refundRequestId = "";
      try {
        let requestRefundResponse = await requestRefund(
          amountMultipleDecimals(amount, 6),
          "Refund"
        );
        if (
          requestRefundResponse &&
          requestRefundResponse.success &&
          requestRefundResponse.data &&
          requestRefundResponse.data.id
        ) {
          this.refundStatus = "Request submitted";
          this.refundRequestId = requestRefundResponse.data.id;
          alert("Request submitted!");
          this.refreshUserBalance();
        } else {
          this.refundStatus = "Failed!";
          alert("Network error!");
        }
      } catch (error) {
        alert("Network error!");
        console.log("Network error:", error);
      }
      this.isLoading = false;
    },
    async getTransactionReceipt(txnHash) {
      console.log("getTransactionReceipt start");
      let txnData = await getTransaction(txnHash);
      console.log("txnData:", txnData);
      if (txnData && txnData.txn_id) {
        if (txnData.status) {
          this.depositStatus = "Success";
          alert("Success!");
          this.refreshUserBalance();
        } else {
          this.depositStatus = "Failed";
          if (txnData.revertReason) {
            alert("Failed! RevertReason:" + txnData.revertReason);
            this.depositStatus = "Failed! RevertReason:" + txnData.revertReason;
          } else {
            alert("Failed!");
          }
        }
        this.depositTxnId = txnHash;
        this.isLoading = false;
      } else {
        setTimeout(() => {
          this.getTransactionReceipt(txnHash);
        }, 3000);
      }
    },
    async refreshUserBalance() {
      console.log("refreshUserBalance start");
      try {
        let getUserBalanceResponse = await getUserBalance();
        console.log("getUserBalanceResponse:", getUserBalanceResponse);
        if (
          getUserBalanceResponse &&
          getUserBalanceResponse.success &&
          getUserBalanceResponse.data &&
          getUserBalanceResponse.data.balance != ""
        ) {
          setBalance(getUserBalanceResponse.data.balance);
          this.balance = amountDivideDecimals(
            getUserBalanceResponse.data.balance,
            6
          );
        }
      } catch (error) {
        alert("Network error!");
        console.log("Network error:", error);
      }
    },
  },
};
</script>

<style scoped>
.home-container {
  max-width: 500px;
}
.txn-id-display {
  word-break: break-all;
}
</style>