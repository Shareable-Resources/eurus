import { createApp } from "vue";
import { Vue } from "vue";
import { router } from "../src/router";
import App from "./App.vue";

createApp(App).mount("#app");

var app = new Vue({
    router,
    render: (h) => h(App),
});
app.$mount("#app");