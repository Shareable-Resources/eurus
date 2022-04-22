<template>

        <div>
            <h1>Dashboard</h1>
            <p>Welcome {{username}}</p>
            <button v-on:click="logout()">Logout</button>
        </div>


</template>

<script>
    import {
        getName,
        clearName,
        clearToken
    } from '@/utils/auth.js'

    export default {
        name: 'login',
        data() {
            return {
                username: getName(),
            }
        },
        methods: {
            logout: function () {
                clearName()
                clearToken()
                this.$router.go('/login')
            }
        }
    }
</script>

<style lang="scss">

</style>