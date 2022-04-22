<template>
  <div class="home">
    <div class="trans flexbox-wrap">
      <div class="trans-details flexbox-left box-shadow bd-rds8">
        <div class="card">
          <div class="name card-content">
            <span>{{name}}</span>
          </div>
          <div class="hash card-content">
            <div class="text ellipsis-text">
              {{details.hash}}
            </div>
            <div class="type bd-rds3">
              {{details.transactionType}}
            </div>
          </div>
          <ul class="addr">
            <li>
              <span class="li-label">From:</span>
              <span class="li-val">{{details.from | hashSplit(5)}}</span>
              <div class="arrow"></div>
            </li>
            <li>
              <span class="li-label">To:</span>
              <span class="li-val">{{details.to | hashSplit(5)}}</span>
            </li>
          </ul>
          <ul class="info">
            <li>
              <span class="li-label">Time:</span>
              <span class="li-val">{{details.verifiedTimestampISO | moment("from", "now")}}</span>
            </li>
            <li>
              <span class="li-label">Status:</span>
              <span class="li-val">{{details.status || ''}}</span>
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
      <div class="trans-value flexbox-right bd-rds8 bg-eun"
      >
        <div class="card">
          <div class="val-header">
            <span>Value</span>
          </div>
          <div class="val-content content">
            <span class="amount">{{details.value}}</span>
            <span class="unit">{{'EUN'}}</span>
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
          <ul class="block-wrap bottom-line">
            <li class="block active" @click="blockJump(details.blockHash)">
              <span class="li-label">Block:</span>
              <span class="li-val">{{ `#${details.blockNumber}`}}</span>
              <span class="drop-icon"></span>
            </li>
            <li class="bytecode" @click="showBytecode">
              <div>
                <span class="li-label">Input Bytecode:</span>
              </div>
              <div class="drop-icon"></div>
            </li>
          </ul>
          <div class="bytecode-content" v-show="bytecodeShow">
            <div class="bytecode-content-wrap">
              <div class="bytecode">{{details.input}}</div>
            </div>
          </div>
          <ul class="nonce-wrap">
            <li class="nonce">
              <span class="li-label">Nonce:</span>
              <span class="li-val">{{details.nonce || ''}}</span>
            </li>
            <li class="transaction">
              <span class="li-label">Transaction Fee:</span>
              <span class="li-val">{{details.cumulativeGasUsed | amountFormatter}}</span>
            </li>
          </ul>
          <ul class="position-wrap">
            <li class="position">
              <span class="li-label">Position:</span>
              <span class="li-val">{{details.transactionIndex}}</span>
            </li>
          </ul>
          <ul class="gas-wrap">
            <li class="gas">
              <span class="li-label">Gas Used:</span>
              <span class="li-val">{{details.gasUsed | amountFormatter}}{{` (${Math.round(details.gasUsed/details.gas*10000)/100}%)`}}</span>
            </li>
          </ul>
          <ul class="price-wrap">
            <li>
              <span class="li-label">Gas Price:</span>
              <span class="li-val">{{details.gasPrice | amountFormatter}}</span>
            </li>
          </ul>
          <ul class="val-wrap">
            <li>
              <span class="li-label">Value:</span>
              <span class="li-val">{{details.ethValue|amountFormatter}}{{` ETH / `}}{{details.value|amountFormatter}}</span>
            </li>
          </ul>
        </div>
      </div>
      <div class="flexbox-right">
      </div>
    </div>
    <div class="content tab-wrap box-shadow bd-rds8">
      <el-tabs
        class="tab-box"
        v-model="activeTab"
        @tab-click="switchTab"
      >
        <el-tab-pane
          label="Function"
          name="function"
        >
          <Empty v-if="!details.functionMeta" :msg="EMPTY_MSG[activeTab]"/>
          <div v-else>
            <div
              class="function event"
            >
              <div class="event-content">
                <div class="name">
                  {{`Name: `}}
                  <span
                    class="black-bold"
                    v-if="details.functionMeta && details.functionMeta.functionName"
                  >{{details.functionMeta.functionName}}</span>
                </div>
                <div class="trans-data bd-rds8">
                  <div class="title black-bold">{{`PARAMETERS`}}</div>
                  <div class="trans">
                    <div class="from-to-mobile">
                      <div
                        v-if="details.functionMeta && details.functionMeta.params"
                      >{{details.functionMeta.params.find(item=> item.name ==='_to') | eventItemFilter}}</div>
                    </div>
                    <div class="value">
                      <div
                        v-if="details.functionMeta && details.functionMeta.params"
                      >{{details.functionMeta.params.find(item=> item.name ==='_value') | eventItemFilter}}</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane
          label="Events"
          name="events"
        >
          <Empty v-if="!events.data || events.paging.totalElements < 1" :msg="EMPTY_MSG[activeTab]"/>
          <div v-else>
            <div
              class="event"
              v-for="(item,index) in events.data"
              :key="`event_${index}`"
            >
              <div class="event-index">{{`#${index+1}`}}</div>
              <div class="event-content">
                <div class="name">
                  {{`Name: `}}
                  <span class="black-bold">{{item.eventName|hashSplit(5)}}</span>
                </div>
                <div class="trans-data bd-rds8">
                  <div class="title black-bold">{{`PARAMETERS`}}</div>
                  <div class="trans">
                    <template v-if="!(item.parameters && item.parameters.length >1)">
                      <span>{{`No details`}}</span>
                    </template>
                    <template v-else-if="!(item.parameters&& item.parameters[0]&& item.parameters[0].name)">
                      <div 
                        class="params"
                        v-for="(tran,tranIndex) in item.parameters"
                        :key="`tran_${tranIndex}`"
                      >
                        <div class="param">{{`[${tran.value}]`}}</div>
                        <div class="param-mobile">{{`[`}}{{tran.value | hashSplit(5)}}{{`]`}}</div>
                      </div>
                    </template>
                    <template v-else>
                      <div class="from-to">
                        <div class="from">{{item.parameters.find(item=> item.name ==='from') | eventItemFilter}}</div>
                        <div class="to">{{item.parameters.find(item=> item.name ==='to') | eventItemFilter}}</div>
                      </div>
                      <div class="from-to-mobile">
                        <div class="from">{{item.parameters.find(item=> item.name ==='from').value | hashSplit(5)}}</div>
                        <div class="arrow"></div>
                        <div class="to">{{item.parameters.find(item=> item.name ==='to').value | hashSplit(5)}}</div>
                      </div>
                      <div class="value">
                        <div>{{item.parameters.find(item=> item.name ==='value') | eventItemFilter}}</div>
                      </div>
                    </template>
                  </div>
                </div>
              </div>
            </div>
            <div class="pagging" v-show="events.paging.totalElements > 0">
              <div class="total">{{`Showing ${events.paging.page+1} out of ${events.paging.totalElements} ${events.paging.totalElements > 1 ? 'events': 'event'}`}}</div>
              <div class="pagination"  v-if="events.paging.totalPages > 1">
                <div class="back">
                </div>
                <div class="pages">
                  <template
                    v-if="events.paging.totalPages < 5"
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
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import {transactionAPI, transactionEventAPI} from '@/api'
import Empty from '@/components/Empty.vue'
export default {
  name: 'Home',
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
      if(!event) return ''
      if(event.name ==='_to') {
        let arr = []
        arr.push(event.value.slice(0,5))
        arr.push(event.value.slice(-5))
        return `${event.name}: ${arr.join('...')}`
      }
      return `${event.name}: ${event.value}`
    }
  },
  data() {
    return {
      name: 'Transaction Details',
      detailTitle: 'Additional Details',
      details:{},
      events:{},
      detailShow: false,
      bytecodeShow: false,
      activeTab: 'function',
      UNIT_MAP: {
        EUN: 'eun',
        ETH: 'eth'
      },
      EMPTY_MSG: {
        function: 'There is no function metadata for this transaction',
        events: 'There is no event for this transaction'
      }
    }
  },
  methods: {
    showDetail() {
      this.detailShow = !this.detailShow;
    },
    switchTab(tab) {
      const QUERY_NAME = 'events'
      let {path} = this.$route
      if(tab.name === QUERY_NAME) {
        this.queryEvent(path)
      }
    },
    blockJump(blockHash) {
      this.$router.push({ path: `/blocks/${blockHash}` , params: { blockHash } })
    },
    showBytecode() {
      this.bytecodeShow = !this.bytecodeShow;
    },
    queryEvent(val, cb=()=>{}) {
      let API = transactionEventAPI(val)
      this.axios.get(API)
      .then(({data}) => {
        this.$set(this,`events`,data)
        return cb()
      })
    },
    queryData(val,cb=()=>{}) {
      let API = transactionAPI(val)
      this.axios.get(API)
      .then(({data}) => {
        this.$set(this,`details`,data)
        return cb()
      })
    },
    initPage(id, cb=()=>{}) {
      this.$set(this,'activeTab','function')
      this.$set(this,'events',{})
      this.queryData(id,cb)
    }
  },
  beforeRouteUpdate(to, from, next) {
    if(to.params && to.params.id) {
      this.initPage(to.path, next)
    }
    next();
  },
  mounted() {
    console.log(process.env.VUE_APP_BASE_URL)
    let {params} = this.$route;
    if(params && params.id) {
      this.initPage(this.$route.path)
    }
  }
}
</script>
