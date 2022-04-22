<template>
  <div id="top">
    <el-header class="header">
      <div class="header-left">
        <router-link to="/" class="logo">
          <img alt="logo" src="../assets/eurusLogo.png">
        </router-link>
      </div>
      <div class="header-right">
        <el-autocomplete
          v-model="input"
          :fetch-suggestions="querySearch"
          :placeholder="inputPlaceHolder"
          :trigger-on-focus="false"
          @select="handleSelect"
        >
          <i slot="prefix" class="el-input__icon el-icon-search"></i>
          <template slot-scope="{ item }">
            <span class="name">{{ `[${item.rel}]` }}</span>
            <span class="addr">{{ item.display }}</span>
          </template>
        </el-autocomplete>
      </div>
    </el-header>
  </div>
</template>

<script>
import {searchQueryAPI} from '@/api'
export default {
  name: 'Top',
  data() {
    return {
      input: '',
      timeout: null,
      inputPlaceHolder: 'Search by address, token, transaction hash, or block number'
    }
  },
  methods: {
    querySearch(queryString, cb) {
      if (this.timeout){
       clearTimeout(this.timeout); 
      }
      this.timeout = setTimeout(() => {
        this.queryData(queryString,cb)
      }, 300); // delay
    },
    handleSelect(item) {
      try {
        if(this.$route.path === item.href) return;
        this.$router.push({ path: `/${item.rel}s/${item.display}`, params: item })
        this.$set(this,'input','')
      } catch(err) {
        console.log(err)
      }
    },
    queryData(val,cb){
      const API = searchQueryAPI(val)
      this.axios.get(API).then(({data}) => {
        this.filterData(data.data)
        .then(arr => {
          cb(arr)
        })
      })
    },
    filterData(arr) {
      return new Promise((resolve) => {
        let MAP = {
          transactions: true,
          blocks: true
        }
        let newArr = []
        arr.forEach(item => {
          MAP[item.type] ? newArr.push(item.link) : null
        });
        resolve(newArr)
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
