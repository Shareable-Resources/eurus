<template>
  <div id="basicTable">
    <el-table
      v-loading="tableData.loading || false"
      class="basic-table"
      style="width: 100%"
      :border="tableData.border || true"
      :show-summary="tableData.getSummaries ? true : false"
      :summary-method="tableData.getSummaries || null"
      :data="tableData.source"
    >
      <el-table-column
        v-if="tableData.index"
        label="#"
        type="index"
        width="50"
      />
      <template
        v-for="col in tableData.columns"
      >
        <!-- show -->
        <el-table-column
          v-if="col.type === 'show'"
          :key="col.prop"
          :align="col.align || 'center'"
          :width="col.width || ''"
          :label="$t(col.label)"
          :prop="col.prop"
        />
        <!-- filter -->
        <el-table-column
          v-if="col.type === 'filter'"
          :key="`${col.prop}_${col.type}`"
          :align="col.align || 'center'"
          :width="col.width || ''"
          :label="$t(col.label)"
          :prop="col.prop"
        >
          <template slot-scope="props">
            <span>{{
              col.filterFn?
                col.filterFn(props.row[col.prop])
                : ''
            }}</span>
          </template>
        </el-table-column>
        <!-- options -->
        <el-table-column
          v-if="col.type === 'options'"
          :key="col.prop"
          :align="col.align || 'center'"
          :width="col.width || ''"
          :label="$t(col.label)"
          :prop="col.prop"
        >
          <template slot-scope="props">
            <span>{{
              props.row[col.prop]?
                $t(`options.${col.options[props.row[col.prop]]}`)
                : ''
            }}</span>
          </template>
        </el-table-column>
        <!-- time -->
        <el-table-column
          v-if="col.type === 'time'"
          :key="col.prop"
          :align="col.align || 'center'"
          :width="col.width || ''"
          :label="$t(col.label)"
          :prop="col.prop"
        >
          <template slot-scope="props">
            <span>
              {{ props.row[col.prop] | moment('timezone', tz, "YYYY-MM-DD HH:mm:ss") }}
            </span>
          </template>
        </el-table-column>
        <!-- btns -->
        <el-table-column
          v-if="col.type === 'btns'"
          :key="col.prop || col.label"
          :align="col.align || 'center'"
          :width="col.width || '100'"
          :label="$t(col.label)"
        >
          <template
            v-for="btn in col.btns"
            slot-scope="scope"
          >
            <el-button
              v-if="btn.type === 'details'"
              :key="`btn_${btn.type}`"
              size="small"
              @click="handleDetails(btn.fn, scope.row)"
            >{{ $t(`common.details`) }}
            </el-button>
            <el-button
              v-if="btn.type === 'retry'"
              :key="`btn_${btn.type}`"
              size="small"
              @click="handleDetails(btn.fn, scope.row)"
            >{{ $t(`common.retry`) }}
            </el-button>
            <el-button
              v-if="btn.type === 'btn'"
              :key="`btn_${btn.type}`"
              size="small"
              @click="handleDetails(btn.fn, scope.row)"
            >{{ $t(btn.label) }}
            </el-button>
          </template>
        </el-table-column>
      </template>
    </el-table>
    <el-pagination
      class="pagination"
      :current-page="tableData.currentPage"
      :page-sizes="PAGE_SIZES"
      :page-size="PAGE_SIZE"
      layout="total, sizes, prev, pager, next"
      :total="tableData.total"
      @size-change="tableData.handleSizeChange"
      @current-change="tableData.handleCurrentChange"
    />
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
export default {
  name: 'BasicTable',
  props: {
    tableData: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      PAGE_SIZES: [20, 50, 100],
      PAGE_SIZE: 20
    }
  },
  computed: {
    ...mapGetters([
      'tz'
    ])
  },
  methods: {
    handleDetails(fn, row) {
      return fn(row)
    }
  }
}
</script>
<style scoped class="scss">
.pagination {
  margin-top: 10px;
}
</style>
