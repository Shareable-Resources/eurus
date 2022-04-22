<template id="inputField">
    <q-input
        outlined 
        :placeholder="$t('login.enter_username')"
      />
</template>

<script>
export default {
    name: "inputField",
    props: {
        fidld: Object
    },
    data : function(){
        return {
            fieldData :{}
        }
    }
}
</script>

<style>

</style>