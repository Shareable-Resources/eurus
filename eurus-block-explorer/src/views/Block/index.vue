<template>
  <div class="home block">
    <div class="trans flexbox-wrap">
      <div class="trans-details flexbox-left box-shadow bd-rds8">
        <div class="card">
          <div class="name card-content">
            <span>{{name}}</span>
          </div>
          <div class="hash card-content">
            <div class="text ellipsis-text">
              {{details.number}}
            </div>
            <div class="copy-icon" @click="copy(details.number)">
            </div>
          </div>
          <ul class="info">
            <li class="hash-icon">
              <span class="li-label">Hash:</span>
              <span class="li-val">{{details.hash | hashSplit(5)}}</span>
              <div class="arrow"></div>
            </li>
            <li class="hash-icon">
              <span class="li-label">Parent Hash:</span>
              <span class="li-val">{{details.parentHash | hashSplit(5)}}</span>
            </li>
          </ul>
          <ul class="info">
            <li class="mined-icon">
              <span class="li-label">Mined By:</span>
              <span class="li-val">{{details.miner | hashSplit(5)}}</span>
            </li>
            <li class="time-icon">
              <span class="li-label">Time:</span>
              <span class="li-val">{{details.timestampVerifiedISO | moment("from", "now")}}</span>
            </li>
            <li class="size-icon">
              <span class="li-label">Size:</span>
              <span class="li-val">{{details.size | amountFormatter}}{{` Bytes`}}</span>
            </li>
          </ul>
          <div class="additional">
            <span
              @click="showDetail"
              class="drop-box"
              :class="!detailShow ? 'icon-drop':'icon-up'"
            >
              <span>
                {{detailTitle}}
              </span>
            </span>
          </div>
        </div>
      </div>
      <div class="trans-value flexbox-right bd-rds8"
      >
        <div class="card">
          <div class="val-header">
            <span>Transaction Count</span>
          </div>
          <div class="val-content content">
            <span class="amount">{{details.transactionCount}}</span>
          </div>
        </div>
      </div>
      <div class="additional-tip"
        @click="showDetail"
      >
        <span
          class="drop-box"
          :class="!detailShow ? 'icon-drop':'icon-up'"
        >
          <span>
            {{detailTitle}}
          </span>
        </span>
        <span>{{!detailShow ?`+`:`-`}}</span>
      </div>
    </div>
    <div class="additional-box flexbox-wrap"
      v-show="detailShow"
    >
      <div class="left box-shadow trans-details flexbox-left bd-rds8">
        <div class="card">
          <ul class="nonce-wrap">
            <li class="nonce ">
              <span class="li-label">Nonce:</span>
              <span class="li-val">{{details.nonce}}</span>
            </li>
          </ul>
          <ul class="nonce-wrap">
            <li class="block">
              <span class="li-label">Difficulty:</span>
              <span class="li-val">{{details.difficulty}}</span>
              <span class="drop-icon"></span>
            </li>
            <li class="bytecode">
              <div>
                <span class="li-label">Total Difficulty:</span>
                <span class="li-val">{{details.totalDifficulty | amountFormatter}}</span>
              </div>
            </li>
          </ul>
          <ul class="nonce-wrap">
            <li class="gas">
              <div>
                <span class="li-label">Gas Used:</span>
                <span class="li-val">{{details.gasUsed | amountFormatter}}</span>
              </div>
            </li>
            <li class="transaction">
              <div>
                <span class="li-label">Gas Limit:</span>
                <span class="li-val">{{details.gasLimit | amountFormatter}}</span>
              </div>
            </li>
          </ul>
        </div>
      </div>
      <div class="flexbox-right">
      </div>
    </div>
    <div class="content tab-wrap box-shadow bd-rds8">
      <div class="tab-box">
        <div class="title blue-title ">
          {{`Transaction`}}
        </div>
        <div class="filter-wrap">
          <div class="filter-pagesize">
            <div class="filter">
              <div  @click="showFilter">{{`filter (${filter.length})`}}</div>
            </div>
            <div class="pagesize">
              <el-select
                size="mini"
                v-model="num"
                placeholder=""
                @change="onPagesizeChange"
              >
                <el-option
                  v-for="item in 5"
                  :key="item"
                  :label="item"
                  :value="item">
                </el-option>
              </el-select>
              <span class="text">{{`Row per page`}}</span>
            </div>
          </div>
          <div class="sorting" v-show="filterShow">
            <div class="sorting-mobile">
              <div class="label">{{`Filter:`}}</div>
              <el-select
                multiple
                :multiple-limit="2"
                v-model="filter"
                placeholder=""
                size="mini"
                @change="onFilterChange"
                @remove-tag="onFilterRemove"
              >
                <el-option
                  v-for="(item,key) in FILTER_BY"
                  :key="key"
                  :label="item"
                  :value="key">
                </el-option>
              </el-select>
            </div>
          </div>
          <div class="sorting">
            <div class="sorting-mobile">
              <div class="label">{{`Sorting:`}}</div>
              <el-select
                v-model="sort"
                placeholder=""
                size="mini"
                @change="onSortingChange"
              >
                <el-option
                  v-for="(item,key) in SORT_OPTIONS"
                  :key="key"
                  :label="item"
                  :value="key">
                </el-option>
              </el-select>
            </div>
          </div>
        </div>
        <div class="transactions">
          <Empty v-if="!events.data" :msg="EMPTY_MSG[activeTab]"/>
          <div v-else>
            <div
              class="event"
              v-for="(item,index) in events.data"
              :key="`event_${index}`"
            >
              <div class="event-index">{{`#${index+1}`}}</div>
              <div class="event-content">
                <div class="name">
                  <span class="trans-name">{{item.hash | hashSplit(5)}}</span>
                  <span class="type">{{FILTER_BY[item.transactionType]}}</span>
                </div>
                <div class="trans-data bd-rds8">
                  <div class="trans">
                    <template v-if="!events.data">
                      <span>{{`No details`}}</span>
                    </template>
                    <template v-else>
                      <div class="from-to">
                        <div class="from">{{`from: ${item.from}`}}</div>
                        <div class="to">{{`to: ${item.to}`}}</div>
                      </div>
                      <div class="from-to-mobile">
                        <div class="from">{{item.from | hashSplit(5)}}</div>
                        <div class="arrow"></div>
                        <div class="to">{{item.to  | hashSplit(5)}}</div>
                      </div>
                      <div class="value">
                        <div class="">{{`Value: ${item.value}`}}</div>
                        <div class="">{{`Function: ${item.direction}`}}</div>
                      </div>
                    </template>
                  </div>
                </div>
              </div>
            </div>
            <div class="pagging">
              <div class="total">{{`Showing ${events.paging.page+1} out of ${events.paging.totalElements} ${events.paging.totalElements > 1 ? 'transactions': 'transaction'}`}}</div>
              <div class="pagination"   v-if="events.paging.totalPages > 1" >
                <div class="back">
                </div>
                <div class="pages">
                  <template
                    v-if="events.paging.totalElements < 5"
                  >
                    <span
                      v-for="item in 4"
                      :key="item"
                      class="page"
                      :class="events.paging.page+1 === item? 'activePage':''"
                    >
                      {{item}}
                    </span>
                  </template>
                  <template
                    v-else
                  >
                    <span
                      v-for="item in 3"
                      :key="item"
                      class="page"
                      :class="events.paging.page+1 === item? 'activePage':''"
                    >
                      {{item}}
                    </span>
                    <span
                      :key="events.paging.totalPages"
                      class="page"
                    >
                      {{events.paging.totalPages}}
                    </span>
                  </template>
                </div>
                <div class="push">
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {blockAPI, blockTransAPI} from '@/api'
import Empty from '@/components/Empty.vue'
export default {
  name: 'Block',
  components: {
    Empty
  },
  filters: {
    hashSplit(value, unit) {
      if (!value) return ''
      if (value.length < 12) return value
      let arr = []
      arr.push(value.slice(0,unit))
      arr.push(value.slice(-unit))
      return arr.join('...')
    },
    amountFormatter(value) {
      if(!value) return value
      return new Intl.NumberFormat('en-US', { style: 'decimal' }).format(value)
    },
    eventItemFilter(event) {
      console.log(event)
      if(!event) return ''
      return `${event.name}: ${event.value}`
    }
  },
  data() {
    return {
      name: 'Block Details',
      detailTitle: 'Additional Details',
      details: {},
      events: {},
      detailShow: false,
      filterShow: false,
      activeTab: 'transactions',
      EMPTY_MSG: {
        transactions: 'There are no transactions in this block'
      },
      SORT_OPTIONS: {
        high: 'Highest first',
        low: 'Lowest first'
      },
      FILTER_BY: {
        ContractCall: 'Contract Call',
        Transfer: 'Transfer',
        ContractCreation: 'Contract Creation'
      },
      sort: '',
      filter: [],
      num: 1
    }
  },
  methods: {
    showFilter() {
      this.filterShow = !this.filterShow;
    },
    showDetail() {
      this.detailShow = !this.detailShow;
    },
    copy(message) {
      var input = document.createElement("input");
      input.value = message;
      document.body.appendChild(input);
      input.select();
      input.setSelectionRange(0, input.value.length), document.execCommand('Copy');
      document.body.removeChild(input);
      this.$alert(`Copy successed!`,`Tips`,{
        closeOnClickModal: true,
        showConfirmButton: false,
        showClose:false,
        confirmButtonText: 'OK',
        callback: action => {
          console.log('action',action)
        }
      })
    },
    handleChange(value) {
      console.log(value);
    },
    onPagesizeChange(val) {
      console.log(val)
    },
    onFilterChange(val) {
      let {params} = this.$route;
      if(val.length === 0) return this.queryTrans(params.id, ()=>{})
      let query = `?filter=(`
      let tempStart = `%20transactionType%20eq%20`
      let tempEnd = `%20`
      val.forEach((item,index)=> {
        index > 0 ? query+=`or${tempStart}${item}${tempEnd}`:
        query+=`${tempStart}${item}${tempEnd}`
      })
      let end = `)&page=0`
      return this.queryTrans(params.id, ()=>{}, `${query}${end}`)
    },
    onFilterRemove() {
      let {filter} = this;
      return this.onFilterChange(filter)
    },
    onSortingChange(val) {
      let {params} = this.$route;
      const MAP = {
        high: `?sort=ethValue&direction=DESC`,
        low: `?sort=ethValue&direction=ASC`
      }
      return this.queryTrans(params.id, ()=>{}, MAP[val])
    },
    queryData(val,cb=()=>{}) {
      let API = blockAPI(val)
      this.axios.get(API)
      .then(({data}) => {
        this.$set(this,`details`,data)
        return cb()
      })
    },
    queryTrans(val, cb=()=>{},query=``) {
      let API = blockTransAPI(val,query)
      this.axios.get(API)
      .then(({data}) => {
        this.$set(this,`events`,data)
        return cb()
      })
    },
    initPage(id, cb=()=>{}) {
      this.queryData(id,cb)
      this.queryTrans(id,cb)
    }
  },
  beforeRouteUpdate(to, from, next) {
    if(to.params && to.params.id) {
      return this.initPage(to.params.id, next)
    }
    next();
  },
  mounted() {
    let {params} = this.$route;
    if(params && params.id) {
      this.initPage(params.id)
    }
  }
}
</script>
