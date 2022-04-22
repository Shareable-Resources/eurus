<template>
  <div>
    <b-navbar toggleable="lg" type="dark" variant="primary">
      <!-- <b-navbar-brand href="#">NavBar</b-navbar-brand> -->

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item
            @click="goToPage('/userList')"
            :active="checkActive('/userList')"
            >User</b-nav-item
          >
          <b-nav-item
            @click="goToPage('/depositList')"
            :active="checkActive('/depositList')"
            >Deposit</b-nav-item
          >
          <b-nav-item
            @click="goToPage('/withdrawList')"
            :active="checkActive('/withdrawList')"
            >Withdraw Requests</b-nav-item
          >
          <b-nav-item
            @click="goToPage('/withdrawHistory')"
            :active="checkActive('/withdrawHistory')"
            >Withdraw History</b-nav-item
          >
        </b-navbar-nav>

        <!-- Right aligned nav items -->
        <b-navbar-nav class="nav-left">
          <b-nav-item-dropdown right>
            <!-- Using 'button-content' slot -->
            <template #button-content> {{ username }} </template>
            <b-dropdown-item @click="logout">Sign Out</b-dropdown-item>
          </b-nav-item-dropdown>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
  </div>
</template>

<script>
import { getUsername, getWalletAddress, clearAll } from "../../store";

export default {
  name: "Navbar",
  data() {
    return {
      username: getUsername(),
      walletAddress: getWalletAddress(),
    };
  },
  methods: {
    checkActive(path) {
      let currentPath = "";
      try {
        currentPath = this.$router.history.current.path;
      } catch (error) {}
      return path === currentPath;
    },
    goToPage(page) {
      this.$router.push({
        path: page,
      });
    },
    logout(event) {
      event.preventDefault();
      clearAll();
      this.$router.push({
        path: "/",
      });
    },
  },
};
</script>

<style scoped>
.nav-left {
  margin-left: auto;
}
</style>
