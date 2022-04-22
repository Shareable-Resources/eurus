<template>
<section class="app-main">
        <titleHeader v-bind:titleData=titleData></titleHeader>
        <router-view :key="$route.path"></router-view>
</section>
</template>
<script>
export default {
    data : function() {
        return {
            titleData : {
                title : "common." + this.$route.name.toLowerCase()
            }
        }
    },
    mounted : function() {
       
    },
    watch : {
        "$route.path" : function() {
            this.titleData.title = "common." + this.$route.name.toLowerCase();
        }
    },

}
</script>

<style>

</style>

