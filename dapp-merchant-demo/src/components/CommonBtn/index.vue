<template id="common-button">
  <el-button
    :type="btnData.type || ''"
    :size="btnData.size || 'mini'"
    @click="btnData.action"
  >
    {{ $t(btnData.label) ? $t(btnData.label) : '' }}
  </el-button>
</template>

<script>
export default {
  name: 'CommonButton',
  props: {
    btnData: {
      type: Object,
      required: true,
      default: function() {
        return {
          type: '',
          size: 'mini',
          label: ''
        }
      }
    }
  }
}
</script>
