<template>
  <el-form
    :inline="true"
    size="small"
    :model="filterFormData"
    class="filter-form demo-form-inline"
  >
    <el-form-item
      v-for="(item,index) in filterForm.columns"
      :key="`filterForm_${index}`"
      :label="item.label? $t(item.label): null"
      :prop="item.model"
    >
      <el-date-picker
        v-if="item.type ==='month'"
        v-model="filterFormData[item.model]"
        type="month"
        :format="item.format"
      />

      <el-date-picker
        v-if="item.type ==='date'"
        v-model="filterFormData[item.model]"
        type="date"
        :format="item.format"
      />

      <el-date-picker
        v-if="item.type ==='datetime'"
        v-model="filterFormData[item.model]"
        type="datetime"
        :format="item.format"
      />

      <el-date-picker
        v-if="item.type ==='datetimerange'"
        v-model="filterFormData[item.model]"
        type="datetimerange"
        range-separator="-"
        start-placeholder=""
        end-placeholder=""
      />

      <template
        v-if="item.type ==='snapShot'"
      >
        <el-col :span="11">
          <el-form-item prop="snapShotDate" class="snap-shot">
            <el-select
              v-model="filterFormData.snapShotDate"
              :placeholder="item.placeholder? $t(item.placeholder): null"
              style="width: 100%;"
            >
              <el-option
                v-for="option in SEVEN_DAYS"
                :key="`date_${option}`"
                :label="option"
                :value="option"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col class="line" :span="1">-</el-col>
        <el-col :span="11">
          <el-form-item prop="snapShotTime" class="snap-shot">
            <el-select
              v-model="filterFormData.snapShotTime"
              :placeholder="item.placeholder? $t(item.placeholder): null"
              style="width: 100%;"
            >
              <el-option
                v-for="count in 24"
                :key="`time_${count}`"
                :label="count < 10 ?`0${count}:00:00`: `${count}:00:00`"
                :value="count < 10 ?`0${count}:00:00`: `${count}:00:00`"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </template>

      <el-select
        v-if="item.type ==='select'"
        v-model="filterFormData[item.model]"
        clearable
        filterable
        :placeholder="item.placeholder? $t(item.placeholder): null"
      >
        <el-option
          v-for="(optionLabel, optionVal) in item.options"
          :key="optionVal"
          :label="$t(`options.${optionLabel}`)"
          :value="optionVal"
        />
      </el-select>

      <el-input
        v-if="item.type ==='input'"
        v-model="filterFormData[item.model]"
        :placeholder="item.placeholder? $t(item.placeholder): null"
      />
    </el-form-item>

    <el-form-item>
      <el-button
        v-if="filterForm.searchAction"
        type="primary"
        @click="onSearch"
      >{{ $t(`common.search`) }}
      </el-button>
      <el-button
        v-if="filterForm.searchAction"
        @click="initFilterForm"
      >{{ $t(`common.reset`) }}
      </el-button>
      <el-button
        v-if="filterForm.addAction"
        type="primary"
        @click="onAdd"
      >{{ $t(`common.create`) }}
      </el-button>
      <el-button
        v-if="filterForm.exportAction"
        @click="onExport"
      >{{ $t(`common.export`) }}
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import { deepClone } from '@/utils'
export default {
  name: 'FilterForm',
  props: {
    filterForm: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      filterFormData: {},
      SEVEN_DAYS: this.get7days()
    }
  },
  computed: {
    snapShot: function() {
      const { snapShotDate, snapShotTime } = this.filterFormData
      if (!snapShotDate || !snapShotTime) return ''
      return `${snapShotDate}-${snapShotTime}`
    }
  },
  mounted() {
    this.initFilterForm()
  },
  methods: {
    get7days() {
      const that = this
      const arr = []
      for (let i = 1; i <= 7; i++) {
        arr.push(that.dayBeforeDate(i))
      }
      return arr
    },
    dayBeforeDate(dayBefore) {
      const today = new Date()
      const date = new Date(today.getTime() - (dayBefore * 24 * 60 * 60 * 1000))
      return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
    },
    initFilterForm() {
      this.filterFormData = deepClone(this.filterForm.default)
    },
    onSearch() {
      return this.filterForm.searchAction({
        ...this.filterFormData,
        snapShot: this.snapShot
      })
    },
    onAdd() {
      return this.filterForm.addAction(this.filterFormData)
    },
    onExport() {
      return this.filterForm.exportAction(this.filterFormData)
    }
  }
}
</script>

<style lang="scss" scoped>
.snap-shot {
  margin-bottom: 0px;
}
</style>
