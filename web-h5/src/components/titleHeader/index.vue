<template id="titleHeader">
<section class="titleHeader">
  <div class="row items-center justify-center">
      <div class="column ">
          <h5 class="pageTitle">{{$t(titleData.title)}}</h5>
        </div>
  </div>
</section>
</template>

<script>
export default {
    name: "titleHeader",
    props : {
         titleData : Object
    }
}
</script>

<style lang="scss" scoped>

.pageTitle {
    color: white;
    margin: 0;
}

</style>