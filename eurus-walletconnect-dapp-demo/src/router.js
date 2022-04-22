import { Vue } from "vue";
import VueRouter from "vue-router";

Vue.use(VueRouter);

export const router = [{
    path: "/connected",
    // name: "connected",
    component: () =>
        import ("../src/views/connected"),
}, ]