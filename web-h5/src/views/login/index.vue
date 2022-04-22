<template>
<section class="main">
  <!-- <section class="pageHeader">
        <titleHeader v-bind:titleData=titleData></titleHeader>
  </section> -->

  <div class="container">
    <div class="langselect">
      <q-btn color="primary" v-on:click="chanlang('en')">ENG</q-btn>
      <q-btn color="primary" v-on:click="chanlang('zh')">CHI</q-btn>
    </div>
    <br>
    <q-form @submit="saveSession" class="q-gutter-md">
      <div>
        <label>{{$t('common.username')}}</label>
        <q-input outlined :placeholder="$t('login.enter_username')" v-model="login.username" />
      </div>
      <div>
        <label>{{$t('common.password')}}</label>
        <q-input outlined type="password" :placeholder="$t('login.enter_password')" v-model="login.password" />

      </div>
      <br>
      <router-link to="/login">{{$t("login.forgot_password")}}</router-link>
      <br>
      <div class="row justify-center">
        <q-btn label="Login" type="submit" class="col-3 col-xs-4" color="primary" />
      </div>
    </q-form>
    <br>

    <div class="column items-center">
      <p class="col">{{$t("login.do_not_have_account")}} <router-link to="/login">{{$t("common.register")}}
        </router-link>
      </p>


    </div>
  </div>
  </section>


</template>

<script>
  import {
    setToken,
    setName,
    setLocale
  } from '@/utils/auth.js'
  export default {
    name: 'login',

    data() {
      return {
        login: {
          username: '',
          password: '',
        },
        titleData: {
          title: "common.login"
        }


      }
    },
    methods: {
      saveSession:  function () {
        //SEND LOGIN AND PASSWORD HERE (WAIT FOR RESPONSE)

        //SERVER PART END

        //SAVE IT TO COOKIES OR SESSIONSTORAGE
        
        setToken();
        setName();
        this.$router.push("/dashboard");
        
         

      },
      chanlang: function (lang) {
        setLocale(lang)
        this.$i18n.locale = lang;
      }

    },
  }
</script>

<style lang="scss" scoped>
  .login {
    & .title {
      text-align: center;

    }
  }

.langselect {
  padding-top: 2rem;
}
  .container {
    max-width: 1440px;
    margin: 0 auto;

  }
</style>